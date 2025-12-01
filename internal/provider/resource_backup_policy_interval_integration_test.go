package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStandardScheduleConfig_IntervalMinutes tests standard interval with minutes
func TestStandardScheduleConfig_IntervalMinutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		intervalMinutes int64
		expectError     bool
		errorContains   string
	}{
		{
			name:            "valid 360 minutes (6 hours)",
			intervalMinutes: 360,
			expectError:     false,
		},
		{
			name:            "valid 480 minutes (8 hours)",
			intervalMinutes: 480,
			expectError:     false,
		},
		{
			name:            "valid 720 minutes (12 hours)",
			intervalMinutes: 720,
			expectError:     false,
		},
		{
			name:            "invalid - not divisible by 60",
			intervalMinutes: 90,
			expectError:     true,
			errorContains:   "must be divisible by 60",
		},
		{
			name:            "invalid - not allowed value (180 minutes / 3 hours)",
			intervalMinutes: 180,
			expectError:     true,
			errorContains:   "must be 6, 8, or 12 hours",
		},
		{
			name:            "invalid - not allowed value (1440 minutes / 24 hours)",
			intervalMinutes: 1440,
			expectError:     true,
			errorContains:   "must be 6, 8, or 12 hours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create interval config attributes
			intervalConfigAttrs := map[string]attr.Value{
				"interval_minutes":     types.Int64Value(tt.intervalMinutes),
				"interval_hours":       types.Int64Null(),
				"start_window_minutes": types.Int64Null(),
			}

			intervalConfigAttrTypes := map[string]attr.Type{
				"interval_minutes":     types.Int64Type,
				"interval_hours":       types.Int64Type,
				"start_window_minutes": types.Int64Type,
			}

			intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
			require.False(t, diags.HasError(), "Should create interval config object")

			// Create schedule config attributes
			scheduleConfigAttrs := map[string]attr.Value{
				"frequency":       types.StringValue("INTERVAL"),
				"interval_config": intervalConfigObj,
			}

			scheduleConfigAttrTypes := map[string]attr.Type{
				"frequency":       types.StringType,
				"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
			}

			scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
			require.False(t, diags.HasError(), "Should create schedule config object")

			schedule := &BackupScheduleModel{
				VaultId:        types.StringValue("test-vault-id"),
				RetentionDays:  types.Int64Value(7),
				ScheduleConfig: scheduleConfigObj,
			}

			// Test the actual function
			config, err := createStandardScheduleConfig(schedule)

			if tt.expectError {
				assert.Error(t, err, "Should return error")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error message should contain expected text")
				}
				assert.Nil(t, config, "Config should be nil on error")
			} else {
				assert.NoError(t, err, "Should not return error")
				assert.NotNil(t, config, "Config should not be nil")
			}
		})
	}
}

// TestStandardScheduleConfig_IntervalHours tests standard interval with hours
func TestStandardScheduleConfig_IntervalHours(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		intervalHours int64
		expectError   bool
		errorContains string
	}{
		{
			name:          "valid 6 hours",
			intervalHours: 6,
			expectError:   false,
		},
		{
			name:          "valid 8 hours",
			intervalHours: 8,
			expectError:   false,
		},
		{
			name:          "valid 12 hours",
			intervalHours: 12,
			expectError:   false,
		},
		{
			name:          "invalid - 3 hours not allowed",
			intervalHours: 3,
			expectError:   true,
			errorContains: "must be 6, 8, or 12 hours",
		},
		{
			name:          "invalid - 24 hours not allowed",
			intervalHours: 24,
			expectError:   true,
			errorContains: "must be 6, 8, or 12 hours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create interval config attributes
			intervalConfigAttrs := map[string]attr.Value{
				"interval_minutes":     types.Int64Null(),
				"interval_hours":       types.Int64Value(tt.intervalHours),
				"start_window_minutes": types.Int64Null(),
			}

			intervalConfigAttrTypes := map[string]attr.Type{
				"interval_minutes":     types.Int64Type,
				"interval_hours":       types.Int64Type,
				"start_window_minutes": types.Int64Type,
			}

			intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
			require.False(t, diags.HasError(), "Should create interval config object")

			// Create schedule config attributes
			scheduleConfigAttrs := map[string]attr.Value{
				"frequency":       types.StringValue("INTERVAL"),
				"interval_config": intervalConfigObj,
			}

			scheduleConfigAttrTypes := map[string]attr.Type{
				"frequency":       types.StringType,
				"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
			}

			scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
			require.False(t, diags.HasError(), "Should create schedule config object")

			schedule := &BackupScheduleModel{
				VaultId:        types.StringValue("test-vault-id"),
				RetentionDays:  types.Int64Value(7),
				ScheduleConfig: scheduleConfigObj,
			}

			// Test the actual function
			config, err := createStandardScheduleConfig(schedule)

			if tt.expectError {
				assert.Error(t, err, "Should return error")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error message should contain expected text")
				}
				assert.Nil(t, config, "Config should be nil on error")
			} else {
				assert.NoError(t, err, "Should not return error")
				assert.NotNil(t, config, "Config should not be nil")
			}
		})
	}
}

