package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMSpringCloudAppCosmosDBBinding_cassandra(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_cosmosdb_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppCosmosDBBinding_cassandra(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "api_type", "cassandra"),
					resource.TestCheckResourceAttr(resourceName, "cassandra_key_space", fmt.Sprintf("acctest-ck-%d", ri)),
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

func TestAccAzureRMSpringCloudAppCosmosDBBinding_mongo(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_cosmosdb_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppCosmosDBBinding_mongo(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "api_type", "mongo"),
					resource.TestCheckResourceAttr(resourceName, "mongo_database_name", fmt.Sprintf("acctest-mongo-%d", ri)),
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

func TestAccAzureRMSpringCloudAppCosmosDBBinding_sql(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_cosmosdb_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppCosmosDBBinding_sql(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "api_type", "sql"),
					resource.TestCheckResourceAttr(resourceName, "sql_database_name", fmt.Sprintf("acctest-sql-%d", ri)),
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

func TestAccAzureRMSpringCloudAppCosmosDBBinding_table(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_cosmosdb_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppCosmosDBBinding_table(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "api_type", "table"),
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

func testAccAzureRMSpringCloudAppCosmosDBBinding_cassandra(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  capabilities {
    name = "EnableCassandra"
  }

  consistency_policy {
    consistency_level = "BoundedStaleness"
  }

  geo_location {
    location            = azurerm_resource_group.test.location
    failover_priority   = 0
  }
}

resource "azurerm_cosmosdb_cassandra_keyspace" "test" {
  name                = "acctest-ck-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_spring_cloud_app_cosmosdb_binding" "test" {
    name                    = "acctest-scacb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    cosmosdb_account_id     = azurerm_cosmosdb_account.test.id
    api_type                = "cassandra"
    cassandra_key_space     = azurerm_cosmosdb_cassandra_keyspace.test.name
    key                     = azurerm_cosmosdb_account.test.primary_master_key
}
`, testAccAzureRMSpringCloudActiveDeployment_basic(rInt, location), rInt, rInt, rInt)
}

func testAccAzureRMSpringCloudAppCosmosDBBinding_mongo(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  consistency_policy {
    consistency_level   = "BoundedStaleness"
  }

  geo_location {
    location            = azurerm_resource_group.test.location
    failover_priority   = 0
  }
}

resource "azurerm_cosmosdb_mongo_database" "test" {
  name                = "acctest-mongo-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_spring_cloud_app_cosmosdb_binding" "test" {
    name                    = "acctest-scacb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    cosmosdb_account_id     = azurerm_cosmosdb_account.test.id
    api_type                = "mongo"
    mongo_database_name     = azurerm_cosmosdb_mongo_database.test.name
    key                     = azurerm_cosmosdb_account.test.primary_master_key
}
`, testAccAzureRMSpringCloudActiveDeployment_basic(rInt, location), rInt, rInt, rInt)
}

func testAccAzureRMSpringCloudAppCosmosDBBinding_sql(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "Eventual"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-sql-%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_spring_cloud_app_cosmosdb_binding" "test" {
    name                    = "acctest-scacb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    cosmosdb_account_id     = azurerm_cosmosdb_account.test.id
    api_type                = "sql"
    sql_database_name       = azurerm_cosmosdb_sql_database.test.name
    key                     = azurerm_cosmosdb_account.test.primary_master_key
}
`, testAccAzureRMSpringCloudActiveDeployment_basic(rInt, location), rInt, rInt, rInt)
}

func testAccAzureRMSpringCloudAppCosmosDBBinding_table(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"


  consistency_policy {
    consistency_level = "BoundedStaleness"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  capabilities {
    name = "EnableTable"
  }
}

resource "azurerm_spring_cloud_app_cosmosdb_binding" "test" {
    name                    = "acctest-scacb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    cosmosdb_account_id     = azurerm_cosmosdb_account.test.id
    api_type                = "table"
    key                     = azurerm_cosmosdb_account.test.primary_master_key
}
`, testAccAzureRMSpringCloudActiveDeployment_basic(rInt, location), rInt, rInt)
}
