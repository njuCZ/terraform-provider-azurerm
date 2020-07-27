package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSynapseRoleAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_role_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseRoleAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseRoleAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSynapseRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_role_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseRoleAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseRoleAssignmentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSynapseRoleAssignment_requiresImport),
		},
	})
}

func testCheckAzureRMSynapseRoleAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		synapseClient := acceptance.AzureProvider.Meta().(*clients.Client).Synapse
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("synapse role assignment not found: %s", resourceName)
		}
		id, err := parse.SynapseRoleAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := synapseClient.AccessControlClient(id.WorkspaceName)
		if resp, err := client.GetRoleAssignmentByID(ctx, id.Id); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Synapse role assignment %q does not exist", id.Id)
			}
			return fmt.Errorf("bad: Get on Synapse.AccessControlClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMSynapseRoleAssignmentDestroy(s *terraform.State) error {
	synapseClient := acceptance.AzureProvider.Meta().(*clients.Client).Synapse
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_synapse_role_assignment" {
			continue
		}
		id, err := parse.SynapseRoleAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := synapseClient.AccessControlClient(id.WorkspaceName)
		if resp, err := client.GetRoleAssignmentByID(ctx, id.Id); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Synapse.AccessControlClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMSynapseRoleAssignment_basic(data acceptance.TestData) string {
	template := testAccAzureRMSynapseRoleAssignment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_role_assignment" "test" {
  workspace_name                       = azurerm_synapse_workspace.test.name
  role_name                            = "Workspace Admin"
  principal_id                         = data.azurerm_client_config.current.object_id
}
`, template)
}

func testAccAzureRMSynapseRoleAssignment_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMSynapseRoleAssignment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_role_assignment" "import" {
  workspace_name                       = azurerm_synapse_role_assignment.test.workspace_name
  role_name                            = azurerm_synapse_role_assignment.test.role_name
  principal_id                         = azurerm_synapse_role_assignment.test.principal_id
}
`, config)
}

func testAccAzureRMSynapseRoleAssignment_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-Synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}

data "azurerm_client_config" "current" {}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
