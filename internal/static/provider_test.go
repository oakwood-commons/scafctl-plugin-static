// Copyright 2025-2026 Oakwood Commons
// SPDX-License-Identifier: Apache-2.0

package static

import (
	"context"
	"testing"

	sdkplugin "github.com/oakwood-commons/scafctl-plugin-sdk/plugin"
	sdkprovider "github.com/oakwood-commons/scafctl-plugin-sdk/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPlugin() *Plugin { return &Plugin{} }

func TestGetProviders(t *testing.T) {
	providers, err := newPlugin().GetProviders(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{ProviderName}, providers)
}

func TestGetProviderDescriptor(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		wantErr     string
		wantName    string
		wantCaps    []sdkprovider.Capability
		wantVersion string
	}{
		{
			name:     "valid static provider",
			provider: "static",
			wantName: "static",
			wantCaps: []sdkprovider.Capability{
				sdkprovider.CapabilityFrom,
				sdkprovider.CapabilityTransform,
			},
			wantVersion: "1.0.0",
		},
		{
			name:     "unknown provider",
			provider: "nonexistent",
			wantErr:  "unknown provider: nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := newPlugin().GetProviderDescriptor(context.Background(), tt.provider)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				assert.Nil(t, desc)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, desc)
			assert.Equal(t, tt.wantName, desc.Name)
			assert.Equal(t, "Static Value Provider", desc.DisplayName)
			assert.Equal(t, "v1", desc.APIVersion)
			assert.Equal(t, tt.wantVersion, desc.Version.String())
			assert.Equal(t, tt.wantCaps, desc.Capabilities)
			assert.NotNil(t, desc.Schema)
			assert.NotEmpty(t, desc.Schema.Properties)
			assert.Contains(t, desc.OutputSchemas, sdkprovider.CapabilityFrom)
			assert.Contains(t, desc.OutputSchemas, sdkprovider.CapabilityTransform)
			assert.NotContains(t, desc.OutputSchemas, sdkprovider.CapabilityAction)
			assert.NotEmpty(t, desc.Examples)
			assert.Equal(t, "Core", desc.Category)
			assert.Equal(t, []string{"static", "constant", "testing", "default"}, desc.Tags)
		})
	}
}

func TestExecuteProvider(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		want    any
		wantErr string
	}{
		{
			name:  "string value",
			input: map[string]any{"value": "test-value"},
			want:  "test-value",
		},
		{
			name:  "numeric value",
			input: map[string]any{"value": 42},
			want:  42,
		},
		{
			name:  "boolean value",
			input: map[string]any{"value": true},
			want:  true,
		},
		{
			name: "object value",
			input: map[string]any{
				"value": map[string]any{
					"key1": "value1",
					"key2": 123,
					"key3": true,
				},
			},
			want: map[string]any{
				"key1": "value1",
				"key2": 123,
				"key3": true,
			},
		},
		{
			name:  "array value",
			input: map[string]any{"value": []any{"item1", "item2", "item3"}},
			want:  []any{"item1", "item2", "item3"},
		},
		{
			name:  "nil value",
			input: map[string]any{"value": nil},
			want:  nil,
		},
		{
			name:    "missing value",
			input:   map[string]any{},
			wantErr: "static: missing required input: value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := newPlugin().ExecuteProvider(context.Background(), ProviderName, tt.input)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				assert.Nil(t, out)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, out)
			assert.Equal(t, tt.want, out.Data)
		})
	}
}

func TestExecuteProvider_UnknownProvider(t *testing.T) {
	_, err := newPlugin().ExecuteProvider(context.Background(), "unknown", nil)
	require.EqualError(t, err, "unknown provider: unknown")
}

func TestDescribeWhatIf(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		input    map[string]any
		want     string
		wantErr  string
	}{
		{
			name:     "with string value",
			provider: ProviderName,
			input:    map[string]any{"value": "hello"},
			want:     "Would return static value: hello",
		},
		{
			name:     "with numeric value",
			provider: ProviderName,
			input:    map[string]any{"value": 42},
			want:     "Would return static value: 42",
		},
		{
			name:     "missing value",
			provider: ProviderName,
			input:    map[string]any{},
			want:     "Would return static value",
		},
		{
			name:     "nil value",
			provider: ProviderName,
			input:    map[string]any{"value": nil},
			want:     "Would return static value: <nil>",
		},
		{
			name:     "unknown provider",
			provider: "other",
			input:    map[string]any{},
			wantErr:  "unknown provider: other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPlugin().DescribeWhatIf(context.Background(), tt.provider, tt.input)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigureProvider(t *testing.T) {
	err := newPlugin().ConfigureProvider(context.Background(), ProviderName, sdkplugin.ProviderConfig{})
	require.NoError(t, err)
}

func TestExecuteProviderStream(t *testing.T) {
	err := newPlugin().ExecuteProviderStream(context.Background(), ProviderName, nil, nil)
	require.ErrorIs(t, err, sdkplugin.ErrStreamingNotSupported)
}

func TestExtractDependencies(t *testing.T) {
	deps, err := newPlugin().ExtractDependencies(context.Background(), ProviderName, nil)
	require.NoError(t, err)
	assert.Nil(t, deps)
}

func TestStopProvider(t *testing.T) {
	err := newPlugin().StopProvider(context.Background(), ProviderName)
	require.NoError(t, err)
}

func BenchmarkExecuteProvider(b *testing.B) {
	p := newPlugin()
	ctx := context.Background()

	b.Run("string_value", func(b *testing.B) {
		input := map[string]any{"value": "hello-world"}
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_, _ = p.ExecuteProvider(ctx, ProviderName, input)
		}
	})

	b.Run("map_value", func(b *testing.B) {
		input := map[string]any{
			"value": map[string]any{
				"host": "localhost",
				"port": 8080,
				"ssl":  true,
			},
		}
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_, _ = p.ExecuteProvider(ctx, ProviderName, input)
		}
	})

	b.Run("list_value", func(b *testing.B) {
		input := map[string]any{
			"value": []any{"dev", "staging", "production"},
		}
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_, _ = p.ExecuteProvider(ctx, ProviderName, input)
		}
	})
}
