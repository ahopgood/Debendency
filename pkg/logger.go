package pkg

import (
	"log"
	"log/slog"
)

func ConfigureLogger(config Config) {
	if config.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	//log.LUTC
	//slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
}
