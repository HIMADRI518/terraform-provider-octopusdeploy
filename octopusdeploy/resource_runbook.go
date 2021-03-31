package octopusdeploy

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
)

func resourceRunbook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRunbookCreate,
		ReadContext:   resourceRunbookRead,
		UpdateContext: resourceRunbookUpdate,
		DeleteContext: resourceRunbookDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"project_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func buildRunbookResource(d *schema.ResourceData, m interface{}) *octopusdeploy.Runbook {
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	Runbook := octopusdeploy.NewRunbook(name, projectID)

	if attr, ok := d.GetOk("description"); ok {
		Runbook.Description = attr.(string)
	}

	return Runbook
}

func resourceRunbookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	Runbook := buildRunbookResource(d, m)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Runbooks.Add(Runbook)
	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceRunbookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.Runbooks.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", resource.Name)
	d.Set("description", resource.Description)

	return nil
}

func resourceRunbookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	Runbook := buildRunbookResource(d, m)
	Runbook.ID = d.Id() // set ID so Octopus API knows which project group to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.Runbooks.Update(Runbook)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceRunbookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Runbooks.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
