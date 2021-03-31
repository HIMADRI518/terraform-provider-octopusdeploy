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

func TestAccOctopusDeployNuGetFeedBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_nuget_feed." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedURI := "http://test.com"
	isEnhancedMode := true
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployNuGetFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testNuGetFeedBasic(localName, name, feedURI, username, password, isEnhancedMode),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployNuGetFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "is_enhanced_mode", strconv.FormatBool(isEnhancedMode)),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
			},
		},
	})
}

func testNuGetFeedBasic(localName string, name string, feedURI string, username string, password string, isEnhancedMode bool) string {
	return fmt.Sprintf(`resource "octopusdeploy_nuget_feed" "%s" {
		feed_uri = "%s"
		is_enhanced_mode = %v
		name = "%s"
		password = "%s"
		username = "%s"
	}`, localName, feedURI, isEnhancedMode, name, password, username)
}

func testOctopusDeployNuGetFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testOctopusDeployNuGetFeedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_nuget_feed" {
			continue
		}

		_, err := client.Feeds.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
