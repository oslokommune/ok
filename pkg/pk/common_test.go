package pk

import (
	"testing"
)

func TestFilterTemplatesByWorkingDir(t *testing.T) {
	templates := []Template{
		{BaseOutputFolder: ".", Subfolder: "app-hello", Name: "app"},
		{BaseOutputFolder: ".", Subfolder: "networking", Name: "networking"},
		{BaseOutputFolder: "infra", Subfolder: "database", Name: "db"},
	}

	tests := []struct {
		name     string
		cwd      string
		repoRoot string
		want     []string // expected subfolder names
	}{
		{
			name:     "at repo root returns nil",
			cwd:      "/repo",
			repoRoot: "/repo",
			want:     nil,
		},
		{
			name:     "exact match on subfolder",
			cwd:      "/repo/app-hello",
			repoRoot: "/repo",
			want:     []string{"app-hello"},
		},
		{
			name:     "within subfolder",
			cwd:      "/repo/app-hello/src",
			repoRoot: "/repo",
			want:     []string{"app-hello"},
		},
		{
			name:     "nested base output folder",
			cwd:      "/repo/infra/database",
			repoRoot: "/repo",
			want:     []string{"database"},
		},
		{
			name:     "within nested output folder",
			cwd:      "/repo/infra/database/migrations",
			repoRoot: "/repo",
			want:     []string{"database"},
		},
		{
			name:     "no match",
			cwd:      "/repo/other-folder",
			repoRoot: "/repo",
			want:     nil,
		},
		{
			name:     "partial name should not match",
			cwd:      "/repo/app-hello-world",
			repoRoot: "/repo",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterTemplatesByWorkingDir(templates, tt.cwd, tt.repoRoot)

			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("expected %d templates, got %d", len(tt.want), len(got))
				return
			}

			for i, tpl := range got {
				if tpl.Subfolder != tt.want[i] {
					t.Errorf("expected subfolder %q, got %q", tt.want[i], tpl.Subfolder)
				}
			}
		})
	}
}
