# terraform-provider-antimetal

The Antimetal provider lets you connect your AWS account to Antimetal using Terraform.
This is done by creating a cross-account IAM role in your AWS account that Antimetal can assume to
access your billing data.

Visit https://antimetal.com/ for more information.

See the documentation in the Terraform registry for the most up-to-date information
and latest release.

This provider is maintained by Antimetal.

## Requirements
- Bash
- Make
- [Go](https://go.dev/doc/install) 1.22
- [Terraform CLI](https://developer.hashicorp.com/terraform/install) 1.x

## Developing

To see the full list of make targets, run:

```
make help
```

### Build

To build the provider, run:

```
make build // Or simply `make`
```

This will build the provider binary under `bin/`.

### Test

To run unit tests, run:

```
make test
```

To run provider acceptance tests, run:

```
make testacc
```

### Install

You can install the provider for local testing using:

```
make install
```

This will build the provider if needed, and install it in `~/.terraform.d/plugins`.

> [!NOTE]
> If you're on Windows, you'll have to manually install the built provider binary under
`%APPDATA%/terraform.d/plugins` or `%APPDATA%/HashiCorp/Terraform/plugins`.

### Updating Docs

To update documentation, edit files in the `template/` and/or `examples/` and then run:

```
make generate
```

This will auto-generate the files in `/docs`.
These files should not be updated manually.

> [!NOTE]
> If you are updating Terraform files in `/examples` make sure to run `terraform fmt` on it before
running `make generate` to make sure the Terraform blocks are properly formatted.

## Contributing
Please review our [contributing guidlines](./CONTRIBUTING.md) and
[Code of Conduct](./CODE_OF_CONDUCT.md) before contributing to this project.

Contributions to the project are [released](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service#6-contributions-under-repository-license)
under the [project's open source license](./LICENSE).

For bug reports and feature requests, please create a new
[Github issue](https://github.com/antimetal/terraform-provider-antimetal/issues/new) and provide
as much detail as you can.
