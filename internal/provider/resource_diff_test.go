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
				Config: `resource "values_diff" "test" {
					values = {
						"1" = "a"
						"2" = "b"
						"3" = "c"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("values_diff.test", "id", "diff"),

					resource.TestCheckResourceAttr("values_diff.test", "is_initiated", "true"),

					resource.TestCheckResourceAttr("values_diff.test", "last_values.%", "0"),

					resource.TestCheckResourceAttr("values_diff.test", "created.#", "3"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "created.*", "1"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "created.*", "2"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "created.*", "3"),

					resource.TestCheckResourceAttr("values_diff.test", "updated.#", "0"),

					resource.TestCheckResourceAttr("values_diff.test", "deleted.#", "0"),

					resource.TestCheckResourceAttr("values_diff.test", "is_value_commited", "true"),
				),
			},
			{
				Config: `resource "values_diff" "test" {
					values = {
						"1" = "a"
						"3" = "cc"
						"4" = "d"
						"5" = "e"
					}

					commit_exp = "[...created, ...updated].length <= 1"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("values_diff.test", "is_initiated", "false"),

					resource.TestCheckResourceAttr("values_diff.test", "last_values.%", "3"),
					resource.TestCheckResourceAttr("values_diff.test", "last_values.1", "a"),
					resource.TestCheckResourceAttr("values_diff.test", "last_values.2", "b"),
					resource.TestCheckResourceAttr("values_diff.test", "last_values.3", "c"),

					resource.TestCheckResourceAttr("values_diff.test", "created.#", "2"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "created.*", "4"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "created.*", "5"),

					resource.TestCheckResourceAttr("values_diff.test", "updated.#", "1"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "updated.*", "3"),

					resource.TestCheckResourceAttr("values_diff.test", "deleted.#", "1"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "deleted.*", "2"),

					resource.TestCheckResourceAttr("values_diff.test", "is_value_commited", "false"),
				),
			},
			{
				Config: `resource "values_diff" "test" {
					values = {
						"1" = "a"
					}

					commit_exp = "[...created, ...updated].length === 0"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("values_diff.test", "is_initiated", "false"),

					resource.TestCheckResourceAttr("values_diff.test", "last_values.%", "3"),
					resource.TestCheckResourceAttr("values_diff.test", "last_values.1", "a"),
					resource.TestCheckResourceAttr("values_diff.test", "last_values.2", "b"),
					resource.TestCheckResourceAttr("values_diff.test", "last_values.3", "c"),

					resource.TestCheckResourceAttr("values_diff.test", "created.#", "0"),

					resource.TestCheckResourceAttr("values_diff.test", "updated.#", "0"),

					resource.TestCheckResourceAttr("values_diff.test", "deleted.#", "2"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "deleted.*", "2"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "deleted.*", "3"),

					resource.TestCheckResourceAttr("values_diff.test", "is_value_commited", "true"),
				),
			},
			{
				Config: `resource "values_diff" "test" {
					values = {
						"1" = "b"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("values_diff.test", "is_initiated", "false"),

					resource.TestCheckResourceAttr("values_diff.test", "last_values.%", "1"),
					resource.TestCheckResourceAttr("values_diff.test", "last_values.1", "a"),

					resource.TestCheckResourceAttr("values_diff.test", "created.#", "0"),

					resource.TestCheckResourceAttr("values_diff.test", "updated.#", "1"),
					resource.TestCheckTypeSetElemAttr("values_diff.test", "updated.*", "1"),

					resource.TestCheckResourceAttr("values_diff.test", "deleted.#", "0"),

					resource.TestCheckResourceAttr("values_diff.test", "is_value_commited", "true"),
				),
			},
		},
	})
}
