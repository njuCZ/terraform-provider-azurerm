package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AssessmentsClient              *security.AssessmentsClient
	AssessmentsMetadataClient      *security.AssessmentsMetadataClient
	ContactsClient                 *security.ContactsClient
	PricingClient                  *security.PricingsClient
	WorkspaceClient                *security.WorkspaceSettingsClient
	AdvancedThreatProtectionClient *security.AdvancedThreatProtectionClient
	AutoProvisioningClient         *security.AutoProvisioningSettingsClient
	SettingClient                  *security.SettingsClient
	AutomationsClient              *security.AutomationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ascLocation := "Global"

	AssessmentsClient := security.NewAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AssessmentsClient.Client, o.ResourceManagerAuthorizer)

	AssessmentsMetadataClient := security.NewAssessmentsMetadataClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AssessmentsMetadataClient.Client, o.ResourceManagerAuthorizer)

	ContactsClient := security.NewContactsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&ContactsClient.Client, o.ResourceManagerAuthorizer)

	PricingClient := security.NewPricingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&PricingClient.Client, o.ResourceManagerAuthorizer)

	WorkspaceClient := security.NewWorkspaceSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&WorkspaceClient.Client, o.ResourceManagerAuthorizer)

	AdvancedThreatProtectionClient := security.NewAdvancedThreatProtectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AdvancedThreatProtectionClient.Client, o.ResourceManagerAuthorizer)

	AutoProvisioningClient := security.NewAutoProvisioningSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AutoProvisioningClient.Client, o.ResourceManagerAuthorizer)

	SettingClient := security.NewSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&SettingClient.Client, o.ResourceManagerAuthorizer)

	AutomationsClient := security.NewAutomationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AutomationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssessmentsClient:              &AssessmentsClient,
		AssessmentsMetadataClient:      &AssessmentsMetadataClient,
		ContactsClient:                 &ContactsClient,
		PricingClient:                  &PricingClient,
		WorkspaceClient:                &WorkspaceClient,
		AdvancedThreatProtectionClient: &AdvancedThreatProtectionClient,
		AutoProvisioningClient:         &AutoProvisioningClient,
		SettingClient:                  &SettingClient,
		AutomationsClient:              &AutomationsClient,
	}
}