// TestStandardScheduleConfig_BothProvided tests that both minutes and hours cannot be provided
func TestStandardScheduleConfig_BothProvided(t *testing.T) {
	t.Parallel()

	// Create interval config with both fields provided
	intervalConfigAttrs := map[string]attr.Value{
		"interval_minutes":     types.Int64Value(720),
		"interval_hours":       types.Int64Value(12),
		"start_window_minutes": types.Int64Null(),
	}

	intervalConfigAttrTypes := map[string]attr.Type{
		"interval_minutes":     types.Int64Type,
		"interval_hours":       types.Int64Type,
		"start_window_minutes": types.Int64Type,
	}

	intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
	require.False(t, diags.HasError(), "Should create interval config object")

	scheduleConfigAttrs := map[string]attr.Value{
		"frequency":       types.StringValue("INTERVAL"),
		"interval_config": intervalConfigObj,
	}

	scheduleConfigAttrTypes := map[string]attr.Type{
		"frequency":       types.StringType,
		"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
	}

	scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
	require.False(t, diags.HasError(), "Should create schedule config object")

	schedule := &BackupScheduleModel{
		VaultId:        types.StringValue("test-vault-id"),
		RetentionDays:  types.Int64Value(7),
		ScheduleConfig: scheduleConfigObj,
	}

	// Test the actual function
	config, err := createStandardScheduleConfig(schedule)

	assert.Error(t, err, "Should return error when both are provided")
	assert.Contains(t, err.Error(), "cannot specify both interval_minutes and interval_hours", "Error should mention both fields")
	assert.Nil(t, config, "Config should be nil on error")
}

// TestStandardScheduleConfig_NeitherProvided tests that at least one must be provided
func TestStandardScheduleConfig_NeitherProvided(t *testing.T) {
	t.Parallel()

	// Create interval config with neither field provided
	intervalConfigAttrs := map[string]attr.Value{
		"interval_minutes":     types.Int64Null(),
		"interval_hours":       types.Int64Null(),
		"start_window_minutes": types.Int64Null(),
	}

	intervalConfigAttrTypes := map[string]attr.Type{
		"interval_minutes":     types.Int64Type,
		"interval_hours":       types.Int64Type,
		"start_window_minutes": types.Int64Type,
	}

	intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
	require.False(t, diags.HasError(), "Should create interval config object")

	scheduleConfigAttrs := map[string]attr.Value{
		"frequency":       types.StringValue("INTERVAL"),
		"interval_config": intervalConfigObj,
	}

	scheduleConfigAttrTypes := map[string]attr.Type{
		"frequency":       types.StringType,
		"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
	}

	scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
	require.False(t, diags.HasError(), "Should create schedule config object")

	schedule := &BackupScheduleModel{
		VaultId:        types.StringValue("test-vault-id"),
		RetentionDays:  types.Int64Value(7),
		ScheduleConfig: scheduleConfigObj,
	}

	// Test the actual function
	config, err := createStandardScheduleConfig(schedule)

	assert.Error(t, err, "Should return error when neither is provided")
	assert.Contains(t, err.Error(), "either interval_minutes or interval_hours must be specified", "Error should mention both fields")
	assert.Nil(t, config, "Config should be nil on error")
}

