resource "diff-state_items" "example" {
  values = {
    "1" = "a"
    "2" = "b"
    "3" = "c"
  }

  lifecycle {
    postcondition {
      condition     = length(concat(self.new, self.updated)) <= 1
      error_message = "Created or changed more than 1 resource"
    }
  }
}