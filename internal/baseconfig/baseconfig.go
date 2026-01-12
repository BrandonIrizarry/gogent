package baseconfig

import (
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/yamlconfig"
)

type BaseConfig struct {
	cliargs.CLIArguments
	yamlconfig.YAMLConfig

	// We're now getting this from a TUI file picker. Eventually,
	// this app will become a full TUI app and so all of its
	// config will come from TUI widgets.
	WorkingDir string
}
