// Copyright 2025-2026 Oakwood Commons
// SPDX-License-Identifier: Apache-2.0

// Package static implements the static provider plugin.
package static

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/google/jsonschema-go/jsonschema"
	sdkplugin "github.com/oakwood-commons/scafctl-plugin-sdk/plugin"
	sdkprovider "github.com/oakwood-commons/scafctl-plugin-sdk/provider"
	sdkhelper "github.com/oakwood-commons/scafctl-plugin-sdk/provider/schemahelper"
)

const (
	// ProviderName is the unique identifier for this provider.
	ProviderName = "static"

	// Version is the provider version.
	Version = "1.0.0"
)

// Plugin implements the scafctl ProviderPlugin interface.
type Plugin struct{}

// GetProviders returns the list of providers exposed by this plugin.
//
//nolint:revive // ctx required by interface
func (p *Plugin) GetProviders(_ context.Context) ([]string, error) {
	return []string{ProviderName}, nil
}

// GetProviderDescriptor returns the descriptor for the named provider.
//
//nolint:revive // ctx required by interface
func (p *Plugin) GetProviderDescriptor(_ context.Context, providerName string) (*sdkprovider.Descriptor, error) {
	if providerName != ProviderName {
		return nil, fmt.Errorf("unknown provider: %s", providerName)
	}

	return &sdkprovider.Descriptor{
		Name:        ProviderName,
		DisplayName: "Static Value Provider",
		Description: "Returns a static value without performing any operations. Useful for constants, defaults, and testing.",
		APIVersion:  "v1",
		Version:     semver.MustParse(Version),
		Category:    "Core",
		Tags:        []string{"static", "constant", "testing", "default"},
		Capabilities: []sdkprovider.Capability{
			sdkprovider.CapabilityFrom,
			sdkprovider.CapabilityTransform,
		},
		Schema: sdkhelper.ObjectSchema(
			[]string{"value"},
			map[string]*jsonschema.Schema{
				"value": sdkhelper.AnyProp(
					"The static value to return (can be any type: string, number, boolean, object, array)",
					sdkhelper.WithExample("example-value"),
				),
			},
		),
		OutputSchemas: map[sdkprovider.Capability]*jsonschema.Schema{
			sdkprovider.CapabilityFrom: sdkhelper.ObjectSchema(nil, map[string]*jsonschema.Schema{
				"value": sdkhelper.AnyProp("The static value that was provided (returned directly)", sdkhelper.WithExample("example-value")),
			}),
			sdkprovider.CapabilityTransform: sdkhelper.ObjectSchema(nil, map[string]*jsonschema.Schema{
				"value": sdkhelper.AnyProp("The static value that was provided (returned directly)", sdkhelper.WithExample("example-value")),
			}),
		},
		Examples: []sdkprovider.Example{
			{
				Name:        "String value",
				Description: "Return a static string value",
				YAML:        "name: environment\ntype: static\nfrom:\n  value: production",
			},
			{
				Name:        "Numeric value",
				Description: "Return a static numeric value",
				YAML:        "name: port\ntype: static\nfrom:\n  value: 8080",
			},
			{
				Name:        "Boolean value",
				Description: "Return a static boolean value",
				YAML:        "name: enabled\ntype: static\nfrom:\n  value: true",
			},
			{
				Name:        "Object value",
				Description: "Return a static object/map value",
				YAML:        "name: config\ntype: static\nfrom:\n  value:\n    host: localhost\n    port: 8080\n    ssl: true",
			},
			{
				Name:        "Array value",
				Description: "Return a static array value",
				YAML:        "name: environments\ntype: static\nfrom:\n  value:\n    - dev\n    - staging\n    - production",
			},
		},
	}, nil
}

// ExecuteProvider executes the static provider, returning the input value directly.
//
//nolint:revive // ctx required by interface
func (p *Plugin) ExecuteProvider(_ context.Context, providerName string, input map[string]any) (*sdkprovider.Output, error) {
	if providerName != ProviderName {
		return nil, fmt.Errorf("unknown provider: %s", providerName)
	}

	value, ok := input["value"]
	if !ok {
		return nil, fmt.Errorf("%s: missing required input: value", ProviderName)
	}

	return &sdkprovider.Output{
		Data: value,
	}, nil
}

// DescribeWhatIf returns a description of what the provider would do.
//
//nolint:revive // ctx required by interface
func (p *Plugin) DescribeWhatIf(_ context.Context, providerName string, input map[string]any) (string, error) {
	if providerName != ProviderName {
		return "", fmt.Errorf("unknown provider: %s", providerName)
	}

	if val, ok := input["value"]; ok {
		return fmt.Sprintf("Would return static value: %v", val), nil
	}
	return "Would return static value", nil
}

// ConfigureProvider stores host-side configuration. The static plugin does not
// require any configuration, so this is a no-op.
//
//nolint:revive // ctx and cfg required by interface
func (p *Plugin) ConfigureProvider(_ context.Context, _ string, _ sdkplugin.ProviderConfig) error {
	return nil
}

// ExecuteProviderStream is not supported by the static plugin.
//
//nolint:revive // all params required by interface
func (p *Plugin) ExecuteProviderStream(_ context.Context, _ string, _ map[string]any, _ func(sdkplugin.StreamChunk)) error {
	return sdkplugin.ErrStreamingNotSupported
}

// ExtractDependencies returns resolver keys this input depends on.
//
//nolint:revive // all params required by interface
func (p *Plugin) ExtractDependencies(_ context.Context, _ string, _ map[string]any) ([]string, error) {
	return nil, nil
}

// StopProvider performs cleanup for the named provider.
//
//nolint:revive // all params required by interface
func (p *Plugin) StopProvider(_ context.Context, _ string) error {
	return nil
}
