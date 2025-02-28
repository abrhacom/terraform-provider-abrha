package certificate

import (
	"context"
	"fmt"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaCertificateRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the certificate",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes

			// When the certificate type is lets_encrypt, the certificate
			// ID will change when it's renewed, so we have to rely on the
			// certificate name as the primary identifier instead.
			// We include the UUID as another computed field for use in the
			// short-term refresh function that waits for it to be ready.
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "uuid of the certificate",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of the certificate",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "current state of the certificate",
			},
			"domains": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "domains for which the certificate was issued",
			},
			"not_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expiration date and time of the certificate",
			},
			"sha1_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA1 fingerprint of the certificate",
			},
		},
	}
}

func dataSourceAbrhaCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	// When the certificate type is lets_encrypt, the certificate
	// ID will change when it's renewed, so we have to rely on the
	// certificate name as the primary identifier instead.
	name := d.Get("name").(string)
	cert, err := FindCertificateByName(client, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cert.Name)
	d.Set("name", cert.Name)
	d.Set("uuid", cert.ID)
	d.Set("type", cert.Type)
	d.Set("state", cert.State)
	d.Set("not_after", cert.NotAfter)
	d.Set("sha1_fingerprint", cert.SHA1Fingerprint)

	if err := d.Set("domains", flattenAbrhaCertificateDomains(cert.DNSNames)); err != nil {
		return diag.Errorf("Error setting `domain`: %+v", err)
	}

	return nil
}

func FindCertificateByName(client *goApiAbrha.Client, name string) (*goApiAbrha.Certificate, error) {
	cert, _, err := client.Certificates.ListByName(context.Background(), name, nil)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving certificates: %s", err)
	}

	if len(cert) == 0 {
		return nil, fmt.Errorf("certificate not found")
	}

	return &cert[0], err
}
