//go:build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/addlicense"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
