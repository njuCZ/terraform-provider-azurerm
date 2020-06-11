package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventHubCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubClusterExists(data.ResourceName),
				),
			},
			{
				PreConfig:         func() { time.Sleep(4 * time.Hour) }, // Cluster cannot be deleted until four hours after its creation time.
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMEventHubCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubClusterExists(data.ResourceName),
				),
			},
			{
				PreConfig:   func() { time.Sleep(4 * time.Hour) }, // Cluster cannot be deleted until four hours after its creation time.
				Config:      testAccAzureRMEventHubCluster_requiresImport(data),
				ExpectError: acceptance.RequiresImportError(data.ResourceType),
			},
		},
	})
}

func TestAccAzureRMEventHubCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubClusterExists(data.ResourceName),
				),
			},
			{
				PreConfig:         func() { time.Sleep(4 * time.Hour) }, // Cluster cannot be deleted until four hours after its creation time.
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMEventHubCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMEventHubCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubClusterExists(data.ResourceName),
				),
			},
			{
				PreConfig:         func() { time.Sleep(4 * time.Hour) }, // Cluster cannot be deleted until four hours after its creation time.
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMEventHubClusterDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.ClustersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := conn.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMEventHubClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.ClustersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := conn.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Event Hub Cluster %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("bad: Get on eventHub Cluster: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMEventHubCluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMEventHubCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_cluster" "test" {
  name                = "acctest-EC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location			  = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMEventHubCluster_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMEventHubCluster_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_cluster" "import" {
  name                = azurerm_eventhub_cluster.test.name
  resource_group_name = azurerm_eventhub_cluster.test.resource_group_name
  location			  = azurerm_eventhub_cluster.test.location
}
`, template)
}

func testAccAzureRMEventHubCluster_complete(data acceptance.TestData) string {
	template := testAccAzureRMEventHubCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_cluster" "test" {
  name                = "acctest-EC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location			  = azurerm_resource_group.test.location

  tags = {
    Env     = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMEventHubCluster_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
