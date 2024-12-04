package composition

import (
	"testing"
)

func TestComposition_Load(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"1", "../test_data/compositions/1.yaml"},
		{"2", "../test_data/compositions/2.yaml"},
		{"3", "../test_data/compositions/3.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewFromFile(tt.path).Render().SaveAsPNG("../test_data/compositions/render/" + tt.name + ".png")
		})
	}
}
