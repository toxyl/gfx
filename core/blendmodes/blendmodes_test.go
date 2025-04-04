package blendmodes

import (
	"strings"
	"testing"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/math"
)

// TestBlendModeRegistration tests the registration of blend modes
func TestBlendModeRegistration(t *testing.T) {
	tests := []struct {
		name     string
		category string
	}{
		{"normal", constants.CategoryBasic},
		{"erase", constants.CategoryBasic},
		{"darken", constants.CategoryDarken},
		{"multiply", constants.CategoryDarken},
		{"colorburn", constants.CategoryDarken},
		{"linearburn", constants.CategoryDarken},
		{"darkercolor", constants.CategoryDarken},
		{"lighten", constants.CategoryLighten},
		{"screen", constants.CategoryLighten},
		{"colordodge", constants.CategoryLighten},
		{"add", constants.CategoryLighten},
		{"lightercolor", constants.CategoryLighten},
		{"overlay", constants.CategoryContrast},
		{"softlight", constants.CategoryContrast},
		{"hardlight", constants.CategoryContrast},
		{"vividlight", constants.CategoryContrast},
		{"linearlight", constants.CategoryContrast},
		{"pinlight", constants.CategoryContrast},
		{"hardmix", constants.CategoryContrast},
		{"difference", constants.CategoryComparative},
		{"exclusion", constants.CategoryComparative},
		{"subtract", constants.CategoryComparative},
		{"divide", constants.CategoryComparative},
		{"negation", constants.CategoryComparative},
		{"contrastnegate", constants.CategoryComparative},
		{"hue", constants.CategoryComponent},
		{"saturation", constants.CategoryComponent},
		{"luminosity", constants.CategoryComponent},
		{"reflect", constants.CategorySpecial},
		{"glow", constants.CategorySpecial},
		{"average", constants.CategorySpecial},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, err := Get(tt.name)
			if err != nil {
				t.Errorf("blend mode %s not registered: %v", tt.name, err)
				return
			}
			if mode.Meta().Category() != tt.category {
				t.Errorf("blend mode %s has wrong category, got %s, want %s", tt.name, mode.Meta().Category(), tt.category)
			}
		})
	}
}

// TestBlendModeBasicOperations tests basic blend mode operations
func TestBlendModeBasicOperations(t *testing.T) {
	tests := []struct {
		name     string
		base     *color.RGBA64
		blend    *color.RGBA64
		alpha    float64
		expected *color.RGBA64
	}{
		{
			name:     "normal_basic",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
			alpha:    1.0,
			expected: mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
		},
		{
			name:     "multiply_dark",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.2, 0.2, 0.2, 1.0),
			alpha:    1.0,
			expected: mustNewRGBA64(t, 0.1, 0.1, 0.1, 1.0),
		},
		{
			name:     "screen_light",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
			alpha:    1.0,
			expected: mustNewRGBA64(t, 0.9, 0.9, 0.9, 1.0),
		},
		{
			name:     "darken_basic",
			base:     mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			blend:    mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			alpha:    1.0,
			expected: mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
		},
		{
			name:     "overlay_basic",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
			alpha:    1.0,
			expected: mustNewRGBA64(t, 0.85, 0.85, 0.85, 1.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode(strings.Split(tt.name, "_")[0])
			if mode == nil {
				t.Fatalf("blend mode not found: %s", tt.name)
			}

			result, err := mode.Blend(tt.base, tt.blend, tt.alpha)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !validateRGBA64(result, tt.expected, 0.2) {
				t.Errorf("unexpected result for %s: got %+v, want %+v", tt.name, result, tt.expected)
			}
		})
	}
}

// Helper functions

type testHelper interface {
	Helper()
	Fatalf(format string, args ...interface{})
}

func mustNewRGBA64(t testHelper, r, g, b, a float64) *color.RGBA64 {
	t.Helper()
	c, err := color.NewRGBA64(r, g, b, a)
	if err != nil {
		t.Fatalf("failed to create RGBA64 color: %v", err)
	}
	return c
}

