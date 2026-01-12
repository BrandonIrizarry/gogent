package baseconfig

type BaseConfig struct {
	// We're now getting this from a TUI file picker. Eventually,
	// this app will become a full TUI app and so all of its
	// config will come from TUI widgets.
	WorkingDir  string
	MaxFilesize int
}
