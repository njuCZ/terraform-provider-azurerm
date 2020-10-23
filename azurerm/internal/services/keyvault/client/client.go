package client

import (
	keyvaultmgmt "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	VaultsClient                     *keyvault.VaultsClient
	PrivateEndpointConnectionsClient *keyvault.PrivateEndpointConnectionsClient
	ManagementClient                 *keyvaultmgmt.BaseClient
}

func NewClient(o *common.ClientOptions) *Client {
	VaultsClient := keyvault.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VaultsClient.Client, o.ResourceManagerAuthorizer)

	PrivateEndpointConnectionsClient := keyvault.NewPrivateEndpointConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VaultsClient.Client, o.ResourceManagerAuthorizer)

	ManagementClient := keyvaultmgmt.New()
	o.ConfigureClient(&ManagementClient.Client, o.KeyVaultAuthorizer)

	return &Client{
		VaultsClient:                     &VaultsClient,
		PrivateEndpointConnectionsClient: &PrivateEndpointConnectionsClient,
		ManagementClient:                 &ManagementClient,
	}
}
