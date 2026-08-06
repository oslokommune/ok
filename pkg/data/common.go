package data

const (
	// DefaultTemplateURL is the default Oslo kommune databricks template
	DefaultTemplateURL = "git@github.com:oslokommune/padda-golden-path"

	// DefaultTemplateDir is the directory within the default template
	// repository that contains the databricks template
	DefaultTemplateDir = "bundle-templates"

	// TemplateURLEnvName allows overriding the default template
	TemplateURLEnvName = "DATA_TEMPLATE_URL"
)