func validateRGBA64(got, want *color.RGBA64, epsilon float64) bool {
	if got == nil || want == nil {
		return got == want
	}
	gotR, gotG, gotB, gotA := got.R, got.G, got.B, got.A
	wantR, wantG, wantB, wantA := want.R, want.G, want.B, want.A

	return math.Abs(gotR-wantR) <= epsilon &&
		math.Abs(gotG-wantG) <= epsilon &&
		math.Abs(gotB-wantB) <= epsilon &&
		math.Abs(gotA-wantA) <= epsilon
}

// TestBlendModeValidation tests validation functions
func TestBlendModeValidation(t *testing.T) {
	tests := []struct {
		name        string
		blendName   string
		description string
		category    string
		blendFunc   func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64
		wantErr     bool
	}{
		{
			name:        "valid_registration",
			blendName:   "test_blend",
			description: "Test blend mode",
			category:    constants.CategoryBasic,
			blendFunc: func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
				// Simple passthrough function for testing
				return &color.RGBA64{
					R: bottom.R,
					G: bottom.G,
					B: bottom.B,
					A: bottom.A,
				}
			},
			wantErr: false,
		},
		{
			name:        "invalid_category",
			blendName:   "test_blend",
			description: "Test blend mode",
			category:    "invalid_category",
			blendFunc:   func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 { return bottom },
			wantErr:     true,
		},
		{
			name:        "empty_name",
			blendName:   "",
			description: "Test blend mode",
			category:    constants.CategoryBasic,
			blendFunc:   func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 { return bottom },
			wantErr:     true,
		},
		{
			name:        "nil_function",
			blendName:   "test_blend",
			description: "Test blend mode",
			category:    constants.CategoryBasic,
			blendFunc:   nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegistration(tt.blendName, tt.description, tt.category, tt.blendFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRegistration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestBlendModeEdgeCases tests edge cases for blend modes
func TestBlendModeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		base     *color.RGBA64
		blend    *color.RGBA64
		alpha    float64
		wantErr  bool
		expected *color.RGBA64
	}{
		{
			name:    "nil_base",
			base:    nil,
			blend:   mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			alpha:   1.0,
			wantErr: true,
		},
		{
			name:    "nil_blend",
			base:    mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:   nil,
			alpha:   1.0,
			wantErr: true,
		},
		{
			name:    "invalid_alpha_low",
			base:    mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:   mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			alpha:   -0.1,
			wantErr: true,
		},
		{
			name:    "invalid_alpha_high",
			base:    mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:   mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			alpha:   1.1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode("normal") // Use normal blend mode for edge cases
			if mode == nil {
				t.Fatal("normal blend mode not found")
			}

			result, err := mode.Blend(tt.base, tt.blend, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !validateRGBA64(result, tt.expected, 0.2) {
				t.Errorf("unexpected result: got %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// TestGetBlendMode tests the GetBlendMode function
func TestGetBlendMode(t *testing.T) {
	tests := []struct {
		name     string
		modeName string
		wantNil  bool
	}{
		{"normal_exists", "normal", false},
		{"multiply_exists", "multiply", false},
		{"nonexistent", "nonexistent", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode(tt.modeName)
			if (mode == nil) != tt.wantNil {
				t.Errorf("GetBlendMode(%q) got nil: %v, want nil: %v", tt.modeName, mode == nil, tt.wantNil)
			}
		})
	}
}

// TestListFunctions tests the List and ListCategories functions
func TestListFunctions(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		modes := List()
		if len(modes) == 0 {
			t.Error("List() returned empty list")
		}

		// Check for required modes
		required := []string{"normal", "multiply", "screen", "overlay"}
		for _, name := range required {
			found := false
			for _, mode := range modes {
				if mode == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("List() missing required mode: %s", name)
			}
		}
	})

	t.Run("ListCategories", func(t *testing.T) {
		categories := ListCategories()
		required := []string{
			constants.CategoryBasic,
			constants.CategoryDarken,
			constants.CategoryLighten,
			constants.CategoryContrast,
			constants.CategoryComparative,
			constants.CategoryComponent,
			constants.CategorySpecial,
		}

		if len(categories) != len(required) {
			t.Errorf("ListCategories() got %d categories, want %d", len(categories), len(required))
		}

		for _, cat := range required {
			found := false
			for _, c := range categories {
				if c == cat {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("ListCategories() missing required category: %s", cat)
			}
		}
	})
}

// validateFloat checks if a float value is within the expected range
func validateFloat(t *testing.T, got, want, epsilon float64, name string) {
	t.Helper()
	if math.Abs(got-want) > epsilon {
		t.Errorf("%s: got %v, want %v ± %v", name, got, want, epsilon)
	}
}

// GetBlendMode returns a blend mode by name for testing
func GetBlendMode(name string) *IBlendMode {
	mode, err := Get(name)
	if err != nil {
		return nil
	}
	return mode
}

func TestBlendModeList(t *testing.T) {
	modes := List()
	if len(modes) == 0 {
		t.Error("expected non-empty list of blend modes")
	}

	// Check for required blend modes
	requiredModes := []string{"normal", "multiply", "screen", "darken", "overlay"}
	for _, required := range requiredModes {
		found := false
		for _, mode := range modes {
			if mode == required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("required blend mode %q not found in list", required)
		}
	}
}

func TestBlendModeListCategories(t *testing.T) {
	categories := ListCategories()
	if len(categories) == 0 {
		t.Error("expected non-empty list of categories")
	}

	// Check for required categories
	requiredCategories := []string{
		constants.CategoryBasic,
		constants.CategoryDarken,
		constants.CategoryLighten,
		constants.CategoryContrast,
		constants.CategoryComparative,
		constants.CategoryComponent,
		constants.CategorySpecial,
	}
	for _, required := range requiredCategories {
		found := false
		for _, category := range categories {
			if category == required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("required category %q not found in list", required)
		}
	}
}

func TestGetBlendModeErrors(t *testing.T) {
	tests := []struct {
		name        string
		modeName    string
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty_name",
			modeName:    "",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:        "nonexistent_mode",
			modeName:    "nonexistent",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:        "invalid_chars",
			modeName:    "invalid!@#$",
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, err := Get(tt.modeName)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
			}
			if tt.wantErr && mode != nil {
				t.Error("expected nil mode but got non-nil")
			}
		})
	}
}

func TestGetByCategory(t *testing.T) {
	tests := []struct {
		name        string
		category    string
		wantErr     bool
		errContains string
		minModes    int
	}{
		{
			name:     "basic_category",
			category: constants.CategoryBasic,
			minModes: 1,
		},
		{
			name:     "darken_category",
			category: constants.CategoryDarken,
			minModes: 3,
		},
		{
			name:     "lighten_category",
			category: constants.CategoryLighten,
			minModes: 3,
		},
		{
			name:        "empty_category",
			category:    "",
			wantErr:     true,
			errContains: "no blend modes found",
		},
		{
			name:        "nonexistent_category",
			category:    "nonexistent",
			wantErr:     true,
			errContains: "no blend modes found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modes, err := GetByCategory(tt.category)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(modes) < tt.minModes {
					t.Errorf("expected at least %d modes, got %d", tt.minModes, len(modes))
				}
				// Check that modes are sorted by name
				for i := 1; i < len(modes); i++ {
					if modes[i-1].meta.Name() > modes[i].meta.Name() {
						t.Error("modes are not sorted by name")
						break
					}
				}
			}
		})
	}
}

