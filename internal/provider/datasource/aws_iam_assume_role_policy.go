// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/antimetal/terraform-provider-antimetal/internal/aws/iam"
)

const (
	antimetalIAMPrincipal = "984606353847"
)

var (
	_ datasource.DataSource = &AWSIAMAssumeRolePolicy{}
)

type AWSIAMAssumeRolePolicyModel struct {
	ExternalID types.String `tfsdk:"external_id"`
	JSON       types.String `tfsdk:"json"`
}

type AWSIAMAssumeRolePolicy struct{}

func NewAWSIAMAssumeRolePolicy() datasource.DataSource {
	return &AWSIAMAssumeRolePolicy{}
}

func (d *AWSIAMAssumeRolePolicy) Metadata(_ context.Context,
	req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_aws_iam_assume_role_policy"
}

func (d *AWSIAMAssumeRolePolicy) Schema(_ context.Context,
	_ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"external_id": schema.StringAttribute{
				Required: true,
				Description: "ID to provide secure access to your AWS account. See " +
					"https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html " +
					"for more information.",
			},
			"json": schema.StringAttribute{
				Computed: true,
				Description: "Standard JSON policy document that allows the Antimetal " +
					"IAM Principal to assume the role.",
			},
		},
	}
}

func (d *AWSIAMAssumeRolePolicy) Read(ctx context.Context,
	req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data AWSIAMAssumeRolePolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy := iam.PolicyDocument{
		Version: "2012-10-17",
		Statements: []iam.PolicyStatement{
			{
				Effect:  "Allow",
				Actions: "sts.AssumeRole",
				Principals: iam.PolicyStatementPrincipal{
					"AWS": antimetalIAMPrincipal,
				},
				Conditions: iam.PolicyStatementCondition{
					"StringEquals": map[string]string{
						"sts:ExternalID": data.ExternalID.ValueString(),
					},
				},
			},
		},
	}

	jsonPolicy, err := json.MarshalIndent(policy, "", "  ")
	// This should never happen since this document is hardcoded
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Serialize JSON for IAM Assume Role Policy Document", err.Error(),
		)
		return
	}

	data.JSON = types.StringValue(string(jsonPolicy))
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
