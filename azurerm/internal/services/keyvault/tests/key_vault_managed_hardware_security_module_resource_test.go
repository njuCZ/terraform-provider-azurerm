package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultManagedHardwareSecurityModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultManagedHardwareSecurityModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultManagedHardwareSecurityModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultManagedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultManagedHardwareSecurityModule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultManagedHardwareSecurityModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultManagedHardwareSecurityModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultManagedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMKeyVaultManagedHardwareSecurityModule_requiresImport),
		},
	})
}

func TestAccAzureRMKeyVaultManagedHardwareSecurityModule_updateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultManagedHardwareSecurityModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultManagedHardwareSecurityModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultManagedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKeyVaultManagedHardwareSecurityModule_updateTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultManagedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMKeyVaultManagedHardwareSecurityModuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagedHsmClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_managed_hardware_security_module" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Key Vault Managed Hardware Security Module still exists:\n%#v", resp.Properties)
	}

	return nil
}

func testCheckAzureRMKeyVaultManagedHardwareSecurityModuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagedHsmClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		vaultName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, vaultName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Key Vault Managed Hardware Security Module %q (resource group: %q) does not exist", vaultName, resourceGroup)
			}

			return fmt.Errorf("bad: Get on ManagedHsmClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultManagedHardwareSecurityModule_basic(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultManagedHardwareSecurityModule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  purge_protection_enabled   = false
  soft_delete_retention_days = 7

  tenant_id        = data.azurerm_client_config.current.tenant_id
  admin_object_ids = [data.azurerm_client_config.current.object_id]
}
`, template, data.RandomInteger)
}

func testAccAzureRMKeyVaultManagedHardwareSecurityModule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultManagedHardwareSecurityModule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_managed_hardware_security_module" "import" {
  name                       = azurerm_key_vault_managed_hardware_security_module.test.name
  resource_group_name        = azurerm_key_vault_managed_hardware_security_module.test.resource_group_name
  location                   = azurerm_key_vault_managed_hardware_security_module.test.location
  sku_name                   = azurerm_key_vault_managed_hardware_security_module.test.sku_name
  purge_protection_enabled   = azurerm_key_vault_managed_hardware_security_module.test.purge_protection_enabled
  soft_delete_retention_days = azurerm_key_vault_managed_hardware_security_module.test.soft_delete_retention_days

  tenant_id        = azurerm_key_vault_managed_hardware_security_module.test.tenant_id
  admin_object_ids = azurerm_key_vault_managed_hardware_security_module.test.admin_object_ids
}
`, template)
}

func testAccAzureRMKeyVaultManagedHardwareSecurityModule_updateTags(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultManagedHardwareSecurityModule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  purge_protection_enabled   = false
  soft_delete_retention_days = 7

  tenant_id        = data.azurerm_client_config.current.tenant_id
  admin_object_ids = [data.azurerm_client_config.current.object_id]

  tags = {
    Env = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMKeyVaultManagedHardwareSecurityModule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-KV-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
