package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KeyVaultPrivateEndpointConnectionId struct {
	ResourceGroup string
	VaultName     string
	Name          string
}

func KeyVaultPrivateEndpointConnectionID(input string) (*KeyVaultPrivateEndpointConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	conn := KeyVaultPrivateEndpointConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if conn.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}

	if conn.Name, err = id.PopSegment("privateEndpointConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &conn, nil
}
