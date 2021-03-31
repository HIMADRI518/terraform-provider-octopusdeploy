package octopusdeploy

import (
	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandRunbookProcess(d *schema.ResourceData) *octopusdeploy.RunbookProcess {
	runbookProcess := octopusdeploy.NewRunbookProcess(d.Get("runbook_id").(string), d.Get("project_id").(string))
	runbookProcess.ID = d.Id()

	if v, ok := d.GetOk("step"); ok {
		steps := v.([]interface{})
		for _, step := range steps {
			runbookStep := expandDeploymentStep(step.(map[string]interface{}))
			runbookProcess.Steps = append(runbookProcess.Steps, runbookStep)
		}
	}

	return runbookProcess
}

func getRunbookProcessSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"last_snapshot_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"runbook_id": {
			Required: true,
			Type:     schema.TypeString,
		},
		"project_id": {
			Required: true,
			Type:     schema.TypeString,
		},
		"step": getDeploymentStepSchema(),
		"version": {
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
