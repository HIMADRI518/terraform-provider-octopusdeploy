package octopusdeploy

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
)

func expandDeployTransactJiraGateAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.TransactJiraGate"

	action.Properties["Octopus.Action.JiraGate.AuthKey"] = flattenedAction["auth_key"].(string)

	if tfRequiredStatusValues, ok := flattenedAction["required_status"]; ok {
		requiredStatusValues := make(map[string]string)

		for _, tfValue := range tfRequiredStatusValues.([]interface{}) {
			tfValueTyped := tfValue.(map[string]interface{})
			requiredStatusValues[tfValueTyped["key"].(string)] = tfValueTyped["value"].(string)
		}

		j, _ := json.Marshal(requiredStatusValues)
		action.Properties["Octopus.Action.JiraGate.RequiredStatus"] = string(j)
	}

	return action
}

func flattenDeployTransactJiraGateAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.JiraGate.AuthKey"]; ok {
		flattenedAction["auth_key"] = v
	}

	if v, ok := action.Properties["Octopus.Action.JiraGate.RequiredStatus"]; ok {
		var requiredStatusValues map[string]string
		json.Unmarshal([]byte(v), &requiredStatusValues)

		flattenedRequiredStatusValues := []interface{}{}
		for requiredStatusKey, requiredStatusValue := range requiredStatusValues {
			flattenedRequiredStatusValues = append(flattenedRequiredStatusValues, map[string]interface{}{
				"key":   requiredStatusKey,
				"value": requiredStatusValue,
			})
		}

		flattenedAction["required_status"] = flattenedRequiredStatusValues
	}

	return flattenedAction
}

func getDeployTransactJiraGateActionSchema() *schema.Schema {

	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	element.Schema["auth_key"] = &schema.Schema{
		Description: "The name of the secret resource",
		Required:    true,
		Type:        schema.TypeString,
	}

	element.Schema["required_status"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	return actionSchema
}
