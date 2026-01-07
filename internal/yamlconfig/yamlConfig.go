package yamlconfig

type YAMLConfig struct {
	MaxIterations int    `yaml:"max_iterations"`
	MaxFilesize   int    `yaml:"max_filesize"`
	WorkingDir    string `yaml:"working_dir"`
}
