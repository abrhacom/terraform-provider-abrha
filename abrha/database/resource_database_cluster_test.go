package database_test

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"testing"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaDatabaseCluster_Basic(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "engine", "pg"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "private_host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "port"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "user"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "password"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "urn"),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "tags.#", "1"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "private_network_uuid"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "project_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "storage_size_mib"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "uri"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "private_uri"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_KafkaConnectionDetails(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterKafka, databaseName, "3.7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "engine", "kafka"),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "port", "25073"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "host"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_WithUpdate(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "size", "db-s-1vcpu-2gb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName),
				Check: resource.TestCheckFunc(
					func(s *terraform.State) error {
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "size", "db-s-2vcpu-4gb"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_WithAdditionalStorage(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "storage_size_mib", "30720"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName),
				Check: resource.TestCheckFunc(
					func(s *terraform.State) error {
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithAdditionalStorage, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "storage_size_mib", "61440"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_WithMigration(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "region", "nyc1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithMigration, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "region", "lon1"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_WithMaintWindow(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithMaintWindow, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "maintenance_window.0.day"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "maintenance_window.0.hour"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_WithSQLMode(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithSQLMode, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "sql_mode",
						"ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ZERO_DATE,NO_ZERO_IN_DATE"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithSQLModeUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "sql_mode",
						"ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_CheckSQLModeSupport(t *testing.T) {
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithRedisSQLModeError, databaseName),
				ExpectError: regexp.MustCompile(`sql_mode is only supported for MySQL`),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_RedisNoVersion(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterRedisNoVersion, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "engine", "redis"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "uri"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "private_uri"),
				),
				ExpectError: regexp.MustCompile(`The argument "version" is required, but no definition was found.`),
			},
		},
	})
}

// Abrha only supports one version of Redis. For backwards compatibility
// the API allows for POST requests that specifies a previous version, but new
// clusters are created with the latest/only supported version, regardless of
// the version specified in the config.
// The provider suppresses diffs when the config version is <= to the latest
// version. New clusters is always created with the latest version .
func TestAccAbrhaDatabaseCluster_oldRedisVersion(t *testing.T) {
	var (
		database goApiAbrha.Database
	)

	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterRedis, databaseName, "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "engine", "redis"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "version"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_RedisWithEvictionPolicy(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			// Create with an eviction policy
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithEvictionPolicy, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "eviction_policy", "volatile_random"),
				),
			},
			// Update eviction policy
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithEvictionPolicyUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "eviction_policy", "allkeys_lru"),
				),
			},
			// Remove eviction policy
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterRedis, databaseName, "6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_CheckEvictionPolicySupport(t *testing.T) {
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithEvictionPolicyError, databaseName),
				ExpectError: regexp.MustCompile(`eviction_policy is only supported for Redis`),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_TagUpdate(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "tags.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigTagUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_WithVPC(t *testing.T) {
	var database goApiAbrha.Database
	vpcName := acceptance.RandomTestName()
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithVPC, vpcName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "private_network_uuid"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_WithBackupRestore(t *testing.T) {
	var originalDatabase goApiAbrha.Database
	var backupDatabase goApiAbrha.Database

	originalDatabaseName := acceptance.RandomTestName()
	backupDatabasename := acceptance.RandomTestName()

	originalDatabaseConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, originalDatabaseName)
	backUpRestoreConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithBackupRestore, backupDatabasename, originalDatabaseName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: originalDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &originalDatabase),
					testAccCheckAbrhaDatabaseClusterAttributes(&originalDatabase, originalDatabaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "region", "nyc1"),
					func(s *terraform.State) error {
						err := waitForDatabaseBackups(originalDatabaseName)
						return err
					},
				),
			},
			{
				Config: originalDatabaseConfig + backUpRestoreConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar_backup", &backupDatabase),
					testAccCheckAbrhaDatabaseClusterAttributes(&backupDatabase, backupDatabasename),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar_backup", "region", "nyc1"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_MongoDBPassword(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigMongoDB, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists(
						"abrha_database_cluster.foobar", &database),
					resource.TestCheckResourceAttrSet(
						"abrha_database_cluster.foobar", "password"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "uri"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "private_uri"),
				),
			},
			// Pause before running CheckDestroy
			{
				Config: " ",
				Check: resource.TestCheckFunc(
					func(s *terraform.State) error {
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_Upgrade(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()
	previousPGVersion := "14"
	latestPGVersion := "15"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				// TODO: Hardcoding the versions here is not ideal.
				// We will need to determine a better way to fetch the last and latest versions dynamically.
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigCustomVersion, databaseName, "pg", previousPGVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists(
						"abrha_database_cluster.foobar", &database),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "engine", "pg"),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "version", previousPGVersion),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigCustomVersion, databaseName, "pg", latestPGVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "version", latestPGVersion),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseCluster_nonDefaultProject(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()
	projectName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigNonDefaultProject, projectName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"abrha_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttrPair(
						"abrha_project.foobar", "id", "abrha_database_cluster.foobar", "project_id"),
				),
			},
		},
	})
}

func testAccCheckAbrhaDatabaseClusterDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_database_cluster" {
			continue
		}

		// Try to find the database
		_, _, err := client.Databases.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("DatabaseCluster still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaDatabaseClusterAttributes(database *goApiAbrha.Database, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if database.Name != name {
			return fmt.Errorf("Bad name: %s", database.Name)
		}

		return nil
	}
}

func testAccCheckAbrhaDatabaseClusterExists(n string, database *goApiAbrha.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DatabaseCluster ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundDatabaseCluster, _, err := client.Databases.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDatabaseCluster.ID != rs.Primary.ID {
			return fmt.Errorf("DatabaseCluster not found")
		}

		*database = *foundDatabaseCluster

		return nil
	}
}

// testAccCheckAbrhaDatabaseClusterURIPassword checks that the password in
// a database cluster's URI or private URI matches the password value stored in
// its password attribute.
func testAccCheckAbrhaDatabaseClusterURIPassword(name string, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		uri, ok := rs.Primary.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("%s not set", attributeName)
		}

		u, err := url.Parse(uri)
		if err != nil {
			return err
		}

		password, ok := u.User.Password()
		if !ok || password == "" {
			return fmt.Errorf("password not set in %s: %s", attributeName, uri)
		}

		return resource.TestCheckResourceAttr(name, "password", password)(s)
	}
}

func waitForDatabaseBackups(originalDatabaseName string) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	var (
		tickerInterval = 10 * time.Second
		timeoutSeconds = 300.0
		timeout        = int(timeoutSeconds / tickerInterval.Seconds())
		n              = 0
		ticker         = time.NewTicker(tickerInterval)
	)

	databases, _, err := client.Databases.List(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error retrieving backups from original cluster")
	}

	// gets original database's ID
	var originalDatabaseID string
	for _, db := range databases {
		if db.Name == originalDatabaseName {
			originalDatabaseID = db.ID
		}
	}

	if originalDatabaseID == "" {
		return fmt.Errorf("Error retrieving backups from cluster")
	}

	for range ticker.C {
		backups, resp, err := client.Databases.ListBackups(context.Background(), originalDatabaseID, nil)
		if resp.StatusCode == 412 {
			continue
		}

		if err != nil {
			ticker.Stop()
			return fmt.Errorf("Error retrieving backups from cluster")
		}

		if len(backups) >= 1 {
			ticker.Stop()
			return nil
		}

		if n > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return fmt.Errorf("Timeout waiting for database cluster to have a backup to be restored from")
}

const testAccCheckAbrhaDatabaseClusterConfigBasic = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterConfigWithBackupRestore = `
resource "abrha_database_cluster" "foobar_backup" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]

  backup_restore {
    database_name = "%s"
  }
}`

