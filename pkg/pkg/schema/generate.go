package schema

import (
	"context"
	"fmt"
	"github.com/oslokommune/ok/pkg/jsonschema"
	"github.com/oslokommune/ok/pkg/pkg/config"
)

func GenerateJsonSchemaForApp(ctx context.Context, downloader config.FileDownloader, stackPath, gitRef string) (*jsonschema.Document, error) {
	stacks, err := config.DownloadBoilerplateTemplatesWithDependencies(ctx, downloader, stackPath)
	if err != nil {
		return nil, fmt.Errorf("downloading boilerplate stacks: %w", err)
	}

	mobules := BuildModuleVariables(stacks)

	schema, err := TransformModulesToJsonSchema(fmt.Sprintf("%s-%s", stackPath, gitRef), mobules)
	if err != nil {
		return nil, fmt.Errorf("transforming modules to json schema: %w", err)
	}

	return schema, nil
}
