package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	elastic6 "gopkg.in/olivere/elastic.v6"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccElasticsearchWatch(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// skip tests on ES < 6
			if v := os.Getenv("ES_VERSION"); v < "6.0.0" {
				t.Skip("Watches only supported on ES >= 6")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckElasticsearchWatchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccElasticsearchWatch,
				Check: resource.ComposeTestCheckFunc(
					testCheckElasticsearchWatchExists("elasticsearch_watch.test_watch"),
				),
			},
		},
	})
}

func testCheckElasticsearchWatchExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No watch ID is set")
		}

		meta := testAccProvider.Meta()

		var err error
		switch meta.(type) {
		case *elastic6.Client:
			client := meta.(*elastic6.Client)
			_, err = client.XPackWatchGet().Id("my_watch").Do(context.TODO())
		default:
		}

		if err != nil {
			return err
		}

		return nil
	}
}

func testCheckElasticsearchWatchDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticsearch_watch" {
			continue
		}

		meta := testAccProvider.Meta()

		var err error
		switch meta.(type) {
		case *elastic6.Client:
			client := meta.(*elastic6.Client)
			_, err = client.XPackWatchGet().Id("my_watch").Do(context.TODO())
		default:
		}

		if err != nil {
			return nil // should be not found error
		}

		return fmt.Errorf("Watch %q still exists", rs.Primary.ID)
	}

	return nil
}

var testAccElasticsearchWatch = `
resource "elasticsearch_watch" "test_watch" {
  watch_id = "my_watch"
  body = <<EOF
{
  "input": {
    "simple": {
      "payload": {
        "send": "yes"
      }
    }
  },
  "condition": {
    "always": {}
  },
  "trigger": {
    "schedule": {
      "hourly": {
        "minute": [0, 5]
      }
    }
  },
  "actions": {
    "test_index": {
      "index": {
        "index": "test",
        "doc_type": "test2"
      }
    }
  }
}
EOF
}
`
