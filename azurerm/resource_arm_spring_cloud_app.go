package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// for some limitations of rest api:
// 1. app can not update without active deployment, rest api will throw errors
// 2. create app doesn't take care parameters, for example: TemporaryDisk
// so move these field to resource `azurerm_spring_cloud_app_active_deployment` and make SpringCloudApp unmodifiable
func resourceArmSpringCloudApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudAppCreate,
		Read:   resourceArmSpringCloudAppRead,
		Update: nil,
		Delete: resourceArmSpringCloudAppDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
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

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSpringCloudAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app", *existing.ID)
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, springCloudName, name, &appplatform.AppResource{})
	if err != nil {
		return fmt.Errorf("Error creating Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q) ID", name, springCloudName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudAppRead(d, meta)
}

func resourceArmSpringCloudAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	name := id.Path["apps"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	if resp.Properties != nil && resp.Properties.CreatedTime != nil {
		d.Set("created_time", (*resp.Properties.CreatedTime).String())
	}

	return nil
}

func resourceArmSpringCloudAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	name := id.Path["apps"]

	if _, err := client.Delete(ctx, resourceGroup, springCloudName, name); err != nil {
		return fmt.Errorf("Error deleting Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	return nil
}
