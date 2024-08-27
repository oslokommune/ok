package config

import (
	"fmt"
	"net/url"

	"github.com/oslokommune/ok/pkg/jsonschema"
)

func BuildJsonSchemaFromConfig(config *BoilerplateConfig, dependencies []BoilerplateConfig) (*jsonschema.Document, error) {
	return nil, fmt.Errorf("not implemented")
}

type Stack struct {
	Name         string
	Config       *BoilerplateConfig
	OutputFolder string
	Dependencies []string
}

func transformConfigsToStacks(rootConfig *BoilerplateConfig, packageConfigs []BoilerplateConfig) ([]*Stack, error) {

	folderStacks := make(map[string][]string)

	rootConfig.Dependencies

	for _, dep := range rootConfig.Dependencies {

	}
	return nil, fmt.Errorf("not implemented")
}

func mustJoinUri(base, path string) string {
	uri, err := url.JoinPath(base, path)
	if err != nil {
		panic(err)
	}
	return uri
}
