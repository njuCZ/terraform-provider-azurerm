package client

import (
	keyvaultmgmt "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	VaultsClient     *keyvault.VaultsClient
	ManagedHsmClient *keyvault.ManagedHsmsClient
	ManagementClient *keyvaultmgmt.BaseClient
}

func NewClient(o *common.ClientOptions) *Client {
	VaultsClient := keyvault.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VaultsClient.Client, o.ResourceManagerAuthorizer)

	ManagedHsmClient := keyvault.NewManagedHsmsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ManagedHsmClient.Client, o.ResourceManagerAuthorizer)

	ManagementClient := keyvaultmgmt.New()
	o.ConfigureClient(&ManagementClient.Client, o.KeyVaultAuthorizer)

	return &Client{
		VaultsClient:     &VaultsClient,
		ManagedHsmClient: &ManagedHsmClient,
		ManagementClient: &ManagementClient,
	}
}
