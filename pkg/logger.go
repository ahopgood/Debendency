package pkg

import (
	"fmt"
	"log"
	"log/slog"
)

func ConfigureLogger(config Config) {
	if config.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		fmt.Println("Have set logger to debug")
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
	//log.LUTC
	//slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
}
