// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/antimetal/terraform-provider-antimetal/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const (
	registry = "registry.terraform.io/antimetal/antimetal"
)

var (
	Version string = "dev" // overwritten on build
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run provider in debug mode")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(Version),
		providerserver.ServeOpts{
			Debug:           debug,
			Address:         registry,
			ProtocolVersion: 6,
		},
	)

	if err != nil {
		slog.Error("error starting provider", "error", err)
		os.Exit(1)
	}
}
