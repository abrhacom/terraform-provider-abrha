package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaDatabaseConnectionPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseConnectionPoolCreate,
		ReadContext:   resourceAbrhaDatabaseConnectionPoolRead,
		DeleteContext: resourceAbrhaDatabaseConnectionPoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseConnectionPoolImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 63),
			},

			"user": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"session",
					"transaction",
					"statement"}, false),
			},

			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"db_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"private_uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAbrhaDatabaseConnectionPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	clusterID := d.Get("cluster_id").(string)
	opts := &goApiAbrha.DatabaseCreatePoolRequest{
		Name:     d.Get("name").(string),
		User:     d.Get("user").(string),
		Mode:     d.Get("mode").(string),
		Size:     d.Get("size").(int),
		Database: d.Get("db_name").(string),
	}

	log.Printf("[DEBUG] DatabaseConnectionPool create configuration: %#v", opts)
	pool, _, err := client.Databases.CreatePool(context.Background(), clusterID, opts)
	if err != nil {
		return diag.Errorf("Error creating DatabaseConnectionPool: %s", err)
	}

	d.SetId(createConnectionPoolID(clusterID, pool.Name))
	log.Printf("[INFO] DatabaseConnectionPool Name: %s", pool.Name)

	err = setConnectionPoolInfo(pool, d)
	if err != nil {
		return diag.Errorf("Error building connection URI: %s", err)
	}

	return resourceAbrhaDatabaseConnectionPoolRead(ctx, d, meta)
}

func resourceAbrhaDatabaseConnectionPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID, poolName := splitConnectionPoolID(d.Id())

	pool, resp, err := client.Databases.GetPool(context.Background(), clusterID, poolName)
	if err != nil {
		// If the pool is somehow already destroyed, mark as
		// successfully gone
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving DatabaseConnectionPool: %s", err)
	}

	d.SetId(createConnectionPoolID(clusterID, pool.Name))
	d.Set("cluster_id", clusterID)
	d.Set("name", pool.Name)
	d.Set("user", pool.User)
	d.Set("mode", pool.Mode)
	d.Set("size", pool.Size)
	d.Set("db_name", pool.Database)

	err = setConnectionPoolInfo(pool, d)
	if err != nil {
		return diag.Errorf("Error building connection URI: %s", err)
	}

	return nil
}

func setConnectionPoolInfo(pool *goApiAbrha.DatabasePool, d *schema.ResourceData) error {
	if pool.Connection != nil {
		d.Set("host", pool.Connection.Host)
		d.Set("port", pool.Connection.Port)

		if pool.Connection.Password != "" {
			d.Set("password", pool.Connection.Password)
		}

		uri, err := buildDBConnectionURI(pool.Connection, d)
		if err != nil {
			return err
		}

		d.Set("uri", uri)
	}

	if pool.PrivateConnection != nil {
		d.Set("private_host", pool.PrivateConnection.Host)

		privateURI, err := buildDBConnectionURI(pool.PrivateConnection, d)
		if err != nil {
			return err
		}

		d.Set("private_uri", privateURI)
	}

	return nil
}

func resourceAbrhaDatabaseConnectionPoolImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(createConnectionPoolID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	} else {
		return nil, errors.New("must use the ID of the source database cluster and the name of the connection pool joined with a comma (e.g. `id,name`)")
	}

	return []*schema.ResourceData{d}, nil
}

func resourceAbrhaDatabaseConnectionPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID, poolName := splitConnectionPoolID(d.Id())

	log.Printf("[INFO] Deleting DatabaseConnectionPool: %s", poolName)
	_, err := client.Databases.DeletePool(context.Background(), clusterID, poolName)
	if err != nil {
		return diag.Errorf("Error deleting DatabaseConnectionPool: %s", err)
	}

	d.SetId("")
	return nil
}

func createConnectionPoolID(clusterID string, poolName string) string {
	return fmt.Sprintf("%s/%s", clusterID, poolName)
}

func splitConnectionPoolID(id string) (string, string) {
	splitID := strings.Split(id, "/")
	clusterID := splitID[0]
	poolName := splitID[1]

	return clusterID, poolName
}
