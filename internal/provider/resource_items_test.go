package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDiffStateItemsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDiffStateItemsResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("diff-state_items.test", "id",         "diff"),
					resource.TestCheckResourceAttr("diff-state_items.test", "previous.%", "0"),
					resource.TestCheckResourceAttr("diff-state_items.test", "new.#",      "0"),
					resource.TestCheckResourceAttr("diff-state_items.test", "updated.#",  "0"),
					resource.TestCheckResourceAttr("diff-state_items.test", "deleted.#",  "0"),
				),
			},
			{
				Config: testAccDiffStateItemsModifiedResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("diff-state_items.test", "previous.%", "3"),
					resource.TestCheckResourceAttr("diff-state_items.test", "previous.1", "a"),
					resource.TestCheckResourceAttr("diff-state_items.test", "previous.2", "b"),
					resource.TestCheckResourceAttr("diff-state_items.test", "previous.3", "c"),

					resource.TestCheckResourceAttr("diff-state_items.test", "new.#",      "2"),
					resource.TestCheckResourceAttr("diff-state_items.test", "new.0",      "4"),
					resource.TestCheckResourceAttr("diff-state_items.test", "new.1",      "5"),
					
					resource.TestCheckResourceAttr("diff-state_items.test", "updated.#",  "1"),
					resource.TestCheckResourceAttr("diff-state_items.test", "updated.0",  "3"),

					resource.TestCheckResourceAttr("diff-state_items.test", "deleted.#",  "1"),
					resource.TestCheckResourceAttr("diff-state_items.test", "deleted.0",  "2"),
				),
			},
		},
	})
}

func testAccDiffStateItemsResourceConfig() string {
	return fmt.Sprintf(`
		resource "diff-state_items" "test" {
			values = {
				"1" = "a"
				"2" = "b"
				"3" = "c"
			}
		}
	`)
}

func testAccDiffStateItemsModifiedResourceConfig() string {
	return fmt.Sprintf(`
		resource "diff-state_items" "test" {
			values = {
				"1" = "a"
				"3" = "cс"
				"4" = "d"
				"5" = "e"
			}
		}
	`)
}