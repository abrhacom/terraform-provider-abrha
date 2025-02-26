package database_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaDatabaseUser_Basic(t *testing.T) {
	var databaseUser goApiAbrha.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()
	databaseUserNameUpdated := databaseUserName + "-up"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigBasic, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckNoResourceAttr(
						"abrha_database_user.foobar_user", "mysql_auth_plugin"),
					resource.TestCheckNoResourceAttr(
						"abrha_database_user.foobar_user", "access_cert"),
					resource.TestCheckNoResourceAttr(
						"abrha_database_user.foobar_user", "access_key"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigBasic, databaseClusterName, databaseUserNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserNotExists("abrha_database_user.foobar_user", databaseUserName),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserNameUpdated),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserNameUpdated),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseUser_MongoDB(t *testing.T) {
	var databaseUser goApiAbrha.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigMongo, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseUser_MongoDBMultiUser(t *testing.T) {
	databaseClusterName := acceptance.RandomTestName()
	users := []string{"foo", "bar", "baz", "one", "two"}
	config := fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigMongoMultiUser,
		databaseClusterName,
		users[0], users[0],
		users[1], users[1],
		users[2], users[2],
		users[3], users[3],
		users[4], users[4],
	)
	userResourceNames := make(map[string]string, len(users))
	for _, u := range users {
		userResourceNames[u] = fmt.Sprintf("abrha_database_user.%s", u)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						userResourceNames[users[0]], "name", users[0]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[1]], "name", users[1]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[2]], "name", users[2]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[3]], "name", users[3]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[4]], "name", users[4]),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseUser_MySQLAuth(t *testing.T) {
	var databaseUser goApiAbrha.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigMySQLAuth, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "mysql_auth_plugin", "mysql_native_password"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigMySQLAuthUpdate, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "mysql_auth_plugin", "caching_sha2_password"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigMySQLAuthRemoved, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "mysql_auth_plugin", "caching_sha2_password"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseUser_KafkaACLs(t *testing.T) {
	var databaseUser goApiAbrha.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigKafkaACL, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "access_cert"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "access_key"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "settings.0.acl.0.id"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.0.topic", "topic-1"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.0.permission", "admin"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "settings.0.acl.1.id"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.1.topic", "topic-2"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.1.permission", "produceconsume"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "settings.0.acl.2.id"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.2.topic", "topic-*"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.2.permission", "produce"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "settings.0.acl.3.id"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.3.topic", "topic-*"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.3.permission", "consume"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigKafkaACLUpdate, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "access_cert"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "access_key"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "settings.0.acl.0.id"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.0.topic", "topic-1"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.acl.0.permission", "produceconsume"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseUser_OpenSearchACLs(t *testing.T) {
	var databaseUser goApiAbrha.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigOpenSearchACL, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.0.index", "index-1"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.0.permission", "admin"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.1.index", "index-2"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.1.permission", "readwrite"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.2.index", "index-*"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.2.permission", "write"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.3.index", "index-*"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.3.permission", "read"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigOpenSearchACLUpdate, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &databaseUser),
					testAccCheckAbrhaDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.0.index", "index-1"),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "settings.0.opensearch_acl.0.permission", "readwrite"),
				),
			},
		},
	})
}

func testAccCheckAbrhaDatabaseUserDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_database_user" {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		// Try to find the database
		_, _, err := client.Databases.GetUser(context.Background(), clusterID, name)

		if err == nil {
			return fmt.Errorf("Database User still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaDatabaseUserExists(n string, databaseUser *goApiAbrha.DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database User ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		foundDatabaseUser, _, err := client.Databases.GetUser(context.Background(), clusterID, name)

		if err != nil {
			return err
		}

		if foundDatabaseUser.Name != name {
			return fmt.Errorf("Database user not found")
		}

		*databaseUser = *foundDatabaseUser

		return nil
	}
}

func testAccCheckAbrhaDatabaseUserNotExists(n string, databaseUserName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database User ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		clusterID := rs.Primary.Attributes["cluster_id"]

		_, resp, err := client.Databases.GetDB(context.Background(), clusterID, databaseUserName)

		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}

		if err == nil {
			return fmt.Errorf("Database User %s still exists", databaseUserName)
		}

		return nil
	}
}

func testAccCheckAbrhaDatabaseUserAttributes(databaseUser *goApiAbrha.DatabaseUser, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if databaseUser.Name != name {
			return fmt.Errorf("Bad name: %s", databaseUser.Name)
		}

		return nil
	}
}

const testAccCheckAbrhaDatabaseUserConfigBasic = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1

  maintenance_window {
    day  = "friday"
    hour = "13:00:00"
  }
}

resource "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckAbrhaDatabaseUserConfigMongo = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mongodb"
  version    = "6"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1

  maintenance_window {
    day  = "friday"
    hour = "13:00:00"
  }
}

resource "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckAbrhaDatabaseUserConfigMongoMultiUser = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mongodb"
  version    = "6"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_user" "%s" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}

resource "abrha_database_user" "%s" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}

resource "abrha_database_user" "%s" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}

resource "abrha_database_user" "%s" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}

resource "abrha_database_user" "%s" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckAbrhaDatabaseUserConfigMySQLAuth = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_user" "foobar_user" {
  cluster_id        = abrha_database_cluster.foobar.id
  name              = "%s"
  mysql_auth_plugin = "mysql_native_password"
}`

const testAccCheckAbrhaDatabaseUserConfigKafkaACL = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "kafka"
  version    = "3.5"
  size       = "db-s-2vcpu-2gb"
  region     = "nyc1"
  node_count = 3
}

resource "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  settings {
    acl {
      topic      = "topic-1"
      permission = "admin"
    }
    acl {
      topic      = "topic-2"
      permission = "produceconsume"
    }
    acl {
      topic      = "topic-*"
      permission = "produce"
    }
    acl {
      topic      = "topic-*"
      permission = "consume"
    }
  }
}`

const testAccCheckAbrhaDatabaseUserConfigKafkaACLUpdate = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "kafka"
  version    = "3.5"
  size       = "db-s-2vcpu-2gb"
  region     = "nyc1"
  node_count = 3
}

resource "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  settings {
    acl {
      topic      = "topic-1"
      permission = "produceconsume"
    }
  }
}`

const testAccCheckAbrhaDatabaseUserConfigOpenSearchACL = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "opensearch"
  version    = "2"
  size       = "db-s-2vcpu-4gb"
  region     = "nyc1"
  node_count = 3
}

resource "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  settings {
    opensearch_acl {
      index      = "index-1"
      permission = "admin"
    }
    opensearch_acl {
      index      = "index-2"
      permission = "readwrite"
    }
    opensearch_acl {
      index      = "index-*"
      permission = "write"
    }
    opensearch_acl {
      index      = "index-*"
      permission = "read"
    }
  }
}`

const testAccCheckAbrhaDatabaseUserConfigOpenSearchACLUpdate = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "opensearch"
  version    = "2"
  size       = "db-s-2vcpu-4gb"
  region     = "nyc1"
  node_count = 3
}

resource "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  settings {
    opensearch_acl {
      index      = "index-1"
      permission = "readwrite"
    }
  }
}`

const testAccCheckAbrhaDatabaseUserConfigMySQLAuthUpdate = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_user" "foobar_user" {
  cluster_id        = abrha_database_cluster.foobar.id
  name              = "%s"
  mysql_auth_plugin = "caching_sha2_password"
}`

const testAccCheckAbrhaDatabaseUserConfigMySQLAuthRemoved = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`
