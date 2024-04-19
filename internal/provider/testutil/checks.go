// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"encoding/json"
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

func TestCheckOutputExists(outputName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Outputs[outputName]
		if !ok {
			return fmt.Errorf("output %s not found", outputName)
		}

		return nil
	}
}

func TestCheckOutputIsJSON(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		o, ok := s.RootModule().Outputs[name]
		if !ok {
			return fmt.Errorf("%s not found", name)
		}

		data, ok := o.Value.(string)
		if !ok {
			return fmt.Errorf("%s output value is not a valid string", o.Value)
		}

		var js json.RawMessage
		if err := json.Unmarshal([]byte(data), &js); err != nil {
			return fmt.Errorf("couldn't unmarshal JSON data: %w", err)
		}

		return nil
	}
}
