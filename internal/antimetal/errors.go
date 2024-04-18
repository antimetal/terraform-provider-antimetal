// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package antimetal

import "fmt"

type HTTPError struct {
	URL        string
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf(
		"request to %s failed; status: %d; body: %s", e.URL, e.StatusCode, e.Body,
	)
}
