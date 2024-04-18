// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package antimetal

type HandshakeAction string

const (
	HandshakeCreate HandshakeAction = "CREATE"
	HandshakeDelete HandshakeAction = "DELETE"
)

type HandshakeRequest struct {
	Action      HandshakeAction `json:"action"`
	HandshakeID string          `json:"handshake_id"`
	ExternalID  string          `json:"external_id"`
	RoleARN     string          `json:"role_arn"`
}
