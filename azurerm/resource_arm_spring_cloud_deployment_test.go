package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSpringCloudDeployment_basic(t *testing.T) {
	resourceName := "azurerm_spring_cloud_deployment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSpringCloudDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudDeployment_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudDeploymentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "cpu", "1"),
					resource.TestCheckResourceAttr(resourceName, "memory_in_gb", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_version", "Java_8"),
					resource.TestCheckResourceAttr(resourceName, "env.%", "0"),
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

func TestAccAzureRMSpringCloudDeployment_complete(t *testing.T) {
	resourceName := "azurerm_spring_cloud_deployment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSpringCloudDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudDeployment_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudDeploymentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_in_gb", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "runtime_version", "Java_11"),
					resource.TestCheckResourceAttr(resourceName, "jvm_options", "-Xms1G -Xmx1G"),
					resource.TestCheckResourceAttr(resourceName, "env.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "env.name1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "env.name2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// jar_file is a input_only field, we can not get the value by read function
				ImportStateVerifyIgnore: []string{"jar_file"},
			},
		},
	})
}

func TestAccAzureRMSpringCloudDeployment_update(t *testing.T) {
	resourceName := "azurerm_spring_cloud_deployment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudDeployment_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cpu", "1"),
					resource.TestCheckResourceAttr(resourceName, "memory_in_gb", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_version", "Java_8"),
					resource.TestCheckResourceAttr(resourceName, "env.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "jvm_options", ""),
				),
			},
			{
				Config: testAccAzureRMSpringCloudDeployment_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_in_gb", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "runtime_version", "Java_11"),
					resource.TestCheckResourceAttr(resourceName, "jvm_options", "-Xms1G -Xmx1G"),
					resource.TestCheckResourceAttr(resourceName, "env.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "env.name1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "env.name2", "value2"),
				),
			},
			{
				Config: testAccAzureRMSpringCloudDeployment_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cpu", "1"),
					resource.TestCheckResourceAttr(resourceName, "memory_in_gb", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_version", "Java_8"),
					resource.TestCheckResourceAttr(resourceName, "env.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "jvm_options", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"jar_file"},
			},
		},
	})
}

func testCheckAzureRMSpringCloudDeploymentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Spring Cloud Deployment not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		springCloudName := rs.Primary.Attributes["spring_cloud_name"]
		appName := rs.Primary.Attributes["spring_cloud_app_name"]

		client := testAccProvider.Meta().(*ArmClient).AppPlatform.DeploymentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resGroup, springCloudName, appName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Spring Cloud Deployment %q (Spring Cloud Name %q / Spring Cloud App %q / Resource Group %q) does not exist", name, springCloudName, appName, resGroup)
			}
			return fmt.Errorf("Bad: Get on DeploymentsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSpringCloudDeploymentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).AppPlatform.DeploymentsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spring_cloud_deployment" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		springCloudName := rs.Primary.Attributes["spring_cloud_name"]
		appName := rs.Primary.Attributes["spring_cloud_app_name"]

		if resp, err := client.Get(ctx, resGroup, springCloudName, appName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on DeploymentsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMSpringCloudDeployment_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_deployment" "test" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  name= "acctestscd-%d"
}
`, testAccAzureRMSpringCloudApp(rInt, location), rInt)
}

func testAccAzureRMSpringCloudDeployment_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_deployment" "test" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  name= "acctestscd-%d"

  cpu                          = 2
  memory_in_gb                 = 2
  instance_count               = 2
  runtime_version              = "Java_11"
  jar_file                     = "testdata/spring_cloud_deployment_test.jar"
  jvm_options                  = "-Xms1G -Xmx1G"

  env = {
      name1 = "value1"
      name2 = "value2"
  }
}
`, testAccAzureRMSpringCloudApp(rInt, location), rInt)
}
