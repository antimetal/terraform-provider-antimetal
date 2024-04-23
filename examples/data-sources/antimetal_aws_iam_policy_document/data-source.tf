data "antimetal_aws_iam_policy_document" "this" {}

resource "aws_iam_policy" "example" {
  name   = "example_policy"
  path   = "/"
  policy = data.antimetal_aws_iam_policy_document.this.json
}
