resource "values_diff" "example" {
  values = {
    "1" = "a"
    "2" = "b"
    "3" = "c"
  }

  lifecycle {
    postcondition {
      condition     = self.is_initiated || length(concat(self.created, self.updated)) <= 1
      error_message = "Created or updated more than 1 item"
    }
  }
}