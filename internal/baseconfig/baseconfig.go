package baseconfig

import (
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/yamlconfig"
)

type BaseConfig struct {
	cliargs.CLIArguments
	yamlconfig.YAMLConfig
}
