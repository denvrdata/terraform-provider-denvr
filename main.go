package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-denvr/internal/provider"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/denvrdata/denvr",
	}

	err := providerserver.Serve(context.Background(), provider.New(), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
