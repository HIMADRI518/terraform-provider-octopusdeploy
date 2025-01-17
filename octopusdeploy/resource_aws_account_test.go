package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_account." + localName

	accessKey := acctest.RandString(acctest.RandIntRange(20, 255))
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secretKey := acctest.RandString(acctest.RandIntRange(20, 255))
	tenantedDeploymentParticipation := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted

	newAccessKey := acctest.RandString(acctest.RandIntRange(20, 3000))

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "access_key", accessKey),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "secret_key", secretKey),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
				),
				Config: testAWSAccountBasic(localName, name, description, accessKey, secretKey, tenantedDeploymentParticipation),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "access_key", newAccessKey),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "secret_key", secretKey),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
				),
				Config: testAWSAccountBasic(localName, name, description, newAccessKey, secretKey, tenantedDeploymentParticipation),
			},
		},
	})
}

func testAWSAccountBasic(localName string, name string, description string, accessKey string, secretKey string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_account" "%s" {
		access_key = "%s"
		description = "%s"
		name = "%s"
		secret_key = "%s"
		tenanted_deployment_participation = "%s"
	}
	
	data "octopusdeploy_accounts" "test" {
		ids = [octopusdeploy_aws_account.%s.id]
	}`, localName, accessKey, description, name, secretKey, tenantedDeploymentParticipation, localName)
}
