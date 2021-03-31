package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestTokenAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_token_account." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	token := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTokenAccountBasic(localName, description, name, tenantedDeploymentParticipation, token),
				Check: resource.ComposeTestCheckFunc(
					testAccAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(resourceName, "token", token),
				),
			},
		},
	})
}

func testTokenAccountBasic(localName string, description string, name string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode, token string) string {
	return fmt.Sprintf(`resource "octopusdeploy_token_account" "%s" {
		description                       = "%s"
		name                              = "%s"
		tenanted_deployment_participation = "%s"
		token                             = "%s"
	}`, localName, description, name, tenantedDeploymentParticipation, token)
}
