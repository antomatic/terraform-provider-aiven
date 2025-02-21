package aiven

import (
	"fmt"
	"os"
	"testing"

	"github.com/aiven/aiven-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAivenElasticsearchACLRule_basic(t *testing.T) {
	resourceName := "aiven_elasticsearch_acl_rule.foo"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAivenElasticsearchACLRuleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchACLRuleResource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project", os.Getenv("AIVEN_PROJECT_NAME")),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-aclrule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "index", "test-index"),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("user-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "permission", "readwrite"),
				),
			},
		},
	})
}

func testAccElasticsearchACLRuleResource(name string) string {
	return fmt.Sprintf(`
    data "aiven_project" "foo" {
      project = "%s"
    }

    resource "aiven_elasticsearch" "bar" {
      project = data.aiven_project.foo.project
      cloud_name = "google-europe-west1"
      plan = "startup-4"
      service_name = "test-acc-sr-aclrule-%s"
      maintenance_window_dow = "monday"
      maintenance_window_time = "10:00:00"
    }

    resource "aiven_service_user" "foo" {
      service_name = aiven_elasticsearch.bar.service_name
      project = data.aiven_project.foo.project
      username = "user-%s"
    }

    resource "aiven_elasticsearch_acl_config" "foo" {
      project = data.aiven_project.foo.project
      service_name = aiven_elasticsearch.bar.service_name
      enabled = true
      extended_acl = false
    }

    resource "aiven_elasticsearch_acl_rule" "foo" {
      project = data.aiven_project.foo.project
      service_name = aiven_elasticsearch.bar.service_name
      username = aiven_service_user.foo.username
      index = "test-index"
      permission = "readwrite"
    }
    `, os.Getenv("AIVEN_PROJECT_NAME"), name, name)
}

func testAccCheckAivenElasticsearchACLRuleResourceDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*aiven.Client)

	// loop through the resources in state, verifying each OS ACL is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aiven_opensearch_acl_rule" {
			continue
		}

		projectName, serviceName, username, index := splitResourceID4(rs.Primary.ID)

		r, err := c.ElasticsearchACLs.Get(projectName, serviceName)
		if err != nil {
			if err.(aiven.Error).Status != 404 {
				return err
			}
		}
		if r == nil {
			return nil
		}

		for _, acl := range r.ElasticSearchACLConfig.ACLs {
			if acl.Username != username {
				continue
			}
			for _, rule := range acl.Rules {
				if rule.Index == index {
					return fmt.Errorf("opensearch acl (%s) still exists", rs.Primary.ID)
				}
			}
		}
	}
	return nil
}
