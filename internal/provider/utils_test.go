package provider

import "testing"

func TestContains(t *testing.T) {
	if !contains([]string{"foo", "bar", "baz"}, "bar") {
		t.Fatal("exepected to find \"bar\" in [foo, bar, baz]")
	}

	if contains([]string{"foo", "bar", "baz"}, "foobar") {
		t.Fatal("expected not to find \"foobar\" in [foo, bar, baz]")
	}
}
