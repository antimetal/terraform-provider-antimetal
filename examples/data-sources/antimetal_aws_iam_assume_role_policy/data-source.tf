data "random_uuid" "external_id" {}

data "antimetal_aws_iam_assume_role_policy" "this" {
  external_id = random_uuid.external_id.result
}

resource "aws_iam_role" "test_role" {
  name = "test_role"

  assume_role_policy = antimetal_aws_iam_assume_role_policy.this.json

  tags = {
    tag-key = "tag-value"
  }
}
