package install

import (
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestFindPackagesFromOutputFolders(t *testing.T) {
	packages := []common.Package{
		{OutputFolder: "out/app-hello"},
		{OutputFolder: "out/networking"},
	}

	tests := []struct {
		name          string
		outputFolders []string
		expectedPkgs  []common.Package
	}{
		{
			name:          "Package found",
			outputFolders: []string{"out/app-hello"},
			expectedPkgs:  []common.Package{packages[0]},
		},
		{
			name:          "Package not found",
			outputFolders: []string{"out/unknown"},
			expectedPkgs:  []common.Package{},
		},
		{
			name:          "Multiple output folders, package found",
			outputFolders: []string{"out/unknown", "out/networking"},
			expectedPkgs:  []common.Package{packages[1]},
		},
		{
			name:          "Multiple output folders, package not found",
			outputFolders: []string{"out/unknown1", "out/unknown2"},
			expectedPkgs:  []common.Package{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := FindPackagesFromOutputFolders(packages, tt.outputFolders)
			assert.Equal(t, tt.expectedPkgs, pkgs)
		})
	}
}
