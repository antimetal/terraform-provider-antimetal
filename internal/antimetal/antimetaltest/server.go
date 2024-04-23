// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package antimetaltest

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/antimetal/terraform-provider-antimetal/internal/antimetal"
)

func NewServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/webhook/terraform", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			handleErr(fmt.Errorf("couldn't read body: %w", err),
				w, http.StatusInternalServerError,
			)
			return
		}

		var req antimetal.HandshakeRequest
		if err := json.Unmarshal(body, &req); err != nil {
			handleErr(fmt.Errorf("couldn't unmarshal json body: %w", err),
				w, http.StatusInternalServerError,
			)
			return
		}

		if req.HandshakeID == "" {
			handleErr(fmt.Errorf("handshake_id must be set"), w, http.StatusBadRequest)
			return
		}

		if req.ExternalID == "" {
			handleErr(fmt.Errorf("external_id must be set"), w, http.StatusBadRequest)
			return
		}

		if req.RoleARN == "" {
			handleErr(fmt.Errorf("role_arn must be set"), w, http.StatusBadRequest)
			return
		}

		switch req.Action {
		case antimetal.HandshakeCreate:
		case antimetal.HandshakeDelete:
		default:
			handleErr(fmt.Errorf("invalid action"), w, http.StatusBadRequest)
			return
		}
	})

	return httptest.NewServer(mux)
}

func handleErr(err error, w http.ResponseWriter, statusCode int) {
	slog.Error(err.Error())
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		slog.Error(err.Error())
	}
}
