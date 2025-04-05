package projections

import (
	"testing"
)

func TestCoordinate(t *testing.T) {
	// Test NewCoordinate
	coord := NewCoordinateMeta("test", "test description", "meters", 100.0, 0.0)
	if coord.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", coord.Name())
	}
	if coord.Description() != "test description" {
		t.Errorf("Expected description 'test description', got '%s'", coord.Description())
	}
	if coord.Unit() != "meters" {
		t.Errorf("Expected unit 'meters', got '%s'", coord.Unit())
	}
	if coord.Min() != 0.0 {
		t.Errorf("Expected min 0.0, got %f", coord.Min())
	}
	if coord.Max() != 100.0 {
		t.Errorf("Expected max 100.0, got %f", coord.Max())
	}

	// Test backward compatibility
	if coord.Desc() != coord.Description() {
		t.Errorf("Desc() should return the same as Description()")
	}
}

func TestProjectionMeta(t *testing.T) {
	// Create test coordinates
	coordA := NewCoordinateMeta("x", "x coordinate", "meters", 100.0, 0.0)
	coordB := NewCoordinateMeta("y", "y coordinate", "meters", 100.0, 0.0)

	// Test NewProjectionMeta
	meta := NewProjectionMeta("test projection", "test projection description", coordA, coordB)
	if meta.Name() != "test projection" {
		t.Errorf("Expected name 'test projection', got '%s'", meta.Name())
	}
	if meta.Description() != "test projection description" {
		t.Errorf("Expected description 'test projection description', got '%s'", meta.Description())
	}

	// Test coordinate accessors
	if meta.NameA() != "x" {
		t.Errorf("Expected NameA() 'x', got '%s'", meta.NameA())
	}
	if meta.DescriptionA() != "x coordinate" {
		t.Errorf("Expected DescriptionA() 'x coordinate', got '%s'", meta.DescriptionA())
	}
	if meta.UnitA() != "meters" {
		t.Errorf("Expected UnitA() 'meters', got '%s'", meta.UnitA())
	}
	if meta.MinA() != 0.0 {
		t.Errorf("Expected MinA() 0.0, got %f", meta.MinA())
	}
	if meta.MaxA() != 100.0 {
		t.Errorf("Expected MaxA() 100.0, got %f", meta.MaxA())
	}

	// Test backward compatibility
	if meta.Desc() != meta.Description() {
		t.Errorf("Desc() should return the same as Description()")
	}
	if meta.DescA() != meta.DescriptionA() {
		t.Errorf("DescA() should return the same as DescriptionA()")
	}
	if meta.DescB() != meta.DescriptionB() {
		t.Errorf("DescB() should return the same as DescriptionB()")
	}

	// Test Doc() method
	doc := meta.Doc()
	if doc == "" {
		t.Error("Doc() returned empty string")
	}
	// Basic checks for expected content
	if !contains(doc, "test projection") {
		t.Error("Doc() should contain projection name")
	}
	if !contains(doc, "test projection description") {
		t.Error("Doc() should contain projection description")
	}
	if !contains(doc, "x") || !contains(doc, "y") {
		t.Error("Doc() should contain coordinate names")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
