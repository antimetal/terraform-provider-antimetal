// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/antimetal/terraform-provider-antimetal/internal/antimetal/antimetaltest"
	"github.com/antimetal/terraform-provider-antimetal/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	testServer = antimetaltest.NewServer()
	defer testServer.Close()

	os.Exit(m.Run())
}

func TestAccResourceAntimetalHandshake_basic(t *testing.T) {
	roleARN := "arn:aws:iam::012345678999:role/test"
	updatedRoleARN := "arn:aws:iam::012345678999:role/updatedrole"

	handshakeResourceName := "this"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random"},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAntimetalHandshake(handshakeResourceName, roleARN, testServer.URL),
				Check: resource.ComposeTestCheckFunc(
					testutil.TestCheckResourceExists(handshakeResource(handshakeResourceName)),
					resource.TestCheckResourceAttrSet(handshakeResource(handshakeResourceName), "external_id"),
					resource.TestCheckResourceAttrSet(handshakeResource(handshakeResourceName), "handshake_id"),
					resource.TestCheckResourceAttr(handshakeResource(handshakeResourceName), "role_arn", roleARN),
				),
			},
			{
				Config: testAccResourceAntimetalHandshake(handshakeResourceName, updatedRoleARN, testServer.URL),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(handshakeResource(handshakeResourceName), plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func testAccResourceAntimetalHandshake(resourceName, roleARN, url string) string {
	return fmt.Sprintf(`
		provider "antimetal" {
			url = "%s"
		}

		resource "random_uuid" "external_id" {}

		resource "random_uuid" "handshake_id" {}

		resource "antimetal_handshake" "%s" {
			external_id  = random_uuid.external_id.result
			handshake_id = "testhandshake"
			role_arn     = "%s"
		}
	`, url, resourceName, roleARN)
}

//nolint:unparam
func handshakeResource(name string) string {
	return fmt.Sprintf("antimetal_handshake.%s", name)
}
