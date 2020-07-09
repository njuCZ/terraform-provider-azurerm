package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2019-11-01-preview/spark"
	"github.com/Azure/go-autorest/autorest/to"
)

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	//log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	//
	//plugin.Serve(&plugin.ServeOpts{
	//	ProviderFunc: azurerm.Provider,
	//})

	// Resource Manager endpoints
	//endpoint := "https://management.azure.com/"
	sender := sender.BuildSender("AzureRM")

	builder := authentication.Builder{
		SubscriptionID: "85b3dbca-5974-4067-9669-67a141095a76",
		ClientID:       "1f311432-ae9f-4efb-82d6-273ffaa7889d",
		TenantID:       "72f988bf-86f1-41af-91ab-2d7cd011db47",
		ClientSecret:   "L6oo5-B5a=-NbUg8?z75j[P[4L/uwmDF",
		// we intentionally only support Client Secret auth for tests (since those variables are used all over)
		SupportsClientSecretAuth: true,
	}
	config, err := builder.Build()

	oauthConfig, err := config.BuildOAuthConfig("https://login.microsoftonline.com/")
	if err != nil {
		panic("Could not get oauth")
	}

	//auth, err := config.GetAuthorizationToken(sender, oauthConfig, "https://management.azure.com/")
	//if err != nil {
	//	panic("Could not get Authorization Token")
	//}

	synapseAuth, err := config.GetAuthorizationToken(sender, oauthConfig, "https://dev.azuresynapse.net")
	if err != nil {
		log.Printf("Could not get synapse data plane auth %+v", err)
		panic(1)
	}

	//synapseAuth, err := authrestAuth.NewAuthorizerFromEnvironmentWithResource("https://dev.azuresynapse.net")
	//if err != nil {
	//	log.Printf("Could not get synapse data plane auth %+v", err)
	//	panic(1)
	//}

	client := spark.NewBatchClient("https://cz.dev.azuresynapse.net", "cz")
	client.Authorizer = synapseAuth

	//sparkbatch, err := client.CreateSparkBatchJob(context.Background(), spark.BatchJobOptions{
	//	File:           to.StringPtr("abfss://shangwei@shangweiadlsgen2.dfs.core.windows.net/synapse/workspaces/shangweisynapseworkspace/sparkpools/sparkpool/spark-examples.jar"),
	//	ClassName:      to.StringPtr("org.apache.spark.examples.SparkPi"),
	//	ExecutorCores:  to.Int32Ptr(2),
	//	ExecutorMemory: to.StringPtr("8g"),
	//	Arguments:      to.StringSlicePtr([]string{"10"}),
	//	Name:           to.StringPtr("gosdkpi"),
	//	ExecutorCount:  to.Int32Ptr(2),
	//	DriverCores:    to.Int32Ptr(2),
	//	DriverMemory:   to.StringPtr("8g"),
	//}, nil)

	sparkbatch, err := client.GetSparkBatchJob(context.Background(), 1, to.BoolPtr(true))

	if err != nil {
		fmt.Sprintf("err: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", sparkbatch)

}
