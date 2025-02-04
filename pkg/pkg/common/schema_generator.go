package common

import "context"

type SchemaGenerator interface {
	CreateJsonSchemaFile(
		ctx context.Context, manifestPackagePrefix string, pkg Package) ([]byte, error)
}
