package main

import (
	"log/slog"
	"os"

	"github.com/opentdf/otdfctl/cmd"
)

func main() {
	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	l := new(slog.LevelVar)
	l.Set(slog.LevelInfo)
	l.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))) //nolint:errcheck // ignore error, just use default level
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	}))

	slog.SetDefault(logger)

	cmd.Execute()
}
