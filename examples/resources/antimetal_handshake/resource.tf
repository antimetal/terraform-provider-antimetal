resource "random_uuid" "external_id" {}

resource "antimetal_handshake" "this" {
  // Visit https://docs.antimetal.com/onboarding/overview to receive handshake_id
  handshake_id = "7ef4badd-6946-455a-8788-12c479b0858c"
  external_id  = random_uuid.external_id.result
  role_arn     = "arn:aws:iam::012345678999:role/example_role"
}
