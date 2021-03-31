package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDockerContainerRegistry(d *schema.ResourceData) *octopusdeploy.DockerContainerRegistry {
	name := d.Get("name").(string)

	var dockerContainerRegistry = octopusdeploy.NewDockerContainerRegistry(name)
	dockerContainerRegistry.ID = d.Id()

	if v, ok := d.GetOk("api_version"); ok {
		dockerContainerRegistry.APIVersion = v.(string)
	}

	if v, ok := d.GetOk("feed_uri"); ok {
		dockerContainerRegistry.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("registry_path"); ok {
		dockerContainerRegistry.RegistryPath = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		dockerContainerRegistry.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		dockerContainerRegistry.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		dockerContainerRegistry.Password = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		dockerContainerRegistry.Username = v.(string)
	}

	return dockerContainerRegistry
}

func getDockerContainerRegistrySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_version": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"feed_uri": {
			Description:      "The URL to a Maven repository. This URL is the same value defined in the `repositories`/`repository`/`url` element of a Maven `settings.xml` file.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"id": getIDSchema(),
		"name": {
			Description:      "A short, memorable, unique name for this feed. Example: ACME Builds.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"password": getPasswordSchema(false),
		"package_acquisition_location_options": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"registry_path": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"username": getUsernameSchema(false),
	}
}

func setDockerContainerRegistry(ctx context.Context, d *schema.ResourceData, dockerContainerRegistry *octopusdeploy.DockerContainerRegistry) error {
	d.Set("api_version", dockerContainerRegistry.APIVersion)
	d.Set("feed_uri", dockerContainerRegistry.FeedURI)
	d.Set("name", dockerContainerRegistry.Name)
	d.Set("registry_path", dockerContainerRegistry.RegistryPath)
	d.Set("space_id", dockerContainerRegistry.SpaceID)
	d.Set("username", dockerContainerRegistry.Username)

	if err := d.Set("package_acquisition_location_options", dockerContainerRegistry.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(dockerContainerRegistry.GetID())

	return nil
}
