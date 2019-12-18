package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSpringCloudApp(t *testing.T) {
	resourceName := "azurerm_spring_cloud_app.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSpringCloudAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudApp(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudAppExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "spring_cloud_name"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
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

func testCheckAzureRMSpringCloudAppExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Spring Cloud App not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		springCloudName := rs.Primary.Attributes["spring_cloud_name"]

		client := testAccProvider.Meta().(*ArmClient).AppPlatform.AppsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resGroup, springCloudName, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Spring Cloud App %q (Spring Cloud Name %q / Resource Group %q) does not exist", name, springCloudName, resGroup)
			}
			return fmt.Errorf("Bad: Get on appsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSpringCloudAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).AppPlatform.AppsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spring_cloud_app" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		springCloudName := rs.Primary.Attributes["spring_cloud_name"]

		if resp, err := client.Get(ctx, resGroup, springCloudName, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on appsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMSpringCloudApp(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app" "test" {
  spring_cloud_name             = azurerm_spring_cloud.test.name
  resource_group_name           = azurerm_resource_group.test.name
  name                          = "acctestsca-%d"
}
`, testAccAzureRMSpringCloud_basic(rInt, location), rInt)
}
