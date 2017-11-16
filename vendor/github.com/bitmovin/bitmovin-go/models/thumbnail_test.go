package models

import "testing"

func TestBuilder(t *testing.T) {
	thumb := NewThumbnail(300, []float64{1, 5, 30}, []Output{Output{}})
	thumb = thumb.Builder().
		Name("Test Thumbnail").
		Build()

	if *thumb.Name != "Test Thumbnail" {
		t.Error("Wanted Thumbnail Name to be `Test Thumbnail` got %s", *thumb.Name)
	}

	if thumb.Description != nil {
		t.Error("Wanted Thumbnail Description to be nil, got %v", thumb.Description)
	}

	// Test that it's manipulating references
	thumb.Builder().Description("My Desc")

	if *thumb.Description != "My Desc" {
		t.Error("Wanted Thumbnail Description to be `My Desc`, got %v", thumb.Description)
	}
}
