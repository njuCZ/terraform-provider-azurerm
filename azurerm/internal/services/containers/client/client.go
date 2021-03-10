package client

import (
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2019-12-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	legacy "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-08-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-12-01/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AgentPoolsClient         *containerservice.AgentPoolsClient
	GroupsClient             *containerinstance.ContainerGroupsClient
	RegistriesClient         *containerregistry.RegistriesClient
	ReplicationsClient       *containerregistry.ReplicationsClient
	ServicesClient           *legacy.ContainerServicesClient
	WebhooksClient           *containerregistry.WebhooksClient

	Environment azure.Environment

	o *common.ClientOptions
}

func NewClient(o *common.ClientOptions) *Client {
	registriesClient := containerregistry.NewRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&registriesClient.Client, o.ResourceManagerAuthorizer)

	webhooksClient := containerregistry.NewWebhooksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webhooksClient.Client, o.ResourceManagerAuthorizer)

	replicationsClient := containerregistry.NewReplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&replicationsClient.Client, o.ResourceManagerAuthorizer)

	groupsClient := containerinstance.NewContainerGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)

	// AKS
	agentPoolsClient := containerservice.NewAgentPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&agentPoolsClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := legacy.NewContainerServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AgentPoolsClient:         &agentPoolsClient,
		GroupsClient:             &groupsClient,
		RegistriesClient:         &registriesClient,
		WebhooksClient:           &webhooksClient,
		ReplicationsClient:       &replicationsClient,
		ServicesClient:           &servicesClient,
		Environment:              o.Environment,
		o:o,
	}
}

func (client Client) NewKubernetesClustersClient(headers map[string]interface{}) *containerservice.ManagedClustersClient {
	kubernetesClustersClient := containerservice.NewManagedClustersClientWithBaseURI(client.o.ResourceManagerEndpoint, client.o.SubscriptionId)
	client.o.ConfigureClient(&kubernetesClustersClient.Client, client.o.ResourceManagerAuthorizer)

	if len(headers) > 0 {
		decorate := kubernetesClustersClient.Client.RequestInspector
		if decorate == nil {
			decorate = autorest.WithNothing()
		}
		kubernetesClustersClient.Client.RequestInspector = func(p autorest.Preparer) autorest.Preparer {
			return autorest.DecoratePreparer(p, decorate, autorest.WithHeaders(headers))
		}
	}

	return &kubernetesClustersClient
}
