// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/antimetal/terraform-provider-antimetal/internal/provider"
)

// ProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"antimetal": providerserver.NewProtocol6WithError(provider.New("test")()),
}
