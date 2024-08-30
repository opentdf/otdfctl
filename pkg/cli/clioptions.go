package cli

type cliOptions struct {
	printerJson bool
}

type cliVariadicOption func(cliOptions) cliOptions

// WithPrintJson is a variadic option that enforces JSON output for the printer
func WithPrintJson() cliVariadicOption {
	return func(o cliOptions) cliOptions {
		o.printerJson = true
		return o
	}
}
