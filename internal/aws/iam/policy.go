// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package iam

type PolicyDocument struct {
	Version    string            `json:",omitempty"`
	Statements []PolicyStatement `json:"Statement,omitempty"`
}

type PolicyStatement struct {
	Effect     string                   `json:",omitempty"`
	Actions    any                      `json:"Action,omitempty"`
	Resources  any                      `json:"Resource,omitempty"`
	Principals PolicyStatementPrincipal `json:"Principal,omitempty"`
	Conditions PolicyStatementCondition `json:"Condition,omitempty"`
}

// TODO: do a better job with the type defintions. It's not super robust and
// pretty opaque but it works for now.

type PolicyStatementPrincipal map[string]any

type PolicyStatementCondition map[string]any
