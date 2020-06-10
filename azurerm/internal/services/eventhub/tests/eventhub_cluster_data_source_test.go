package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMEventHubCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMEventHubCluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMEventHubCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_cluster" "test" {
  name                = azurerm_eventhub_cluster.test.name
  resource_group_name = azurerm_eventhub_cluster.test.resource_group_name
}
`, template)
}
