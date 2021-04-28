package octopusdeploy

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
)

func expandDeployTransactJiraGateAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.TransactJiraGate"

	action.Properties["Octopus.Action.JiraGate.AuthKey"] = octopusdeploy.NewPropertyValue(flattenedAction["auth_key"].(string), false)

	if tfRequiredStatusValues, ok := flattenedAction["required_status"]; ok {
		requiredStatusValues := tfRequiredStatusValues.(map[string]interface{})

		j, _ := json.Marshal(requiredStatusValues)
		action.Properties["Octopus.Action.JiraGate.RequiredStatus"] = octopusdeploy.NewPropertyValue(string(j), false)
	}

	return action
}

func flattenDeployTransactJiraGateAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.JiraGate.AuthKey"]; ok {
		flattenedAction["auth_key"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.JiraGate.RequiredStatus"]; ok {
		var requiredStatusValues map[string]string
		json.Unmarshal([]byte(v.Value), &requiredStatusValues)

		flattenedAction["required_status"] = requiredStatusValues
	}

	return flattenedAction
}

func getDeployTransactJiraGateActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	addExecutionLocationSchema(element)
	element.Schema["auth_key"] = &schema.Schema{
		Description: "Auth Key secret",
		Required:    true,
		Type:        schema.TypeString,
	}

	element.Schema["required_status"] = &schema.Schema{
		Elem:     &schema.Schema{Type: schema.TypeString},
		Required: true,
		Type:     schema.TypeMap,
	}

	return actionSchema
}
