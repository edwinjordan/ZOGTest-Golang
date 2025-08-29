package utils

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	got := Slugify("Hello World")
	want := "hello-world"

	if got != want {
		t.Errorf("Slugify() = %v, want %v", got, want)
	}
}
