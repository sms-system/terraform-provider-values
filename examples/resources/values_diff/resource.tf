resource "values_diff" "example" {
  values = {
    "1" = "a"
    "2" = "b"
    "3" = "c"
  }

  commit_exp = <<-EOT
    is_initiated || [...created, ...updated].length <= 1
  EOT

  lifecycle {
    postcondition {
      condition     = self.is_value_commited
      error_message = "Created or updated more than 1 item"
    }
  }
}