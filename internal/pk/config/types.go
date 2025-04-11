package config

// Config represents the top-level structure of the YAML configuration
type Config struct {
	Common    Common     `yaml:"common"`
	Templates []Template `yaml:"templates"`
}

// Common represents the common configuration settings
type Common struct {
	Repo             string   `yaml:"repo"`
	Path             string   `yaml:"path"`
	NonInteractive   bool     `yaml:"non_interactive"`
	Ref              string   `yaml:"ref"`
	VarFiles         []string `yaml:"var_files"`
	BaseOutputFolder string   `yaml:"base_output_folder"`
}

// Template represents an individual template configuration
type Template struct {
	Name           string   `yaml:"name"`
	VarFiles       []string `yaml:"var_files,omitempty"`
	NonInteractive *bool    `yaml:"non_interactive,omitempty"`
	Ref            string   `yaml:"ref"`
	Subfolder      string   `yaml:"subfolder"`
}
