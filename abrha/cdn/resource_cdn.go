package cdn

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/certificate"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	needsCloudflareCert = "needs-cloudflare-cert"
)

func ResourceAbrhaCDN() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaCDNCreate,
		ReadContext:   resourceAbrhaCDNRead,
		UpdateContext: resourceAbrhaCDNUpdate,
		DeleteContext: resourceAbrhaCDNDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAbrhaCDNv0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCDNStateV0toV1,
				Version: 0,
			},
		},

		Schema: resourceAbrhaCDNv1(),
	}
}

func resourceAbrhaCDNv1() map[string]*schema.Schema {
	cdnV1Schema := map[string]*schema.Schema{
		"certificate_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}

	for k, v := range resourceAbrhaCDNv0().Schema {
		cdnV1Schema[k] = v
	}
	cdnV1Schema["certificate_id"].Computed = true
	cdnV1Schema["certificate_id"].Deprecated = "Certificate IDs may change, for example when a Let's Encrypt certificate is auto-renewed. Please specify 'certificate_name' instead."

	return cdnV1Schema
}

func resourceAbrhaCDNv0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"origin": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "fully qualified domain name (FQDN) for the origin server",
				ValidateFunc: validation.NoZeroValues,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The amount of time the content is cached in the CDN",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of a Abrha managed TLS certificate for use with custom domains",
			},
			"custom_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "fully qualified domain name (FQDN) for custom subdomain, (requires certificate_id)",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "fully qualified domain name (FQDN) to serve the CDN content",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time (ISO8601) of when the CDN endpoint was created.",
			},
		},
	}
}

func migrateCDNStateV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if len(rawState) == 0 {
		log.Println("[DEBUG] Empty state; nothing to migrate.")
		return rawState, nil
	}

	// When the certificate type is lets_encrypt, the certificate
	// ID will change when it's renewed, so we have to rely on the
	// certificate name as the primary identifier instead.
	certID := rawState["certificate_id"].(string)
	if certID != "" {
		log.Println("[DEBUG] Migrating CDN schema from v0 to v1.")
		client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
		cert, _, err := client.Certificates.Get(context.Background(), certID)
		if err != nil {
			return rawState, err
		}

		rawState["certificate_id"] = cert.Name
		rawState["certificate_name"] = cert.Name
	}

	return rawState, nil
}

func resourceAbrhaCDNCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	cdnRequest := &goApiAbrha.CDNCreateRequest{
		Origin: d.Get("origin").(string),
	}

	if v, ok := d.GetOk("ttl"); ok {
		cdnRequest.TTL = uint32(v.(int))
	}

	if v, ok := d.GetOk("custom_domain"); ok {
		cdnRequest.CustomDomain = v.(string)
	}

	if name, nameOk := d.GetOk("certificate_name"); nameOk {
		certName := name.(string)
		if certName != "" {
			if certName == needsCloudflareCert {
				cdnRequest.CertificateID = needsCloudflareCert
			} else {
				cert, err := certificate.FindCertificateByName(client, certName)
				if err != nil {
					return diag.FromErr(err)
				}
				cdnRequest.CertificateID = cert.ID
			}
		}
	}

	if id, idOk := d.GetOk("certificate_id"); idOk && cdnRequest.CertificateID == "" {
		// When the certificate type is lets_encrypt, the certificate
		// ID will change when it's renewed, so we have to rely on the
		// certificate name as the primary identifier instead.
		certName := id.(string)
		if certName != "" {
			cert, err := certificate.FindCertificateByName(client, certName)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					log.Println("[DEBUG] Certificate not found looking up by name. Falling back to lookup by ID.")
					cert, _, err = client.Certificates.Get(context.Background(), certName)
					if err != nil {
						return diag.FromErr(err)
					}
				} else {
					return diag.FromErr(err)
				}
			}

			cdnRequest.CertificateID = cert.ID
		}
	}

	log.Printf("[DEBUG] CDN create request: %#v", cdnRequest)
	cdn, _, err := client.CDNs.Create(context.Background(), cdnRequest)
	if err != nil {
		return diag.Errorf("Error creating CDN: %s", err)
	}

	d.SetId(cdn.ID)
	log.Printf("[INFO] CDN created, ID: %s", d.Id())

	return resourceAbrhaCDNRead(ctx, d, meta)
}

func resourceAbrhaCDNRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	cdn, resp, err := getCDNWithRetryBackoff(ctx, client, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] CDN  (%s) was not found - removing from state", d.Id())
			d.SetId("")
		}
		return diag.Errorf("Error reading CDN: %s", err)
	}

	d.SetId(cdn.ID)
	d.Set("origin", cdn.Origin)
	d.Set("ttl", cdn.TTL)
	d.Set("endpoint", cdn.Endpoint)
	d.Set("created_at", cdn.CreatedAt.UTC().String())
	d.Set("custom_domain", cdn.CustomDomain)

	if cdn.CertificateID != "" && cdn.CertificateID != needsCloudflareCert {
		// When the certificate type is lets_encrypt, the certificate
		// ID will change when it's renewed, so we have to rely on the
		// certificate name as the primary identifier instead.
		cert, _, err := client.Certificates.Get(context.Background(), cdn.CertificateID)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("certificate_id", cert.Name)
		d.Set("certificate_name", cert.Name)
	}

	if cdn.CertificateID == needsCloudflareCert {
		d.Set("certificate_id", cdn.CertificateID)
		d.Set("certificate_name", cdn.CertificateID)
	}

	return nil
}

func getCDNWithRetryBackoff(ctx context.Context, client *goApiAbrha.Client, id string) (*goApiAbrha.CDN, *goApiAbrha.Response, error) {
	var (
		cdn     *goApiAbrha.CDN
		resp    *goApiAbrha.Response
		timeout = 30 * time.Second
		err     error
	)
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		cdn, resp, err = client.CDNs.Get(ctx, id)
		if err != nil {
			if util.IsAbrhaError(err, http.StatusNotFound, "") {
				log.Printf("[DEBUG] Received %s, retrying CDN", err.Error())
				return retry.RetryableError(err)
			}

			if util.IsAbrhaError(err, http.StatusTooManyRequests, "") {
				log.Printf("[DEBUG] Received %s, backing off", err.Error())
				time.Sleep(10 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return nil, resp, err
	}

	return cdn, resp, nil
}

func resourceAbrhaCDNUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	d.Partial(true)

	if d.HasChange("ttl") {
		ttlUpdateRequest := &goApiAbrha.CDNUpdateTTLRequest{
			TTL: uint32(d.Get("ttl").(int)),
		}
		_, _, err := client.CDNs.UpdateTTL(context.Background(), d.Id(), ttlUpdateRequest)

		if err != nil {
			return diag.Errorf("Error updating CDN TTL: %s", err)
		}
		log.Printf("[INFO] Updated TTL on CDN")
	}

	if d.HasChanges("certificate_id", "custom_domain", "certificate_name") {
		cdnUpdateRequest := &goApiAbrha.CDNUpdateCustomDomainRequest{
			CustomDomain: d.Get("custom_domain").(string),
		}

		certName := d.Get("certificate_name").(string)
		if certName != "" {
			if certName == needsCloudflareCert {
				cdnUpdateRequest.CertificateID = needsCloudflareCert
			} else {
				cert, err := certificate.FindCertificateByName(client, certName)
				if err != nil {
					return diag.FromErr(err)
				}

				cdnUpdateRequest.CertificateID = cert.ID
			}
		}

		_, _, err := client.CDNs.UpdateCustomDomain(context.Background(), d.Id(), cdnUpdateRequest)

		if err != nil {
			return diag.Errorf("Error updating CDN custom domain: %s", err)
		}
		log.Printf("[INFO] Updated custom domain/certificate on CDN")
	}

	d.Partial(false)
	return resourceAbrhaCDNRead(ctx, d, meta)
}

func resourceAbrhaCDNDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	resourceID := d.Id()

	timeout := 30 * time.Second
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err := client.CDNs.Delete(context.Background(), resourceID)
		if err != nil {
			if util.IsAbrhaError(err, http.StatusTooManyRequests, "") {
				log.Printf("[DEBUG] Received %s, backing off", err.Error())
				time.Sleep(10 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error deleting CDN: %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] CDN deleted, ID: %s", resourceID)

	return nil
}
