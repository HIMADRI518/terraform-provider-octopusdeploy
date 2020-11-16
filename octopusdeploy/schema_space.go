package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
)

func expandSpace(d *schema.ResourceData) *octopusdeploy.Space {
	name := d.Get("name").(string)

	space := octopusdeploy.NewSpace(name)
	space.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		space.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		space.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("space_managers_team_members"); ok {
		space.SpaceManagersTeamMembers = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("space_managers_teams"); ok {
		space.SpaceManagersTeams = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("task_queue_stopped"); ok {
		space.TaskQueueStopped = v.(bool)
	}

	return space
}

func flattenSpace(ctx context.Context, d *schema.ResourceData, space *octopusdeploy.Space) {
	d.Set("description", space.Description)
	d.Set("is_default", space.IsDefault)
	d.Set("name", space.Name)
	d.Set("space_managers_team_members", space.SpaceManagersTeamMembers)
	d.Set("space_managers_teams", space.SpaceManagersTeams)
	d.Set("task_queue_stopped", space.TaskQueueStopped)

	d.SetId(space.GetID())
}

func getSpaceDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ids": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"partial_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"spaces": {
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"is_default": {
						Optional: true,
						Type:     schema.TypeBool,
					},
					"id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"name": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"space_managers_team_members": {
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
						Type:     schema.TypeList,
					},
					"space_managers_teams": {
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
						Type:     schema.TypeList,
					},
					"task_queue_stopped": {
						Optional: true,
						Type:     schema.TypeBool,
					},
				},
			},
			Type: schema.TypeList,
		},
	}
}

func getSpaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"is_default": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"name": {
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"space_managers_team_members": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"space_managers_teams": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"task_queue_stopped": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
