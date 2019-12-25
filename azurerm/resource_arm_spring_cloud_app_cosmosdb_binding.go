package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudAppCosmosDBBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudAppCosmosDBBindingCreate,
		Read:   resourceArmSpringCloudAppCosmosDBBindingRead,
		Update: resourceArmSpringCloudAppCosmosDBBindingUpdate,
		Delete: resourceArmSpringCloudAppCosmosDBBindingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"spring_cloud_app_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"cosmosdb_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			// because there is no cosmosdb `gremlin` resource for terraform, so currently there is only 4 types available
			"api_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cassandra",
					"mongo",
					"sql",
					"table",
				}, false),
			},

			// the return key is like: abc******
			"key": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: azure.SuppressKeyDiff,
			},

			"cassandra_key_space": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"mongo_database_name", "sql_database_name"},
			},

			"mongo_database_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"sql_database_name", "cassandra_key_space"},
			},

			"sql_database_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"mongo_database_name", "cassandra_key_space"},
			},

			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"generated_properties": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSpringCloudAppCosmosDBBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Spring Cloud App Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app_CosmosDB_binding", *existing.ID)
		}
	}

	cosmosDBAccountID := d.Get("cosmosdb_account_id").(string)
	apiType := d.Get("api_type").(string)
	key := d.Get("key").(string)

	bindingParameters := map[string]interface{}{
		"apiType": apiType,
	}
	if apiType == "sql" {
		bindingParameters["databaseName"] = d.Get("sql_database_name").(string)
	} else if apiType == "mongo" {
		bindingParameters["databaseName"] = d.Get("mongo_database_name").(string)
	} else if apiType == "cassandra" {
		bindingParameters["keySpace"] = d.Get("cassandra_key_space").(string)
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: bindingParameters,
			Key:               utils.String(key),
			ResourceID:        utils.String(cosmosDBAccountID),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, springCloudName, appName, name, &bindingResource); err != nil {
		return fmt.Errorf("Error creating Spring Cloud App CosmosDB Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Spring Cloud App CosmosDB Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Spring Cloud App CosmosDB Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q) ID", name, appName, springCloudName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudAppCosmosDBBindingRead(d, meta)
}

func resourceArmSpringCloudAppCosmosDBBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]
	name := id.Path["bindings"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App CosmosDB Binding %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud App CosmosDB Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", appName)
	if bindingResourceProperties := resp.Properties; bindingResourceProperties != nil {
		d.Set("cosmosdb_account_id", bindingResourceProperties.ResourceID)
		d.Set("resource_name", bindingResourceProperties.ResourceName)
		d.Set("resource_type", bindingResourceProperties.ResourceType)
		d.Set("generated_properties", bindingResourceProperties.GeneratedProperties)

		databaseName := ""
		keySpace := ""
		if v, ok := bindingResourceProperties.BindingParameters["databaseName"]; ok {
			databaseName = v.(string)
		}
		if v, ok := bindingResourceProperties.BindingParameters["keySpace"]; ok {
			keySpace = v.(string)
		}

		if v, ok := bindingResourceProperties.BindingParameters["apiType"]; ok {
			apiType := v.(string)
			d.Set("api_type", apiType)
			if apiType == "mongo" {
				d.Set("mongo_database_name", databaseName)
				d.Set("key", azure.FindValueInGeneratedProperties(*bindingResourceProperties.GeneratedProperties, *bindingResourceProperties.ResourceName, "@"))
			} else if apiType == "sql" {
				d.Set("sql_database_name", databaseName)
				d.Set("key", azure.FindValueInGeneratedProperties(*bindingResourceProperties.GeneratedProperties, "azure.cosmosdb.key", "\n"))
			} else if apiType == "cassandra" {
				d.Set("cassandra_key_space", keySpace)
				d.Set("key", azure.FindValueInGeneratedProperties(*bindingResourceProperties.GeneratedProperties, "spring.data.cassandra.password", "\n"))
			} else if apiType == "table" {
				d.Set("key", azure.FindValueInGeneratedProperties(*bindingResourceProperties.GeneratedProperties, "AccountKey", ";"))
			}
		}
	}
	return nil
}

func resourceArmSpringCloudAppCosmosDBBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)
	cosmosDBAccountID := d.Get("cosmosdb_account_id").(string)
	apiType := d.Get("api_type").(string)
	key := d.Get("key").(string)

	bindingParameters := map[string]interface{}{
		"apiType": apiType,
	}
	if apiType == "sql" {
		bindingParameters["databaseName"] = d.Get("sql_database_name").(string)
	} else if apiType == "mongo" {
		bindingParameters["databaseName"] = d.Get("mongo_database_name").(string)
	} else if apiType == "cassandra" {
		bindingParameters["keySpace"] = d.Get("cassandra_key_space").(string)
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: bindingParameters,
			Key:               utils.String(key),
			ResourceID:        utils.String(cosmosDBAccountID),
		},
	}

	if _, err := client.Update(ctx, resourceGroup, springCloudName, appName, name, &bindingResource); err != nil {
		return fmt.Errorf("Error updating Spring Cloud App CosmosDB Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	return resourceArmSpringCloudAppCosmosDBBindingRead(d, meta)
}

func resourceArmSpringCloudAppCosmosDBBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]
	name := id.Path["bindings"]

	if _, err := client.Delete(ctx, resourceGroup, springCloudName, appName, name); err != nil {
		return fmt.Errorf("Error deleting Spring Cloud App CosmosDB Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	return nil
}
