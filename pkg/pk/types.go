package pk

type CommonConfig struct {
	Repo             string   `yaml:"repo"`
	Path             string   `yaml:"path"`
	NonInteractive   bool     `yaml:"non_interactive"`
	Ref              string   `yaml:"ref"`
	VarFiles         []string `yaml:"var_files"`
	BaseOutputFolder string   `yaml:"base_output_folder"`
}

type Template struct {
	Name           string   `yaml:"name"`
	VarFiles       []string `yaml:"var_files,omitempty"`
	NonInteractive bool     `yaml:"non_interactive,omitempty"`
	Ref            string   `yaml:"ref"`
	Subfolder      string   `yaml:"subfolder"`
}

type Config struct {
	Common    CommonConfig `yaml:"common"`
	Templates []Template   `yaml:"templates"`
}
