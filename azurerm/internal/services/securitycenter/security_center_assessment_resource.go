package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSecurityCenterAssessment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterAssessmentCreateUpdate,
		Read:   resourceArmSecurityCenterAssessmentRead,
		Update: resourceArmSecurityCenterAssessmentCreateUpdate,
		Delete: resourceArmSecurityCenterAssessmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AssessmentID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"assessment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(security.CustomerManaged),
				ValidateFunc: validation.StringInSlice([]string{
					string(security.CustomerManaged),
				}, false),
			},

			"severity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(security.SeverityMedium),
				ValidateFunc: validation.StringInSlice([]string{
					string(security.SeverityLow),
					string(security.SeverityMedium),
					string(security.SeverityHigh),
				}, false),
			},

			"implementation_effort": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.ImplementationEffortLow),
					string(security.ImplementationEffortModerate),
					string(security.ImplementationEffortHigh),
				}, false),
			},

			"is_preview": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"remediation_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"threats": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"AccountBreach",
						"DataExfiltration",
						"DataSpillage",
						"MaliciousInsider",
						"ElevationOfPrivilege",
						"ThreatResistance",
						"MissingCoverage",
						"DenialOfService",
					}, false),
				},
			},

			"user_impact": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.UserImpactLow),
					string(security.UserImpactModerate),
					string(security.UserImpactHigh),
				}, false),
			},
		},
	}
}

func resourceArmSecurityCenterAssessmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAssessmentID(d.Get("target_resource_id").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.TargetResourceID, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Security Center Assessments %q : %+v", id.ID(), err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_assessment", id.ID())
		}
	}

	params := security.Assessment{
		AssessmentProperties: &security.AssessmentProperties{
			Status:

		},
	}

	if v, ok := d.GetOk("threats"); ok {
		threats := make([]security.Threats, 0)
		for _, item := range *(utils.ExpandStringSlice(v.(*schema.Set).List())) {
			threats = append(threats, (security.Threats)(item))
		}
		params.AssessmentProperties.Threats = &threats
	}

	if v, ok := d.GetOk("implementation_effort"); ok {
		params.AssessmentProperties.ImplementationEffort = security.ImplementationEffort(v.(string))
	}

	if v, ok := d.GetOk("is_preview"); ok {
		params.AssessmentProperties.Preview = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("remediation_description"); ok {
		params.AssessmentProperties.RemediationDescription = utils.String(v.(string))
	}

	if v, ok := d.GetOk("user_impact"); ok {
		params.AssessmentProperties.UserImpact = security.UserImpact(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id.TargetResourceID, id.Name, params); err != nil {
		return fmt.Errorf("creating/updating Security Center Assessment  %q : %+v", name, err)
	}

	d.SetId(id.ID())

	return resourceArmSecurityCenterAssessmentRead(d, meta)
}

func resourceArmSecurityCenterAssessmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetInSubscription(ctx, id.AssessmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] security Center Assessment  %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Security Center Assessment  %q : %+v", id.AssessmentName, err)
	}

	d.Set("name", id.AssessmentName)

	if props := resp.AssessmentProperties; props != nil {
		d.Set("assessment_type", props.AssessmentType)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("severity", props.Severity)
		d.Set("implementation_effort", props.ImplementationEffort)
		d.Set("is_preview", props.Preview)
		d.Set("remediation_description", props.RemediationDescription)
		d.Set("user_impact", props.UserImpact)

		threats := make([]string, 0)
		if props.Threats != nil {
			for _, item := range *props.Threats {
				threats = append(threats, (string)(item))
			}
		}
		d.Set("threats", utils.FlattenStringSlice(&threats))
	}

	return nil
}

func resourceArmSecurityCenterAssessmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteInSubscription(ctx, id.AssessmentName); err != nil {
		return fmt.Errorf("deleting Security Center Assessment  %q : %+v", id.AssessmentName, err)
	}

	return nil
}