const testAccCheckAbrhaDatabaseClusterConfigWithUpdate = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-2vcpu-4gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterConfigWithAdditionalStorage = `
resource "abrha_database_cluster" "foobar" {
  name             = "%s"
  engine           = "pg"
  version          = "15"
  size             = "db-s-1vcpu-2gb"
  region           = "nyc1"
  node_count       = 1
  tags             = ["production"]
  storage_size_mib = 61440
}`

const testAccCheckAbrhaDatabaseClusterConfigWithMigration = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "lon1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterConfigWithMaintWindow = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]

  maintenance_window {
    day  = "friday"
    hour = "13:00"
  }
}`

const testAccCheckAbrhaDatabaseClusterConfigWithSQLMode = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "lon1"
  node_count = 1
  sql_mode   = "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ZERO_DATE,NO_ZERO_IN_DATE"
}`

const testAccCheckAbrhaDatabaseClusterConfigWithSQLModeUpdate = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "lon1"
  node_count = 1
  sql_mode   = "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE"
}`

const testAccCheckAbrhaDatabaseClusterConfigWithRedisSQLModeError = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "redis"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  sql_mode   = "ANSI"
}`

const testAccCheckAbrhaDatabaseClusterRedisNoVersion = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "redis"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterRedis = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "redis"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterKafka = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "kafka"
  version    = "%s"
  size       = "db-s-2vcpu-2gb"
  region     = "nyc1"
  node_count = 3
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterMySQL = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterPostgreSQL = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterMongoDB = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mongodb"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterOpensearch = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "opensearch"
  version    = "%s"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckAbrhaDatabaseClusterConfigWithEvictionPolicy = `
resource "abrha_database_cluster" "foobar" {
  name            = "%s"
  engine          = "redis"
  version         = "5"
  size            = "db-s-1vcpu-1gb"
  region          = "nyc1"
  node_count      = 1
  tags            = ["production"]
  eviction_policy = "volatile_random"
}
`

const testAccCheckAbrhaDatabaseClusterConfigWithEvictionPolicyUpdate = `
resource "abrha_database_cluster" "foobar" {
  name            = "%s"
  engine          = "redis"
  version         = "5"
  size            = "db-s-1vcpu-1gb"
  region          = "nyc1"
  node_count      = 1
  tags            = ["production"]
  eviction_policy = "allkeys_lru"
}
`

const testAccCheckAbrhaDatabaseClusterConfigWithEvictionPolicyError = `
resource "abrha_database_cluster" "foobar" {
  name            = "%s"
  engine          = "pg"
  version         = "15"
  size            = "db-s-1vcpu-1gb"
  region          = "nyc1"
  node_count      = 1
  eviction_policy = "allkeys_lru"
}
`

const testAccCheckAbrhaDatabaseClusterConfigTagUpdate = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production", "foo"]
}`

const testAccCheckAbrhaDatabaseClusterConfigWithVPC = `
resource "abrha_vpc" "foobar" {
  name   = "%s"
  region = "nyc1"
}

resource "abrha_database_cluster" "foobar" {
  name                 = "%s"
  engine               = "pg"
  version              = "15"
  size                 = "db-s-1vcpu-2gb"
  region               = "nyc1"
  node_count           = 1
  tags                 = ["production"]
  private_network_uuid = abrha_vpc.foobar.id
}`

const testAccCheckAbrhaDatabaseClusterConfigMongoDB = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mongodb"
  version    = "6"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc3"
  node_count = 1
}`

const testAccCheckAbrhaDatabaseClusterConfigCustomVersion = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "%s"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc3"
  node_count = 1
}`

const testAccCheckAbrhaDatabaseClusterConfigNonDefaultProject = `
resource "abrha_project" "foobar" {
  name = "%s"
}

resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  project_id = abrha_project.foobar.id
}`
