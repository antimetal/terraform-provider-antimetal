// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package datasource_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/antimetal/terraform-provider-antimetal/internal/provider/testutil"
)

func TestAccDatasourceAntimetalAWSIAMPolicyDocument_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.ProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAntimetalAWSIAMPolicyDocument(),
				Check: resource.ComposeTestCheckFunc(
					testutil.TestCheckOutputExists("policy"),
					testutil.TestCheckOutputIsJSON("policy"),
				),
			},
		},
	})
}

func testAccDatasourceAntimetalAWSIAMPolicyDocument() string {
	return `
		data "antimetal_aws_iam_policy_document" "this" {}

		output "policy" {
			value = data.antimetal_aws_iam_policy_document.this.json
		}
	`
}
