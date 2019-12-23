package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMSpringCloudAppMysqlBinding_basic(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_mysql_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppMysqlBinding_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttrSet(resourceName, "password"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DBforMySQL"),
					resource.TestCheckResourceAttr(resourceName, "username", "mysqladminun"),
					resource.TestCheckResourceAttr(resourceName, "database_name", fmt.Sprintf("acctest-db-%d", ri)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSpringCloudAppMysqlBinding_update(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_mysql_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppMysqlBinding_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_name", fmt.Sprintf("acctest-db-%d", ri)),
				),
			},
			{
				Config: testAccAzureRMSpringCloudAppMysqlBinding_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_name", fmt.Sprintf("acctest-db1-%d", ri)),
				),
			},
			{
				Config: testAccAzureRMSpringCloudAppMysqlBinding_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_name", fmt.Sprintf("acctest-db-%d", ri)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMSpringCloudAppMysqlBinding_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "test" {
  name                = "acctest-mysql-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 5120
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mysql_database" "test" {
  name                = "acctest-db-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = "azurerm_mysql_server.test.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}

resource "azurerm_spring_cloud_app_mysql_binding" "test" {
    name                    = "acctest-scamb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    db_server_id            = azurerm_mysql_server.test.id
    database_name           = azurerm_mysql_database.test.name
    username                = azurerm_mysql_server.test.administrator_login
    password                = azurerm_mysql_server.test.administrator_login_password
}
`, testAccAzureRMSpringCloudApp(rInt, location), rInt, rInt, rInt)
}

func testAccAzureRMSpringCloudAppMysqlBinding_update(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "test" {
  name                = "acctest-mysql-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 5120
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mysql_database" "test1" {
  name                = "acctest-db1-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = "azurerm_mysql_server.test.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}

resource "azurerm_spring_cloud_app_mysql_binding" "test" {
    name                    = "acctest-scamb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    db_server_id            = azurerm_mysql_server.test.id
    database_name           = azurerm_mysql_database.test1.name
    username                = azurerm_mysql_server.test.administrator_login
    password                = azurerm_mysql_server.test.administrator_login_password
}
`, testAccAzureRMSpringCloudApp(rInt, location), rInt, rInt)
}
