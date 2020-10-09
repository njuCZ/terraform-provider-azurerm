package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMKustoTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoTable_basic(data),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckAzureRMKustoClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKustoTable_basic(data acceptance.TestData) string {
	return ""
}