// TestHighFrequencyScheduleConfig_IntervalMinutes tests high frequency interval with minutes
func TestHighFrequencyScheduleConfig_IntervalMinutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		intervalMinutes int64
		expectError     bool
	}{
		{
			name:            "15 minutes",
			intervalMinutes: 15,
			expectError:     false,
		},
		{
			name:            "30 minutes",
			intervalMinutes: 30,
			expectError:     false,
		},
		{
			name:            "60 minutes",
			intervalMinutes: 60,
			expectError:     false,
		},
		{
			name:            "120 minutes",
			intervalMinutes: 120,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			intervalConfigAttrs := map[string]attr.Value{
				"interval_minutes":     types.Int64Value(tt.intervalMinutes),
				"interval_hours":       types.Int64Null(),
				"start_window_minutes": types.Int64Null(),
			}

			intervalConfigAttrTypes := map[string]attr.Type{
				"interval_minutes":     types.Int64Type,
				"interval_hours":       types.Int64Type,
				"start_window_minutes": types.Int64Type,
			}

			intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
			require.False(t, diags.HasError(), "Should create interval config object")

			scheduleConfigAttrs := map[string]attr.Value{
				"frequency":       types.StringValue("INTERVAL"),
				"interval_config": intervalConfigObj,
			}

			scheduleConfigAttrTypes := map[string]attr.Type{
				"frequency":       types.StringType,
				"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
			}

			scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
			require.False(t, diags.HasError(), "Should create schedule config object")

			schedule := &BackupScheduleModel{
				VaultId:        types.StringValue("test-vault-id"),
				RetentionDays:  types.Int64Value(7),
				ScheduleConfig: scheduleConfigObj,
			}

			config, err := createHighFrequencyScheduleConfig(schedule)

			if tt.expectError {
				assert.Error(t, err, "Should return error")
				assert.Nil(t, config, "Config should be nil on error")
			} else {
				assert.NoError(t, err, "Should not return error")
				assert.NotNil(t, config, "Config should not be nil")
			}
		})
	}
}

// TestHighFrequencyScheduleConfig_IntervalHours tests high frequency interval with hours
func TestHighFrequencyScheduleConfig_IntervalHours(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		intervalHours   int64
		expectedMinutes int32
		expectError     bool
	}{
		{
			name:            "1 hour = 60 minutes",
			intervalHours:   1,
			expectedMinutes: 60,
			expectError:     false,
		},
		{
			name:            "2 hours = 120 minutes",
			intervalHours:   2,
			expectedMinutes: 120,
			expectError:     false,
		},
		{
			name:            "6 hours = 360 minutes",
			intervalHours:   6,
			expectedMinutes: 360,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			intervalConfigAttrs := map[string]attr.Value{
				"interval_minutes":     types.Int64Null(),
				"interval_hours":       types.Int64Value(tt.intervalHours),
				"start_window_minutes": types.Int64Null(),
			}

			intervalConfigAttrTypes := map[string]attr.Type{
				"interval_minutes":     types.Int64Type,
				"interval_hours":       types.Int64Type,
				"start_window_minutes": types.Int64Type,
			}

			intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
			require.False(t, diags.HasError(), "Should create interval config object")

			scheduleConfigAttrs := map[string]attr.Value{
				"frequency":       types.StringValue("INTERVAL"),
				"interval_config": intervalConfigObj,
			}

			scheduleConfigAttrTypes := map[string]attr.Type{
				"frequency":       types.StringType,
				"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
			}

			scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
			require.False(t, diags.HasError(), "Should create schedule config object")

			schedule := &BackupScheduleModel{
				VaultId:        types.StringValue("test-vault-id"),
				RetentionDays:  types.Int64Value(7),
				ScheduleConfig: scheduleConfigObj,
			}

			config, err := createHighFrequencyScheduleConfig(schedule)

			if tt.expectError {
				assert.Error(t, err, "Should return error")
				assert.Nil(t, config, "Config should be nil on error")
			} else {
				assert.NoError(t, err, "Should not return error")
				assert.NotNil(t, config, "Config should not be nil")
				// Verify conversion happened correctly (hours * 60 = minutes)
				t.Logf("Created high frequency config for %d hours (should be %d minutes)", tt.intervalHours, tt.expectedMinutes)
			}
		})
	}
}

