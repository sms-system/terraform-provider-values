package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/sms-system/terraform-provider-diff-state/internal/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/sms-system/diff-state",
	})

	if err != nil {
		log.Fatal(err)
	}
}
