// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package datasource_test

import (
	"testing"

	"github.com/antimetal/terraform-provider-antimetal/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceAntimetalAWSIAMAssumeRolePolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.ProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAntimetalAWSIAMAssumeRolePolicy(),
				Check: resource.ComposeTestCheckFunc(
					testutil.TestCheckOutputExists("assume_role"),
					testutil.TestCheckOutputIsJSON("assume_role"),
				),
			},
		},
	})
}

func testAccDatasourceAntimetalAWSIAMAssumeRolePolicy() string {
	return `
		data "antimetal_aws_iam_assume_role_policy" "this" {
			external_id = "123456789"
		}

		output "assume_role" {
			value = data.antimetal_aws_iam_assume_role_policy.this.json
		}
	`
}
