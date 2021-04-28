package octopusdeploy

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/transactcampus/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployEnvironmentBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_environment." + localName

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	resource.Test(t, resource.TestCase{
		CheckDestroy: testEnvironmentDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testEnvironmentExists(prefix),
					resource.TestCheckResourceAttr(prefix, "allow_dynamic_infrastructure", strconv.FormatBool(allowDynamicInfrastructure)),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
					resource.TestCheckResourceAttr(prefix, "use_guided_failure", strconv.FormatBool(useGuidedFailure)),
				),
				Config: testEnvironmentBasic(localName, name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure),
			},
		},
	})
}

func TestAccOctopusDeployEnvironmentMinimum(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_environment." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testEnvironmentDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testEnvironmentExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testEnvironmentMinimum(localName, name),
			},
		},
	})
}

func testEnvironmentMinimum(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}`, localName, name)
}

func testEnvironmentBasic(localName string, name string, description string, allowDynamicInfrastructure bool, sortOrder int, useGuidedFailure bool) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "%s" {
		allow_dynamic_infrastructure = "%v"
		description                  = "%s"
		name                         = "%s"
		sort_order                   = %v
		use_guided_failure           = "%v"
	}`, localName, allowDynamicInfrastructure, description, name, sortOrder, useGuidedFailure)
}

func testEnvironmentExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		environmentID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Environments.GetByID(environmentID); err != nil {
			return err
		}

		return nil
	}
}

func testEnvironmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_environment" {
			continue
		}

		environmentID := rs.Primary.ID
		environment, err := client.Environments.GetByID(environmentID)
		if err == nil {
			if environment != nil {
				return fmt.Errorf("environment (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
