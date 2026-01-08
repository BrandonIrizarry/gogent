package functions

import (
	"fmt"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	workingDir := "~/repos/gogent"

	tests := []struct {
		relpath, want string
	}{
		{".env", "~/repos/gogent/.env"},
		{".", "~/repos/gogent"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("Test %s", tt.relpath)
		t.Run(name, func(t *testing.T) {
			ans, err := normalizePath(tt.relpath, workingDir)
			if err != nil {
				t.Fatal(err)
			}

			if ans != tt.want {
				t.Errorf("Got %s, want %s", ans, tt.want)
			}

		})
	}
}
