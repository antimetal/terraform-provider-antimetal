//go:build gen

//go:generate go run github.com/google/addlicense -c "Antimetal LLC" -l mpl -s=only -ignore **/tools/* -ignore **/gen/* -ignore **/*.yml -ignore **/*.tf ../
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name antimetal -provider-dir ..

package gen
