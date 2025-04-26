package main

import (
	"context"
	"log"

	"terraform-provider-denvr/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go tool tfplugindocs generate --provider-dir . -provider-name denvr

func main() {
	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/denvrdata/denvr",
	}

	err := providerserver.Serve(context.Background(), provider.New(), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
