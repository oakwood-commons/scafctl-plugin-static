// Copyright 2025-2026 Oakwood Commons
// SPDX-License-Identifier: Apache-2.0

// Package main is the entry point for the scafctl-plugin-static plugin.
package main

import (
	"github.com/oakwood-commons/scafctl-plugin-static/internal/static"

	sdkplugin "github.com/oakwood-commons/scafctl-plugin-sdk/plugin"
)

func main() {
	sdkplugin.Serve(&static.Plugin{})
}
