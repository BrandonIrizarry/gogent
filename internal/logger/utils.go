package logger

func satisfies(verbosity, mask LogMode) bool {
	return verbosity&mask == mask
}
