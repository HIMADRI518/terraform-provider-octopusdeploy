package octopusdeploy

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
)

func resourceRunbookProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRunbookProcessCreate,
		DeleteContext: resourceRunbookProcessDelete,
		Importer:      getImporter(),
		ReadContext:   resourceRunbookProcessRead,
		Schema:        getRunbookProcessSchema(),
		UpdateContext: resourceRunbookProcessUpdate,
	}
}

func resourceRunbookProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	runbookProcess := expandRunbookProcess(d)

	client := m.(*octopusdeploy.Client)
	runbook, err := client.Runbooks.GetByID(runbookProcess.RunbookID)
	if err != nil {
		diag.FromErr(err)
	}

	current, err := client.RunbookProcesses.GetByID(runbook.RunbookProcessID)
	if err != nil {
		diag.FromErr(err)
	}

	runbookProcess.ID = current.ID
	runbookProcess.Version = current.Version

	resource, err := client.RunbookProcesses.Update(*runbookProcess)
	if err != nil {
		diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceRunbookProcessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	resource, err := client.RunbookProcesses.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	logResource(constRunbookProcess, m)

	return nil
}

func resourceRunbookProcessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	runbookProcess := expandRunbookProcess(d)

	client := m.(*octopusdeploy.Client)
	current, err := client.RunbookProcesses.GetByID(runbookProcess.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	runbookProcess.Version = current.Version
	resource, err := client.RunbookProcesses.Update(*runbookProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceRunbookProcessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	current, err := client.RunbookProcesses.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	runbookProcess := &octopusdeploy.RunbookProcess{
		Version: current.Version,
	}
	runbookProcess.ID = d.Id()

	runbookProcess, err = client.RunbookProcesses.Update(*runbookProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
