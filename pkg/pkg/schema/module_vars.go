package schema

import (
	"github.com/oslokommune/ok/pkg/pkg/config"
	"log"
	"strings"
)

type Stack struct {
	Name         string
	Config       *config.BoilerplateConfig
	OutputFolder string
	Dependencies []string
}

type ModuleVariables struct {
	Namespace string
	Variables []config.BoilerplateVariable
}

type CombinedVariables struct {
	OutputFolder string
	Namespace    string
	Variables    []config.BoilerplateVariable
}

func BuildModuleVariables(configs []*config.BoilerplateTemplate) []*ModuleVariables {
	if len(configs) == 0 {
		return nil
	}

	return buildModuleVariables("", configs[0], configs, "some/output/folder")
}

func buildModuleVariables(namespace string, currentConfig *config.BoilerplateTemplate, configs []*config.BoilerplateTemplate, outputFolder string) []*ModuleVariables {
	// ensure input arguments follow the correct format to avoid creating invalid namespaces
	namespace = JoinNamespaces(namespace)
	outputFolder = config.JoinPath(outputFolder, currentConfig.Path)

	namespaceVariables := make(map[string][]config.BoilerplateVariable)
	namespaceVariables[namespace] = currentConfig.Config.Variables

	for _, dep := range currentConfig.Config.Dependencies {
		depPath := config.JoinPath(currentConfig.Path, dep.TemplateUrl)
		depConfig, ok := findConfigFromPath(depPath, configs)
		if !ok {
			log.Printf("dependency %s not found in configs referenced by %s", depPath, currentConfig.Path)
			continue
		}

		depNamespace := namespace
		depOutputFolder := config.JoinPath(outputFolder, dep.OutputFolder)
		// if we move to a different output folder, then we need to create a new namespace
		if depOutputFolder != outputFolder {
			depNamespace = JoinNamespaces(depNamespace, dep.Name)
		}
		subModuleVariables := buildModuleVariables(depNamespace, depConfig, configs, depOutputFolder)
		for _, m := range subModuleVariables {
			namespaceVariables[m.Namespace] = append(namespaceVariables[m.Namespace], m.Variables...)
		}
	}

	moduleVariables := make([]*ModuleVariables, 0, len(namespaceVariables))
	for namespace, variables := range namespaceVariables {
		moduleVariables = append(moduleVariables, &ModuleVariables{
			Namespace: namespace,
			Variables: variables,
		})
	}

	return moduleVariables
}

func findConfigFromPath(path string, configs []*config.BoilerplateTemplate) (*config.BoilerplateTemplate, bool) {
	for _, c := range configs {
		if c.Path == path {
			return c, true
		}
	}
	return nil, false
}

func JoinNamespaces(namespaces ...string) string {
	filtered := make([]string, 0)
	for _, n := range namespaces {
		if n != "" {
			filtered = append(filtered, n)
		}
	}

	return strings.Join(filtered, ".")
}
