// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"

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

func TestAccResourceAntimetalHandshake(t *testing.T) {
	roleARN := "arn:aws:iam::012345678999:role/test"
	updatedRoleARN := "arn:aws:iam::012345678999:role/updatedrole"

	handshakeResourceName := "this"
	rsrc := handshakeResource(handshakeResourceName)
	handshakeID := uuid.New().String()
	externalID := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAntimetalHandshake(
					handshakeResourceName, roleARN, testServer.URL, handshakeID, externalID,
				),
				Check: resource.ComposeTestCheckFunc(
					testutil.TestCheckResourceExists(rsrc),
					resource.TestCheckResourceAttrSet(rsrc, "external_id"),
					resource.TestCheckResourceAttrSet(rsrc, "handshake_id"),
					resource.TestCheckResourceAttr(rsrc, "role_arn", roleARN),
				),
			},
			{
				ResourceName:                         rsrc,
				ImportState:                          true,
				ImportStateId:                        fmt.Sprintf("%s;%s;%s", handshakeID, externalID, roleARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "handshake_id",
			},
			{
				Config: testAccResourceAntimetalHandshake(
					handshakeResourceName, updatedRoleARN, testServer.URL, handshakeID, externalID,
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rsrc, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func testAccResourceAntimetalHandshake(resourceName, roleARN, url, handshakeID, externalID string) string {
	return fmt.Sprintf(`
		provider "antimetal" {
			url = "%s"
		}

		resource "antimetal_handshake" "%s" {
			external_id  = "%s"
			handshake_id = "%s"
			role_arn     = "%s"
		}
	`, url, resourceName, externalID, handshakeID, roleARN)
}

//nolint:unparam
func handshakeResource(name string) string {
	return fmt.Sprintf("antimetal_handshake.%s", name)
}
