package octopusdeploy

import (
	"context"

	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLifecycle() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLifecycleCreate,
		DeleteContext: resourceLifecycleDelete,
		Importer:      getImporter(),
		ReadContext:   resourceLifecycleRead,
		Schema:        getLifecycleSchema(),
		UpdateContext: resourceLifecycleUpdate,
	}
}

func resourceLifecycleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lifecycle := expandLifecycle(d)

	client := m.(*octopusdeploy.Client)
	createdLifecycle, err := client.Lifecycles.Add(lifecycle)
	if err != nil {
		return diag.FromErr(err)
	}

	setLifecycle(ctx, d, createdLifecycle)
	return nil
}

func resourceLifecycleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	lifecycle, err := client.Lifecycles.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setLifecycle(ctx, d, lifecycle)
	return nil
}

func resourceLifecycleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lifecycle := expandLifecycle(d)

	client := m.(*octopusdeploy.Client)
	updatedLifecycle, err := client.Lifecycles.Update(lifecycle)
	if err != nil {
		return diag.FromErr(err)
	}

	setLifecycle(ctx, d, updatedLifecycle)
	return nil
}

func resourceLifecycleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Lifecycles.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
