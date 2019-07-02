package main

import (
	"context"
	"fmt"
	"testing"
	"errors"

	elastic7 "github.com/olivere/elastic/v7"
	elastic6 "gopkg.in/olivere/elastic.v6"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccElasticsearchIndexLifecyclePolicy(t *testing.T) {
	provider := Provider().(*schema.Provider)
	err := provider.Configure(&terraform.ResourceConfig{})
	if err != nil {
		t.Skipf("err: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckElasticsearchIndexLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccElasticsearchIndexLifecyclePolicy,
				Check: resource.ComposeTestCheckFunc(
					testCheckElasticsearchIndexLifecyclePolicyExists("elasticsearch_index_lifecycle_policy.test"),
				),
			},
		},
	})
}

func testCheckElasticsearchIndexLifecyclePolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No index lifecycle policy ID is set")
		}

		meta := testAccProvider.Meta()

		var err error
		switch meta.(type) {
		case *elastic7.Client:
			client := meta.(*elastic7.Client)
			_, err = client.XPackIlmGetLifecycle().Policy(rs.Primary.ID).Do(context.TODO())
		case *elastic6.Client:
			client := meta.(*elastic6.Client)
			_, err = client.XPackIlmGetLifecycle().Policy(rs.Primary.ID).Do(context.TODO())
		default:
			err = errors.New("Index Lifecycle Management is only supported by the elastic library >= v6!")
		}

		if err != nil {
			return err
		}

		return nil
	}
}

func testCheckElasticsearchIndexLifecyclePolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticsearch_index_lifecycle_policy" {
			continue
		}

		meta := testAccProvider.Meta()

		var err error
		switch meta.(type) {
		case *elastic7.Client:
			client := meta.(*elastic7.Client)
			_, err = client.XPackIlmGetLifecycle().Policy(rs.Primary.ID).Do(context.TODO())
		case *elastic6.Client:
			client := meta.(*elastic6.Client)
			_, err = client.XPackIlmGetLifecycle().Policy(rs.Primary.ID).Do(context.TODO())
		default:
			err = errors.New("Index Lifecycle Management is only supported by the elastic library >= v6!")
		}

		if err != nil {
			return nil // should be not found error
		}

		return fmt.Errorf("Index lifecycle policy %q still exists", rs.Primary.ID)
	}

	return nil
}

var testAccElasticsearchIndexLifecyclePolicy = `
resource "elasticsearch_index_lifecycle_policy" "test" {
  name = "terraform-test"
  body = <<EOF
{
  "policy": {
    "phases": {
      "warm": {
        "min_age": "10d",
        "actions": {
          "forcemerge": {
            "max_num_segments": 1
          }
        }
      },
      "delete": {
        "min_age": "30d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
EOF
}
`
