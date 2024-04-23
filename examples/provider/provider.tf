terraform {
  required_providers {
    antimetal = {
      source = "antimetal/antimetal"
    }

    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "antimetal" {}

provider "aws" {
  region = "us-east-1"
}

data "random_uuid" "external_id" {}

resource "antimetal_handshake" "this" {
  // Visit https://docs.antimetal.com/onboarding/overview to receive handshake_id
  handshake_id = "7ef4badd-6946-455a-8788-12c479b0858c"
  external_id  = random_uuid.external_id.result
  role_arn     = aws_iam_role.test_role.arn
}

data "antimetal_aws_iam_assume_role_policy" "this" {
  external_id = random_uuid.external_id.result
}

resource "aws_iam_role" "test_role" {
  name = "test_role"

  assume_role_policy  = antimetal_aws_iam_assume_role_policy.this.json
  managed_policy_arns = [aws_iam_policy.antimetal_billing_policy.arn]

  tags = {
    tag-key = "tag-value"
  }
}

data "antimetal_aws_iam_policy_document" "this" {}

resource "aws_iam_policy" "antimetal_billing_policy" {
  name   = "antimetal_billing_policy"
  path   = "/"
  policy = data.antimetal_aws_iam_policy_document.this.json
}
