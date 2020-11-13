package octopusdeploy

import (
	"strconv"

	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getDeploymentActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addActionTypeSchema(element)
	addExecutionLocationSchema(element)
	element.Schema[constActionType] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The type of action",
		Required:    true,
	}
	addContainerSchema(element)
	addWorkerPoolSchema(element)
	addPackagesSchema(element, false)

	return actionSchema
}

func getCommonDeploymentActionSchema() (*schema.Schema, *schema.Resource) {
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the action",
				Required:    true,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Whether this step is disabled",
				Optional:    true,
				Default:     false,
			},
			constRequired: {
				Type:        schema.TypeBool,
				Description: "Whether this step is required and cannot be skipped",
				Optional:    true,
				Default:     false,
			},
			constEnvironments: {
				Description: "The environments that this step will run in",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constExcludedEnvironments: {
				Description: "The environments that this step will be skipped in",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constChannels: {
				Description: "The channels that this step applies to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_tags": {
				Description: "The tags for the tenants that this step applies to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constProperty: getPropertySchema(),
		},
	}

	actionSchema := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     element,
	}

	return actionSchema, element
}

func addExecutionLocationSchema(element *schema.Resource) {
	element.Schema["run_on_server"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Whether this step runs on a worker or on the target",
		Optional:    true,
		Default:     false,
	}
}

func addActionTypeSchema(element *schema.Resource) {
	element.Schema["action_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The type of action",
		Required:    true,
	}
}

func addWorkerPoolSchema(element *schema.Resource) {
	element.Schema["worker_pool_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Which worker pool to run on",
		Optional:    true,
	}
}

func addContainerSchema(element *schema.Resource) {
	element.Schema["container_image_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "ID of docker container image to use",
		Optional:    true,
	}
	element.Schema["container_image_feed_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "ID of docker container image to use",
		Optional:    true,
	}
}

func buildDeploymentActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := octopusdeploy.DeploymentAction{
		Channels:             getSliceFromTerraformTypeList(tfAction[constChannels]),
		Environments:         getSliceFromTerraformTypeList(tfAction[constEnvironments]),
		ExcludedEnvironments: getSliceFromTerraformTypeList(tfAction[constExcludedEnvironments]),
		IsDisabled:           tfAction["disabled"].(bool),
		IsRequired:           tfAction[constRequired].(bool),
		Name:                 tfAction["name"].(string),
		Properties:           map[string]string{},
		TenantTags:           getSliceFromTerraformTypeList(tfAction["tenant_tags"]),
	}

	actionType := tfAction[constActionType]
	if actionType != nil {
		action.ActionType = actionType.(string)
	}

	// Even though not all actions have these properties, we'll keep them here.
	// They will just be ignored if the action doesn't have it
	runOnServer := tfAction[constRunOnServer]
	if runOnServer != nil {
		action.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(runOnServer.(bool))
	}

	workerPoolID := tfAction[constWorkerPoolID]
	if workerPoolID != nil {
		action.WorkerPoolID = workerPoolID.(string)
	}

	containerImage := tfAction["container_image_id"]
	containerImageFeedId := tfAction["container_image_feed_id"]

	// if containerImage != nil {
	// 	// action.ContainerImage = containerImage.(string)
		action.Container =	map[string]string{
			"Image":  containerImage.(string),
			"FeedId":  containerImageFeedId.(string),
		}

	// }	

	if primaryPackage, ok := tfAction[constPrimaryPackage]; ok {
		tfPrimaryPackage := primaryPackage.(*schema.Set).List()
		if len(tfPrimaryPackage) > 0 {
			primaryPackage := buildPackageReferenceResource(tfPrimaryPackage[0].(map[string]interface{}))
			action.Packages = append(action.Packages, primaryPackage)
		}
	}

	if tfPkgs, ok := tfAction[constPackage]; ok {
		for _, tfPkg := range tfPkgs.(*schema.Set).List() {
			pkg := buildPackageReferenceResource(tfPkg.(map[string]interface{}))
			action.Packages = append(action.Packages, pkg)
		}
	}

	if tfProps, ok := tfAction[constProperty]; ok {
		for _, tfProp := range tfProps.(*schema.Set).List() {
			tfPropi := tfProp.(map[string]interface{})
			action.Properties[tfPropi[constKey].(string)] = tfPropi[constValue].(string)
		}
	}

	return action
}
