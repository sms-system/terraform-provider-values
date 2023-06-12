package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDiffStateItemsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `resource "diff-state_items" "test" {
					values = {
						"1" = "a"
						"2" = "b"
						"3" = "c"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("diff-state_items.test", "id", "diff"),

					resource.TestCheckResourceAttr("diff-state_items.test", "is_initiated", "true"),

					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.%", "0"),

					resource.TestCheckResourceAttr("diff-state_items.test", "created.#", "3"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "created.*", "1"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "created.*", "2"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "created.*", "3"),

					resource.TestCheckResourceAttr("diff-state_items.test", "updated.#", "0"),

					resource.TestCheckResourceAttr("diff-state_items.test", "deleted.#", "0"),

					resource.TestCheckResourceAttr("diff-state_items.test", "is_value_commited", "true"),
				),
			},
			{
				Config: `resource "diff-state_items" "test" {
					values = {
						"1" = "a"
						"3" = "cc"
						"4" = "d"
						"5" = "e"
					}

					commit_exp = "[...created, ...updated].length <= 1"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("diff-state_items.test", "is_initiated", "false"),

					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.%", "3"),
					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.1", "a"),
					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.2", "b"),
					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.3", "c"),

					resource.TestCheckResourceAttr("diff-state_items.test", "created.#", "2"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "created.*", "4"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "created.*", "5"),

					resource.TestCheckResourceAttr("diff-state_items.test", "updated.#", "1"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "updated.*", "3"),

					resource.TestCheckResourceAttr("diff-state_items.test", "deleted.#", "1"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "deleted.*", "2"),

					resource.TestCheckResourceAttr("diff-state_items.test", "is_value_commited", "false"),
				),
			},
			{
				Config: `resource "diff-state_items" "test" {
					values = {
						"1" = "a"
					}

					commit_exp = "[...created, ...updated].length === 0"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("diff-state_items.test", "is_initiated", "false"),

					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.%", "3"),
					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.1", "a"),
					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.2", "b"),
					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.3", "c"),

					resource.TestCheckResourceAttr("diff-state_items.test", "created.#", "0"),

					resource.TestCheckResourceAttr("diff-state_items.test", "updated.#", "0"),

					resource.TestCheckResourceAttr("diff-state_items.test", "deleted.#", "2"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "deleted.*", "2"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "deleted.*", "3"),

					resource.TestCheckResourceAttr("diff-state_items.test", "is_value_commited", "true"),
				),
			},
			{
				Config: `resource "diff-state_items" "test" {
					values = {
						"1" = "b"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("diff-state_items.test", "is_initiated", "false"),

					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.%", "1"),
					resource.TestCheckResourceAttr("diff-state_items.test", "last_values.1", "a"),

					resource.TestCheckResourceAttr("diff-state_items.test", "created.#", "0"),

					resource.TestCheckResourceAttr("diff-state_items.test", "updated.#", "1"),
					resource.TestCheckTypeSetElemAttr("diff-state_items.test", "updated.*", "1"),

					resource.TestCheckResourceAttr("diff-state_items.test", "deleted.#", "0"),

					resource.TestCheckResourceAttr("diff-state_items.test", "is_value_commited", "true"),
				),
			},
		},
	})
}
