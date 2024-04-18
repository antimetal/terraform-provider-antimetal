// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("id not set for resource %s", resourceName)
		}

		return nil
	}
}

