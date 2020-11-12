package octopusdeploy

import (
	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandVariable(d *schema.ResourceData) *octopusdeploy.Variable {
	varName := d.Get("name").(string)
	varType := d.Get(constType).(string)

	var varDesc, varValue string
	var varSensitive bool

	if varDescInterface, ok := d.GetOk("description"); ok {
		varDesc = varDescInterface.(string)
	}

	if varSensitiveInterface, ok := d.GetOk(constIsSensitive); ok {
		varSensitive = varSensitiveInterface.(bool)
	}

	if varSensitive {
		varValue = d.Get(constSensitiveValue).(string)
	} else {
		varValue = d.Get(constValue).(string)
	}

	varScopeInterface := tfVariableScopetoODVariableScope(d)

	newVar := octopusdeploy.NewVariable(varName, varType, varValue, varDesc, varScopeInterface, varSensitive)
	newVar.ID = d.Id()

	varPrompt, ok := d.GetOk(constPrompt)
	if ok {
		tfPromptSettings := varPrompt.(*schema.Set)
		if len(tfPromptSettings.List()) == 1 {
			tfPromptList := tfPromptSettings.List()[0].(map[string]interface{})
			newPrompt := octopusdeploy.VariablePromptOptions{
				Description: tfPromptList["description"].(string),
				Label:       tfPromptList["label"].(string),
				Required:    tfPromptList["is_required"].(bool),
			}
			newVar.Prompt = &newPrompt
		}
	}

	return newVar
}

func getVariableSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Required: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"type": {
			Required: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				"String",
				"Sensitive",
				"Certificate",
				"AmazonWebServicesAccount",
				"AzureAccount",
			}, false)),
		},
		"value": {
			ConflictsWith: []string{"sensitive_value"},
			Optional:      true,
			Type:          schema.TypeString,
		},
		"sensitive_value": {
			ConflictsWith: []string{"value"},
			Optional:      true,
			Sensitive:     true,
			Type:          schema.TypeString,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"scope": schemaVariableScope,
		"is_sensitive": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"prompt": {
			Type:     schema.TypeSet,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"label": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"is_required": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"pgp_key": {
			Type:      schema.TypeString,
			Optional:  true,
			ForceNew:  true,
			Sensitive: true,
		},
		"key_fingerprint": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"encrypted_value": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
