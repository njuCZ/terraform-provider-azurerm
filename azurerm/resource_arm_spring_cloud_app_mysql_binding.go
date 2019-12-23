package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudAppMysqlBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudAppMysqlBindingCreate,
		Read:   resourceArmSpringCloudAppMysqlBindingRead,
		Update: resourceArmSpringCloudAppMysqlBindingUpdate,
		Delete: resourceArmSpringCloudAppMysqlBindingDelete,

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

			"db_server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			// the return key is like: abc******
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: azure.SuppressKeyDiff,
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

func resourceArmSpringCloudAppMysqlBindingCreate(d *schema.ResourceData, meta interface{}) error {
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
			return tf.ImportAsExistsError("azurerm_spring_cloud_app_mysql_binding", *existing.ID)
		}
	}

	dbServerID := d.Get("db_server_id").(string)
	databaseName := d.Get("database_name").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	bindingParameters := map[string]interface{}{
		"databaseName": databaseName,
		"username":     username,
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: bindingParameters,
			Key:               utils.String(password),
			ResourceID:        utils.String(dbServerID),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, springCloudName, appName, name, &bindingResource); err != nil {
		return fmt.Errorf("Error creating Spring Cloud App Mysql Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Spring Cloud App Mysql Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Spring Cloud App Mysql Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q) ID", name, appName, springCloudName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudAppMysqlBindingRead(d, meta)
}

func resourceArmSpringCloudAppMysqlBindingRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] Spring Cloud App Mysql Binding %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud App Mysql Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", appName)
	if bindingResourceProperties := resp.Properties; bindingResourceProperties != nil {
		d.Set("generated_properties", bindingResourceProperties.GeneratedProperties)
		if bindingResourceProperties.GeneratedProperties != nil {
			d.Set("password", azure.FindValueInGeneratedProperties(*bindingResourceProperties.GeneratedProperties, "spring.datasource.password", "\n"))
		}
		d.Set("db_server_id", bindingResourceProperties.ResourceID)
		d.Set("resource_name", bindingResourceProperties.ResourceName)
		d.Set("resource_type", bindingResourceProperties.ResourceType)
		if databaseName, ok := bindingResourceProperties.BindingParameters["databaseName"]; ok {
			d.Set("database_name", databaseName.(string))
		}
		if username, ok := bindingResourceProperties.BindingParameters["username"]; ok {
			d.Set("username", username.(string))
		}
	}
	return nil
}

func resourceArmSpringCloudAppMysqlBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	dbServerID := d.Get("db_server_id").(string)
	databaseName := d.Get("database_name").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	bindingParameters := map[string]interface{}{
		"databaseName": databaseName,
		"username":     username,
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: bindingParameters,
			Key:               utils.String(password),
			ResourceID:        utils.String(dbServerID),
		},
	}

	if _, err := client.Update(ctx, resourceGroup, springCloudName, appName, name, &bindingResource); err != nil {
		return fmt.Errorf("Error updating Spring Cloud App Mysql Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	return resourceArmSpringCloudAppMysqlBindingRead(d, meta)
}

func resourceArmSpringCloudAppMysqlBindingDelete(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error deleting Spring Cloud App Mysql Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	return nil
}
