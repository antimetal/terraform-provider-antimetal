// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"encoding/json"

	"github.com/antimetal/terraform-provider-antimetal/internal/aws/iam"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &AWSIAMPolicyDocument{}
)

type AWSIAMPolicyDocumentModel struct {
	JSON types.String `tfsdk:"json"`
}

type AWSIAMPolicyDocument struct{}

func NewAWSIAMPolicyDocument() datasource.DataSource {
	return &AWSIAMPolicyDocument{}
}

func (d *AWSIAMPolicyDocument) Metadata(_ context.Context,
	req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_aws_iam_policy_document"
}

func (d *AWSIAMPolicyDocument) Schema(_ context.Context,
	_ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"json": schema.StringAttribute{
				Computed: true,
				Description: "Standard JSON policy document that allows Antimetal to access billing " +
					"information in your AWS account.",
			},
		},
	}
}

func (d *AWSIAMPolicyDocument) Read(ctx context.Context,
	req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data AWSIAMPolicyDocumentModel

	policy := iam.PolicyDocument{
		Version: "2012-10-17",
		Statements: []iam.PolicyStatement{
			{
				Effect:    "Allow",
				Resources: "*",
				Actions: []string{
					"application-autoscaling:Describe*",
					"autoscaling:Describe*",
					"ce:Describe*",
					"ce:Get*",
					"ce:List*",
					"cur:Get*",
					"cur:Describe*",
					"cloudwatch:GetMetricData",
					"ec2:Describe*",
					"ec2:AcceptReservedInstancesExchangeQuote",
					"ec2:CancelReservedInstancesListing",
					"ec2:CreateReservedInstancesListing",
					"ec2:DeleteQueuedReservedInstances",
					"ec2:ModifyReservedInstances",
					"ec2:PurchaseHostReservation",
					"ec2:PurchaseReservedInstancesOffering",
					"elasticache:List*",
					"elasticache:Describe*",
					"elasticache:PurchaseReservedCacheNodesOffering",
					"es:Describe*",
					"es:List*",
					"es:PurchaseReservedInstanceOffering",
					"organizations:InviteAccountToOrganization",
					"organizations:List*",
					"organizations:Describe*",
					"organizations:AcceptHandshake",
					"iam:CreateServiceLinkedRole",
					"pricing:DescribeServices",
					"pricing:GetAttributeValues",
					"pricing:GetProducts",
					"rds:Describe*",
					"rds:List*",
					"rds:PurchaseReservedDbInstancesOffering",
					"savingsplans:Describe*",
					"savingsplans:List*",
					"savingsplans:CreateSavingsPlan",
					"servicequotas:Get*",
					"servicequotas:List*",
					"sagemaker:Describe*",
					"sagemaker:List*",
					"medialive:Describe*",
					"medialive:List*",
					"medialive:PurchaseOffering",
					"redshift:Describe*",
					"redshift:List*",
					"redshift:PurchaseReservedNodeOffering",
					"support:*",
				},
			},
		},
	}

	jsonPolicy, err := json.MarshalIndent(policy, "", "  ")
	// This should never happen since this is the document is hardcoded
	if err != nil {
		resp.Diagnostics.AddError("Failed to Serialize JSON for IAM Policy Document", err.Error())
		return
	}

	data.JSON = types.StringValue(string(jsonPolicy))
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
