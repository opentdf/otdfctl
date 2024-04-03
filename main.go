package main

import (
	"log/slog"
	"os"

	"github.com/opentdf/tructl/cmd"
)

func main() {
	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	l := new(slog.LevelVar)
	l.Set(slog.LevelDebug)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	}))

	slog.SetDefault(logger)

	cmd.Execute()
}
