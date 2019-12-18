package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMSpringCloudActiveDeployment_basic(t *testing.T) {
	resourceName := "azurerm_spring_cloud_active_deployment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudActiveDeployment_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "spring_cloud_deployment_name", fmt.Sprintf("acctestscd-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "persistent_disk.0.size_in_gb", "0"),
					resource.TestCheckResourceAttr(resourceName, "temporary_disk.0.size_in_gb", "5"),
					resource.TestCheckResourceAttr(resourceName, "public", "false"),
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

func TestAccAzureRMSpringCloudActiveDeployment_complete(t *testing.T) {
	resourceName := "azurerm_spring_cloud_active_deployment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudActiveDeployment_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "spring_cloud_deployment_name", fmt.Sprintf("acctestscd-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "persistent_disk.0.size_in_gb", "50"),
					resource.TestCheckResourceAttr(resourceName, "temporary_disk.0.size_in_gb", "3"),
					resource.TestCheckResourceAttr(resourceName, "public", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
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

func TestAccAzureRMSpringCloudActiveDeployment_update(t *testing.T) {
	resourceName := "azurerm_spring_cloud_active_deployment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudActiveDeployment_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "spring_cloud_deployment_name", fmt.Sprintf("acctestscd-%d", ri)),
				),
			},
			{
				Config: testAccAzureRMSpringCloudActiveDeployment_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "spring_cloud_deployment_name", fmt.Sprintf("acctestscd1-%d", ri)),
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

func testAccAzureRMSpringCloudActiveDeployment_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_active_deployment" "test" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  spring_cloud_deployment_name  = azurerm_spring_cloud_deployment.test.name
}
`, testAccAzureRMSpringCloudDeployment_basic(rInt, location))
}

func testAccAzureRMSpringCloudActiveDeployment_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_active_deployment" "test" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  spring_cloud_deployment_name  = azurerm_spring_cloud_deployment.test.name

  temporary_disk {
    mount_path = "/tmp"
    size_in_gb = 3
  }

  persistent_disk {
    mount_path = "/persistent"
    size_in_gb = 50
  }

  public = true
}
`, testAccAzureRMSpringCloudDeployment_basic(rInt, location))
}

func testAccAzureRMSpringCloudActiveDeployment_update(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_deployment" "test1" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  name= "acctestscd1-%d"

  cpu                          = 2
  memory_in_gb                 = 2
  instance_count               = 2
  runtime_version              = "Java_11"
  jvm_options                  = "-Xms1G -Xmx1G"

  env = {
      name1 = "value1"
      name2 = "value2"
  }
}

resource "azurerm_spring_cloud_active_deployment" "test" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  spring_cloud_deployment_name  = azurerm_spring_cloud_deployment.test.name
}
`, testAccAzureRMSpringCloudDeployment_basic(rInt, location), rInt)
}

func testAccAzureRMSpringCloudActiveDeployment_updated(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_deployment" "test1" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  name= "acctestscd1-%d"

  cpu                          = 2
  memory_in_gb                 = 2
  instance_count               = 2
  runtime_version              = "Java_11"
  jvm_options                  = "-Xms1G -Xmx1G"

  env = {
      name1 = "value1"
      name2 = "value2"
  }
}

resource "azurerm_spring_cloud_active_deployment" "test" {
  resource_group_name           = azurerm_resource_group.test.name
  spring_cloud_name             = azurerm_spring_cloud.test.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.test.name
  spring_cloud_deployment_name  = azurerm_spring_cloud_deployment.test1.name
}
`, testAccAzureRMSpringCloudDeployment_basic(rInt, location), rInt)
}
