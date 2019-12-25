package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSpringCloudAppBinding() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSpringCloudAppBindingRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

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

			"binding_parameters": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"generated_properties": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSpringCloudAppBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Spring Cloud App Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q) was not found", name, appName, springCloudName, resourceGroup)
		}
		return fmt.Errorf("Error reading Spring Cloud App Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", appName)
	if bindingResourceProperties := resp.Properties; bindingResourceProperties != nil {
		d.Set("binding_parameters", bindingResourceProperties.BindingParameters)
		d.Set("generated_properties", bindingResourceProperties.GeneratedProperties)
		d.Set("key", bindingResourceProperties.Key)
		d.Set("resource_id", bindingResourceProperties.ResourceID)
		d.Set("resource_name", bindingResourceProperties.ResourceName)
		d.Set("resource_type", bindingResourceProperties.ResourceType)
	}

	return nil
}
