package install

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilterPackages(t *testing.T) {
	tests := []struct {
		name          string
		packages      []common.Package
		outputFolders []string
		expected      []common.Package
	}{
		{
			name: "no output folders specified",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
			},
			outputFolders: []string{},
			expected:      []common.Package{},
		},
		{
			name: "single output folder specified",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
			},
			outputFolders: []string{"out/folder1"},
			expected: []common.Package{
				{OutputFolder: "out/folder1"},
			},
		},
		{
			name: "multiple output folders specified",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
				{OutputFolder: "out/folder3"},
			},
			outputFolders: []string{"out/folder1", "out/folder3"},
			expected: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder3"},
			},
		},
		{
			name: "no matching output folders",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
			},
			outputFolders: []string{"out/folder3"},
			expected:      []common.Package{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FindPackageFromOutputFolders(tc.packages, tc.outputFolders)
			require.Equal(t, tc.expected, result)
		})
	}
}