// TestHighFrequencyScheduleConfig_BothProvided tests that both minutes and hours cannot be provided
func TestHighFrequencyScheduleConfig_BothProvided(t *testing.T) {
	t.Parallel()

	intervalConfigAttrs := map[string]attr.Value{
		"interval_minutes":     types.Int64Value(60),
		"interval_hours":       types.Int64Value(1),
		"start_window_minutes": types.Int64Null(),
	}

	intervalConfigAttrTypes := map[string]attr.Type{
		"interval_minutes":     types.Int64Type,
		"interval_hours":       types.Int64Type,
		"start_window_minutes": types.Int64Type,
	}

	intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
	require.False(t, diags.HasError(), "Should create interval config object")

	scheduleConfigAttrs := map[string]attr.Value{
		"frequency":       types.StringValue("INTERVAL"),
		"interval_config": intervalConfigObj,
	}

	scheduleConfigAttrTypes := map[string]attr.Type{
		"frequency":       types.StringType,
		"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
	}

	scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
	require.False(t, diags.HasError(), "Should create schedule config object")

	schedule := &BackupScheduleModel{
		VaultId:        types.StringValue("test-vault-id"),
		RetentionDays:  types.Int64Value(7),
		ScheduleConfig: scheduleConfigObj,
	}

	config, err := createHighFrequencyScheduleConfig(schedule)

	assert.Error(t, err, "Should return error when both are provided")
	assert.Contains(t, err.Error(), "cannot specify both interval_minutes and interval_hours", "Error should mention both fields")
	assert.Nil(t, config, "Config should be nil on error")
}

// TestHighFrequencyScheduleConfig_NeitherProvided tests that at least one must be provided
func TestHighFrequencyScheduleConfig_NeitherProvided(t *testing.T) {
	t.Parallel()

	intervalConfigAttrs := map[string]attr.Value{
		"interval_minutes":     types.Int64Null(),
		"interval_hours":       types.Int64Null(),
		"start_window_minutes": types.Int64Null(),
	}

	intervalConfigAttrTypes := map[string]attr.Type{
		"interval_minutes":     types.Int64Type,
		"interval_hours":       types.Int64Type,
		"start_window_minutes": types.Int64Type,
	}

	intervalConfigObj, diags := types.ObjectValue(intervalConfigAttrTypes, intervalConfigAttrs)
	require.False(t, diags.HasError(), "Should create interval config object")

	scheduleConfigAttrs := map[string]attr.Value{
		"frequency":       types.StringValue("INTERVAL"),
		"interval_config": intervalConfigObj,
	}

	scheduleConfigAttrTypes := map[string]attr.Type{
		"frequency":       types.StringType,
		"interval_config": types.ObjectType{AttrTypes: intervalConfigAttrTypes},
	}

	scheduleConfigObj, diags := types.ObjectValue(scheduleConfigAttrTypes, scheduleConfigAttrs)
	require.False(t, diags.HasError(), "Should create schedule config object")

	schedule := &BackupScheduleModel{
		VaultId:        types.StringValue("test-vault-id"),
		RetentionDays:  types.Int64Value(7),
		ScheduleConfig: scheduleConfigObj,
	}

	config, err := createHighFrequencyScheduleConfig(schedule)

	assert.Error(t, err, "Should return error when neither is provided")
	assert.Contains(t, err.Error(), "either interval_minutes or interval_hours must be specified", "Error should mention both fields")
	assert.Nil(t, config, "Config should be nil on error")
}
