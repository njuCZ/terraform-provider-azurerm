package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMSpringCloudAppRedisBinding_basic(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_redis_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppRedisBinding_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "generated_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_name"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "Microsoft.Cache"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "true"),
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

func TestAccAzureRMSpringCloudAppRedisBinding_update(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app_redis_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudAppRedisBinding_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMSpringCloudAppRedisBinding_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "false"),
				),
			},
			{
				Config: testAccAzureRMSpringCloudAppRedisBinding_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "true"),
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

func testAccAzureRMSpringCloudAppRedisBinding_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_cache" "test" {
  name                     = "acctestredis-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  capacity           	   = 0
  family                   = "C"
  sku_name                 = "Basic"
  enable_non_ssl_port      = true
}

resource "azurerm_spring_cloud_app_redis_binding" "test" {
	name 					= "acctestscarb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    redis_cache_id          = azurerm_redis_cache.test.id
    key                     = azurerm_redis_cache.test.primary_access_key
}
`, testAccAzureRMSpringCloudApp(rInt, location), rInt, rInt)
}

func testAccAzureRMSpringCloudAppRedisBinding_update(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_cache" "test" {
  name                     = "acctestredis-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  capacity           	   = 0
  family                   = "C"
  sku_name                 = "Basic"
  enable_non_ssl_port      = true
}

resource "azurerm_spring_cloud_app_redis_binding" "test" {
	name 					= "acctestscarb-%d"
    resource_group_name    	= azurerm_resource_group.test.name
    spring_cloud_name       = azurerm_spring_cloud.test.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.test.name
    redis_cache_id          = azurerm_redis_cache.test.id
    key                     = azurerm_redis_cache.test.secondary_access_key
    ssl_enabled             = false
}
`, testAccAzureRMSpringCloudApp(rInt, location), rInt, rInt)
}
