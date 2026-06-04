package lang

import (
	"testing"
	"testing/fstest"
)

func TestPackageDefault(t *testing.T) {
	// No default set yet: helpers must echo the key, never panic.
	SetDefault(nil)
	if got := Trans("user.greeting"); got != "user.greeting" {
		t.Errorf("Trans before SetDefault = %q, want echoed key", got)
	}
	if got := TransBy("fr", "user.greeting"); got != "user.greeting" {
		t.Errorf("TransBy before SetDefault = %q, want echoed key", got)
	}

	tr, err := New(
		fstest.MapFS{
			"en.yaml": &fstest.MapFile{Data: []byte("user:\n  greeting: \"Hello {name}\"\n")},
		},
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	SetDefault(tr)
	t.Cleanup(func() { SetDefault(nil) })

	if got := Trans("user.greeting", map[string]string{"name": "Sam"}); got != "Hello Sam" {
		t.Errorf("Trans = %q, want %q", got, "Hello Sam")
	}
	if got := TransBy("en", "user.greeting", map[string]string{"name": "Sam"}); got != "Hello Sam" {
		t.Errorf("TransBy = %q, want %q", got, "Hello Sam")
	}
}
