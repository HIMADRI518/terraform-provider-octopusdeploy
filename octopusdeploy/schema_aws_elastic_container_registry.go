package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAwsElasticContainerRegistry(d *schema.ResourceData) *octopusdeploy.AwsElasticContainerRegistry {
	accessKey := d.Get("access_key").(string)
	name := d.Get("name").(string)
	secretKey := octopusdeploy.NewSensitiveValue(d.Get("secret_key").(string))
	region := d.Get("region").(string)

	var awsElasticContainerRegistry = octopusdeploy.NewAwsElasticContainerRegistry(name, accessKey, secretKey, region)
	awsElasticContainerRegistry.ID = d.Id()

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		awsElasticContainerRegistry.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("space_id"); ok {
		awsElasticContainerRegistry.SpaceID = v.(string)
	}

	return awsElasticContainerRegistry
}

func getAwsElasticContainerRegistrySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Description: "The AWS access key to use when authenticating against Amazon Web Services.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"id": {
			Computed:    true,
			Description: "The unique ID for this feed.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"name": {
			Description:      "A short, memorable, unique name for this feed. Example: ACME Builds.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"package_acquisition_location_options": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"region": {
			Description: "The AWS region where the registry resides.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"secret_key": {
			Description: "The AWS secret key to use when authenticating against Amazon Web Services.",
			Required:    true,
			Sensitive:   true,
			Type:        schema.TypeString,
		},
		"space_id": {
			Computed:    true,
			Description: "The space ID associated with this feed.",
			Optional:    true,
			Type:        schema.TypeString,
		},
	}
}

func setAwsElasticContainerRegistry(ctx context.Context, d *schema.ResourceData, awsElasticContainerRegistry *octopusdeploy.AwsElasticContainerRegistry) error {
	d.Set("access_key", awsElasticContainerRegistry.AccessKey)
	d.Set("name", awsElasticContainerRegistry.Name)
	d.Set("space_id", awsElasticContainerRegistry.SpaceID)
	d.Set("region", awsElasticContainerRegistry.Region)

	if err := d.Set("package_acquisition_location_options", awsElasticContainerRegistry.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(awsElasticContainerRegistry.GetID())

	return nil
}