func TestBlendModeBlendErrors(t *testing.T) {
	tests := []struct {
		name        string
		bottom      *color.RGBA64
		top         *color.RGBA64
		alpha       float64
		wantErr     bool
		errContains string
	}{
		{
			name:        "nil_bottom",
			bottom:      nil,
			top:         mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			alpha:       1.0,
			wantErr:     true,
			errContains: "cannot be nil",
		},
		{
			name:        "nil_top",
			bottom:      mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			top:         nil,
			alpha:       1.0,
			wantErr:     true,
			errContains: "cannot be nil",
		},
		{
			name:        "negative_alpha",
			bottom:      mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			top:         mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			alpha:       -0.1,
			wantErr:     true,
			errContains: "between 0.0 and 1.0",
		},
		{
			name:        "alpha_too_high",
			bottom:      mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			top:         mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			alpha:       1.1,
			wantErr:     true,
			errContains: "between 0.0 and 1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, err := Get("normal")
			if err != nil {
				t.Fatalf("failed to get normal blend mode: %v", err)
			}

			result, err := mode.Blend(tt.bottom, tt.top, tt.alpha)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected non-nil result")
				}
			}
		})
	}
}

func TestBlendModeCombinations(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		base     *color.RGBA64
		blend    *color.RGBA64
		expected *color.RGBA64
	}{
		{
			name:     "screen_white_on_black",
			mode:     "screen",
			base:     mustNewRGBA64(t, 0.0, 0.0, 0.0, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			expected: mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
		},
		{
			name:     "multiply_black_on_white",
			mode:     "multiply",
			base:     mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			blend:    mustNewRGBA64(t, 0.0, 0.0, 0.0, 1.0),
			expected: mustNewRGBA64(t, 0.0, 0.0, 0.0, 1.0),
		},
		{
			name:     "overlay_gray_on_gray",
			mode:     "overlay",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "darken_different_colors",
			mode:     "darken",
			base:     mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			blend:    mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			expected: mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
		},
		{
			name:     "lighten_different_colors",
			mode:     "lighten",
			base:     mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
		},
		{
			name:     "add_colors",
			mode:     "add",
			base:     mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.4, 0.4, 0.4, 1.0),
			expected: mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
		},
		{
			name:     "subtract_colors",
			mode:     "subtract",
			base:     mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			blend:    mustNewRGBA64(t, 0.4, 0.4, 0.4, 1.0),
			expected: mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
		},
		{
			name:     "average_colors",
			mode:     "average",
			base:     mustNewRGBA64(t, 0.2, 0.2, 0.2, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "difference_colors",
			mode:     "difference",
			base:     mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			blend:    mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			expected: mustNewRGBA64(t, 0.4, 0.4, 0.4, 1.0),
		},
		{
			name:     "exclusion_colors",
			mode:     "exclusion",
			base:     mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			blend:    mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			expected: mustNewRGBA64(t, 0.58, 0.58, 0.58, 1.0),
		},
		{
			name:     "divide_colors",
			mode:     "divide",
			base:     mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
			blend:    mustNewRGBA64(t, 0.4, 0.4, 0.4, 1.0),
			expected: mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
		},
		{
			name:     "colorburn_colors",
			mode:     "colorburn",
			base:     mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			blend:    mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			expected: mustNewRGBA64(t, 0.0, 0.0, 0.0, 1.0),
		},
		{
			name:     "linearburn_colors",
			mode:     "linearburn",
			base:     mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			blend:    mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			expected: mustNewRGBA64(t, 0.0, 0.0, 0.0, 1.0),
		},
		{
			name:     "softlight_colors",
			mode:     "softlight",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.65, 0.65, 0.65, 1.0),
		},
		{
			name:     "hardlight_colors",
			mode:     "hardlight",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.85, 0.85, 0.85, 1.0),
		},
		{
			name:     "vividlight_colors",
			mode:     "vividlight",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.83, 0.83, 0.83, 1.0),
		},
		{
			name:     "linearlight_colors",
			mode:     "linearlight",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.9, 0.9, 0.9, 1.0),
		},
		{
			name:     "pinlight_colors",
			mode:     "pinlight",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
		},
		{
			name:     "hardmix_colors",
			mode:     "hardmix",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
		},
		{
			name:     "reflect_colors",
			mode:     "reflect",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.82, 0.82, 0.82, 1.0),
		},
		{
			name:     "glow_colors",
			mode:     "glow",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			expected: mustNewRGBA64(t, 0.82, 0.82, 0.82, 1.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, err := Get(tt.mode)
			if err != nil {
				t.Fatalf("failed to get blend mode %q: %v", tt.mode, err)
			}

			result, err := mode.Blend(tt.base, tt.blend, 1.0)
			if err != nil {
				t.Fatalf("failed to blend colors: %v", err)
			}
			if !validateRGBA64(result, tt.expected, 0.2) {
				t.Errorf("unexpected result for %s: got %+v, want %+v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestBlendModeAdvancedOperations tests advanced blend mode operations
func TestBlendModeAdvancedOperations(t *testing.T) {
	tests := []struct {
		name     string
		base     *color.RGBA64
		blend    *color.RGBA64
		alpha    float64
		wantErr  bool
		expected *color.RGBA64
	}{
		{
			name:     "partial_alpha",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			alpha:    0.5,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.75, 0.75, 0.75, 1.0),
		},
		// ... other test cases ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode("normal")
			if mode == nil {
				t.Fatal("normal blend mode not found")
			}

			result, err := mode.Blend(tt.base, tt.blend, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !validateRGBA64(result, tt.expected, 0.2) {
				t.Errorf("unexpected result: got %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// TestBlendModeComponentOperations tests component blend modes
func TestBlendModeComponentOperations(t *testing.T) {
	tests := []struct {
		name     string
		base     *color.RGBA64
		blend    *color.RGBA64
		alpha    float64
		wantErr  bool
		expected *color.RGBA64
	}{
		{
			name:     "hue_basic",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
		},
		{
			name:     "hue_grayscale_base",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
		},
		{
			name:     "hue_partial_alpha",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    0.5,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.65, 0.3, 0.45, 1.0),
		},
		{
			name:     "saturation_basic",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.7, 0.3, 0.4, 1.0),
		},
		{
			name:     "saturation_grayscale_base",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "saturation_partial_alpha",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    0.5,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.6, 0.25, 0.35, 1.0),
		},
		{
			name:     "luminosity_basic",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
		},
		{
			name:     "luminosity_grayscale",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.6, 0.6, 0.6, 1.0),
		},
		{
			name:     "luminosity_partial_alpha",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 1.0),
			alpha:    0.5,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.65, 0.3, 0.45, 1.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode(strings.Split(tt.name, "_")[0])
			if mode == nil {
				t.Fatalf("blend mode not found: %s", tt.name)
			}

			result, err := mode.Blend(tt.base, tt.blend, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !validateRGBA64(result, tt.expected, 0.2) {
				t.Errorf("unexpected result: got %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// TestBlendModeErrorHandling tests error cases for blend modes
func TestBlendModeErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		modeName string
		base     *color.RGBA64
		blend    *color.RGBA64
		alpha    float64
		wantErr  bool
	}{
		{
			name:     "nil_base_color",
			modeName: "normal",
			base:     nil,
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			alpha:    1.0,
			wantErr:  true,
		},
		{
			name:     "nil_blend_color",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    nil,
			alpha:    1.0,
			wantErr:  true,
		},
		{
			name:     "invalid_alpha_negative",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			alpha:    -0.1,
			wantErr:  true,
		},
		{
			name:     "invalid_alpha_above_one",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			alpha:    1.1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode(tt.modeName)
			if mode == nil {
				t.Fatalf("blend mode not found: %s", tt.modeName)
			}

			_, err := mode.Blend(tt.base, tt.blend, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestBlendModeSpecialCases tests special blend mode cases
func TestBlendModeSpecialCases(t *testing.T) {
	tests := []struct {
		name     string
		modeName string
		base     *color.RGBA64
		blend    *color.RGBA64
		alpha    float64
		wantErr  bool
		expected *color.RGBA64
	}{
		{
			name:     "zero_alpha_blend",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 0.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "half_alpha_blend",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 0.5),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.75, 0.75, 0.75, 1.0),
		},
		{
			name:     "half_alpha_operation",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			alpha:    0.5,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.75, 0.75, 0.75, 1.0),
		},
		{
			name:     "zero_alpha_base",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 0.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
		},
		{
			name:     "both_half_alpha",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 0.5),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 0.5),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.75, 0.75, 0.75, 0.75),
		},
		{
			name:     "multiply_zero_alpha_blend",
			modeName: "multiply",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.0, 0.0, 0.0, 0.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "screen_zero_alpha_blend",
			modeName: "screen",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 0.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "add_half_alpha_blend",
			modeName: "add",
			base:     mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.4, 0.4, 0.4, 0.5),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "overlay_half_alpha_operation",
			modeName: "overlay",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 1.0, 1.0, 1.0, 1.0),
			alpha:    0.5,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.75, 0.75, 0.75, 1.0),
		},
		{
			name:     "hue_zero_alpha_blend",
			modeName: "hue",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 0.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
		},
		{
			name:     "saturation_zero_alpha_blend",
			modeName: "saturation",
			base:     mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.4, 0.6, 0.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.2, 0.3, 1.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode(tt.modeName)
			if mode == nil {
				t.Fatalf("blend mode not found: %s", tt.modeName)
			}

			result, err := mode.Blend(tt.base, tt.blend, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !validateRGBA64(result, tt.expected, 0.2) {
				t.Errorf("unexpected result: got %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// TestBlendModeColorSpaceHandling tests color space conversions in blend modes
func TestBlendModeColorSpaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		modeName string
		base     *color.RGBA64
		blend    *color.RGBA64
		alpha    float64
		wantErr  bool
		expected *color.RGBA64
	}{
		{
			name:     "linear_space_multiply",
			modeName: "multiply",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.4, 0.4, 0.4, 1.0),
		},
		{
			name:     "linear_space_screen",
			modeName: "screen",
			base:     mustNewRGBA64(t, 0.2, 0.2, 0.2, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.8, 0.8, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.84, 0.84, 0.84, 1.0),
		},
		{
			name:     "linear_space_overlay",
			modeName: "overlay",
			base:     mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.42, 0.42, 0.42, 1.0),
		},
		{
			name:     "linear_space_add",
			modeName: "add",
			base:     mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.4, 0.4, 0.4, 1.0),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.7, 0.7, 0.7, 1.0),
		},
		{
			name:     "linear_space_normal_partial_alpha",
			modeName: "normal",
			base:     mustNewRGBA64(t, 0.3, 0.3, 0.3, 1.0),
			blend:    mustNewRGBA64(t, 0.7, 0.7, 0.7, 0.5),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
		},
		{
			name:     "linear_space_multiply_partial_alpha",
			modeName: "multiply",
			base:     mustNewRGBA64(t, 0.5, 0.5, 0.5, 1.0),
			blend:    mustNewRGBA64(t, 0.8, 0.8, 0.8, 0.5),
			alpha:    1.0,
			wantErr:  false,
			expected: mustNewRGBA64(t, 0.45, 0.45, 0.45, 1.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := GetBlendMode(tt.modeName)
			if mode == nil {
				t.Fatalf("blend mode not found: %s", tt.modeName)
			}

			result, err := mode.Blend(tt.base, tt.blend, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !validateRGBA64(result, tt.expected, 0.2) {
				t.Errorf("unexpected result: got %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// Benchmarks

func BenchmarkBlendModeBasic(b *testing.B) {
	base := mustNewRGBA64(b, 0.5, 0.5, 0.5, 1.0)
	blend := mustNewRGBA64(b, 0.8, 0.8, 0.8, 1.0)
	mode := GetBlendMode("normal")
	if mode == nil {
		b.Fatal("normal blend mode not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mode.Blend(base, blend, 1.0)
		if err != nil {
			b.Fatalf("failed to blend colors: %v", err)
		}
	}
}

func BenchmarkBlendModeWithAlpha(b *testing.B) {
	base := mustNewRGBA64(b, 0.5, 0.5, 0.5, 1.0)
	blend := mustNewRGBA64(b, 0.8, 0.8, 0.8, 0.5)
	mode := GetBlendMode("normal")
	if mode == nil {
		b.Fatal("normal blend mode not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mode.Blend(base, blend, 0.5)
		if err != nil {
			b.Fatalf("failed to blend colors: %v", err)
		}
	}
}

func BenchmarkBlendModeComplex(b *testing.B) {
	base := mustNewRGBA64(b, 0.3, 0.5, 0.7, 1.0)
	blend := mustNewRGBA64(b, 0.8, 0.2, 0.4, 0.5)
	mode := GetBlendMode("overlay")
	if mode == nil {
		b.Fatal("overlay blend mode not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mode.Blend(base, blend, 0.75)
		if err != nil {
			b.Fatalf("failed to blend colors: %v", err)
		}
	}
}
