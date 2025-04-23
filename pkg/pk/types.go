package pk

type Template struct {
	Repo             string   `yaml:"repo"`
	Path             string   `yaml:"path"`
	NonInteractive   bool     `yaml:"non_interactive"`
	Ref              string   `yaml:"ref"`
	VarFiles         []string `yaml:"var_files"`
	BaseOutputFolder string   `yaml:"base_output_folder"`
	Name             string   `yaml:"name"`
	Subfolder        string   `yaml:"subfolder"`
}

type Config struct {
	Common    Template   `yaml:"common"`
	Templates []Template `yaml:"templates"`
}
