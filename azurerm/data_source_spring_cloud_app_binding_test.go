package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_cassandra(t *testing.T) {
	resourceName := "data.azurerm_spring_cloud_app_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_cassandra(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.apiType", "cassandra"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.keySpace", fmt.Sprintf("acctest-ck-%d", ri)),
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

func TestAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_mongo(t *testing.T) {
	resourceName := "data.azurerm_spring_cloud_app_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_mongo(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.apiType", "mongo"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.databaseName", fmt.Sprintf("acctest-mongo-%d", ri)),
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

func TestAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_sql(t *testing.T) {
	resourceName := "data.azurerm_spring_cloud_app_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_sql(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.apiType", "sql"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.databaseName", fmt.Sprintf("acctest-sql-%d", ri)),
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

func TestAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_table(t *testing.T) {
	resourceName := "data.azurerm_spring_cloud_app_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_table(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DocumentDB"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.apiType", "table"),
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

func TestAccDataSourceAzureRMSpringCloudAppRedisBinding(t *testing.T) {
	resourceName := "data.azurerm_spring_cloud_app_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSpringCloudAppRedisBinding(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.Cache"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.useSsl", "true"),
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

func TestAccDataSourceAzureRMSpringCloudAppMysqlBinding(t *testing.T) {
	resourceName := "data.azurerm_spring_cloud_app_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSpringCloudAppMysqlBinding(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.DBforMySQL"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.databaseName", "exampledb"),
					resource.TestCheckResourceAttr(resourceName, "binding_parameters.username", "mysqladminun"),
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

func testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_cassandra(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app_binding" "test" {
    name                    = azurerm_spring_cloud_app_cosmosdb_binding.test.name
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
}
`, testAccAzureRMSpringCloudAppCosmosDBBinding_cassandra(rInt, location))
}

func testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_mongo(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app_binding" "test" {
    name                    = azurerm_spring_cloud_app_cosmosdb_binding.test.name
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
}
`, testAccAzureRMSpringCloudAppCosmosDBBinding_mongo(rInt, location))
}

func testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_sql(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app_binding" "test" {
    name                    = azurerm_spring_cloud_app_cosmosdb_binding.test.name
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
}
`, testAccAzureRMSpringCloudAppCosmosDBBinding_sql(rInt, location))
}

func testAccDataSourceAzureRMSpringCloudAppCosmosDBBinding_table(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app_binding" "test" {
    name                    = azurerm_spring_cloud_app_cosmosdb_binding.test.name
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
}
`, testAccAzureRMSpringCloudAppCosmosDBBinding_table(rInt, location))
}

func testAccDataSourceAzureRMSpringCloudAppRedisBinding(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app_binding" "test" {
    name                    = azurerm_spring_cloud_app_redis_binding.test.name
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
}
`, testAccAzureRMSpringCloudAppRedisBinding_basic(rInt, location))
}

func testAccDataSourceAzureRMSpringCloudAppMysqlBinding(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app_binding" "test" {
    name                    = azurerm_spring_cloud_app_mysql_binding.test.name
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
}
`, testAccAzureRMSpringCloudAppMysqlBinding_basic(rInt, location))
}
