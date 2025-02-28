package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/internal/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var mutexKV = mutexkv.NewMutexKV()

func ResourceAbrhaDatabaseUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseUserCreate,
		ReadContext:   resourceAbrhaDatabaseUserRead,
		UpdateContext: resourceAbrhaDatabaseUserUpdate,
		DeleteContext: resourceAbrhaDatabaseUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseUserImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"mysql_auth_plugin": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					goApiAbrha.SQLAuthPluginNative,
					goApiAbrha.SQLAuthPluginCachingSHA2,
				}, false),
				// Prevent diffs when default is used and not specified in the config.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == goApiAbrha.SQLAuthPluginCachingSHA2 && new == ""
				},
			},
			"settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     userACLSchema(),
						},
						"opensearch_acl": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     userOpenSearchACLSchema(),
						},
					},
				},
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"access_cert": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func userACLSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"permission": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"admin",
					"consume",
					"produce",
					"produceconsume",
				}, false),
			},
		},
	}
}

func userOpenSearchACLSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"index": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"permission": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"deny",
					"admin",
					"read",
					"write",
					"readwrite",
				}, false),
			},
		},
	}
}

func resourceAbrhaDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &goApiAbrha.DatabaseCreateUserRequest{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("mysql_auth_plugin"); ok {
		opts.MySQLSettings = &goApiAbrha.DatabaseMySQLUserSettings{
			AuthPlugin: v.(string),
		}
	}

	if v, ok := d.GetOk("settings"); ok {
		opts.Settings = expandUserSettings(v.([]interface{}))
	}

	// Prevent parallel creation of users for same cluster.
	key := fmt.Sprintf("abrha_database_cluster/%s/users", clusterID)
	mutexKV.Lock(key)
	defer mutexKV.Unlock(key)

	log.Printf("[DEBUG] Database User create configuration: %#v", opts)
	user, _, err := client.Databases.CreateUser(context.Background(), clusterID, opts)
	if err != nil {
		return diag.Errorf("Error creating Database User: %s", err)
	}

	d.SetId(makeDatabaseUserID(clusterID, user.Name))
	log.Printf("[INFO] Database User Name: %s", user.Name)

	// set userSettings only on CreateUser, due to CreateUser responses including `settings` but GetUser responses not including `settings`
	if err := d.Set("settings", flattenUserSettings(user.Settings)); err != nil {
		return diag.Errorf("Error setting user settings: %#v", err)
	}

	setDatabaseUserAttributes(d, user)

	return nil
}

func resourceAbrhaDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	// Check if the database user still exists
	user, resp, err := client.Databases.GetUser(context.Background(), clusterID, name)
	if err != nil {
		// If the database user is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Database User: %s", err)
	}

	setDatabaseUserAttributes(d, user)

	return nil
}

func setDatabaseUserAttributes(d *schema.ResourceData, user *goApiAbrha.DatabaseUser) {
	// Default to "normal" when not set.
	if user.Role == "" {
		d.Set("role", "normal")
	} else {
		d.Set("role", user.Role)
	}

	// This will be blank when GETing MongoDB clusters post-create.
	// Don't overwrite the password set on create.
	if user.Password != "" {
		d.Set("password", user.Password)
	}

	if user.MySQLSettings != nil {
		d.Set("mysql_auth_plugin", user.MySQLSettings.AuthPlugin)
	}

	// This is only set for kafka users
	if user.AccessCert != "" {
		d.Set("access_cert", user.AccessCert)
	}
	if user.AccessKey != "" {
		d.Set("access_key", user.AccessKey)
	}
}

func resourceAbrhaDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if d.HasChange("mysql_auth_plugin") {
		authReq := &goApiAbrha.DatabaseResetUserAuthRequest{}
		if d.Get("mysql_auth_plugin").(string) != "" {
			authReq.MySQLSettings = &goApiAbrha.DatabaseMySQLUserSettings{
				AuthPlugin: d.Get("mysql_auth_plugin").(string),
			}
		} else {
			// If blank, restore default value.
			authReq.MySQLSettings = &goApiAbrha.DatabaseMySQLUserSettings{
				AuthPlugin: goApiAbrha.SQLAuthPluginCachingSHA2,
			}
		}

		_, _, err := client.Databases.ResetUserAuth(context.Background(), d.Get("cluster_id").(string), d.Get("name").(string), authReq)
		if err != nil {
			return diag.Errorf("Error updating mysql_auth_plugin for DatabaseUser: %s", err)
		}
	}
	if d.HasChange("settings") {
		updateReq := &goApiAbrha.DatabaseUpdateUserRequest{}
		if v, ok := d.GetOk("settings"); ok {
			updateReq.Settings = expandUserSettings(v.([]interface{}))
		}
		_, _, err := client.Databases.UpdateUser(context.Background(), d.Get("cluster_id").(string), d.Get("name").(string), updateReq)
		if err != nil {
			return diag.Errorf("Error updating settings for DatabaseUser: %s", err)
		}
	}

	return resourceAbrhaDatabaseUserRead(ctx, d, meta)
}

func resourceAbrhaDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	// Prevent parallel deletion of users for same cluster.
	key := fmt.Sprintf("abrha_database_cluster/%s/users", clusterID)
	mutexKV.Lock(key)
	defer mutexKV.Unlock(key)

	log.Printf("[INFO] Deleting Database User: %s", d.Id())
	_, err := client.Databases.DeleteUser(context.Background(), clusterID, name)
	if err != nil {
		return diag.Errorf("Error deleting Database User: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceAbrhaDatabaseUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeDatabaseUserID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	} else {
		return nil, errors.New("must use the ID of the source database cluster and the name of the user joined with a comma (e.g. `id,name`)")
	}

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseUserID(clusterID string, name string) string {
	return fmt.Sprintf("%s/user/%s", clusterID, name)
}

func expandUserSettings(raw []interface{}) *goApiAbrha.DatabaseUserSettings {
	if len(raw) == 0 || raw[0] == nil {
		return &goApiAbrha.DatabaseUserSettings{}
	}
	userSettingsConfig := raw[0].(map[string]interface{})

	userSettings := &goApiAbrha.DatabaseUserSettings{
		ACL:           expandUserACLs(userSettingsConfig["acl"].([]interface{})),
		OpenSearchACL: expandOpenSearchUserACLs(userSettingsConfig["opensearch_acl"].([]interface{})),
	}
	return userSettings
}

func expandUserACLs(rawACLs []interface{}) []*goApiAbrha.KafkaACL {
	acls := make([]*goApiAbrha.KafkaACL, 0, len(rawACLs))
	for _, rawACL := range rawACLs {
		a := rawACL.(map[string]interface{})
		acl := &goApiAbrha.KafkaACL{
			Topic:      a["topic"].(string),
			Permission: a["permission"].(string),
		}
		acls = append(acls, acl)
	}
	return acls
}

func expandOpenSearchUserACLs(rawACLs []interface{}) []*goApiAbrha.OpenSearchACL {
	acls := make([]*goApiAbrha.OpenSearchACL, 0, len(rawACLs))
	for _, rawACL := range rawACLs {
		a := rawACL.(map[string]interface{})
		acl := &goApiAbrha.OpenSearchACL{
			Index:      a["index"].(string),
			Permission: a["permission"].(string),
		}
		acls = append(acls, acl)
	}
	return acls
}

func flattenUserSettings(settings *goApiAbrha.DatabaseUserSettings) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if settings != nil {
		r := make(map[string]interface{})
		r["acl"] = flattenUserACLs(settings.ACL)
		r["opensearch_acl"] = flattenOpenSearchUserACLs(settings.OpenSearchACL)
		result = append(result, r)
	}
	return result
}

func flattenUserACLs(acls []*goApiAbrha.KafkaACL) []map[string]interface{} {
	result := make([]map[string]interface{}, len(acls))
	for i, acl := range acls {
		item := make(map[string]interface{})
		item["id"] = acl.ID
		item["topic"] = acl.Topic
		item["permission"] = normalizePermission(acl.Permission)
		result[i] = item
	}
	return result
}

func flattenOpenSearchUserACLs(acls []*goApiAbrha.OpenSearchACL) []map[string]interface{} {
	result := make([]map[string]interface{}, len(acls))
	for i, acl := range acls {
		item := make(map[string]interface{})
		item["index"] = acl.Index
		item["permission"] = normalizeOpenSearchPermission(acl.Permission)
		result[i] = item
	}
	return result
}

func normalizePermission(p string) string {
	pLower := strings.ToLower(p)
	switch pLower {
	case "admin", "produce", "consume":
		return pLower
	case "produceconsume", "produce_consume", "readwrite", "read_write":
		return "produceconsume"
	}
	return ""
}

func normalizeOpenSearchPermission(p string) string {
	pLower := strings.ToLower(p)
	switch pLower {
	case "admin", "deny", "read", "write":
		return pLower
	case "readwrite", "read_write":
		return "readwrite"
	}
	return ""
}
