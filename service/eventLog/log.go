package eventLog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func InitLog() {
	zerolog.SetGlobalLevel(-1)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func NewLog() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	// output.FormatLevel = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	// }
	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		switch level {
		case "| TRACE |", "| DEBUG |":
			return "\x1b[36m" + level + "\x1b[0m" // Cyan
		case "| INFO  |":
			return "\x1b[32m" + level + "\x1b[0m" // Green
		case "| WARN  |":
			return "\x1b[33m" + level + "\x1b[0m" // Yellow
		case "| ERROR |", "| FATAL |", "| PANIC |":
			return "\x1b[31m" + level + "\x1b[0m" // Red
		default:
			return level
		}
	}
	// output.FormatMessage = func(i interface{}) string {
	// 	return fmt.Sprintf("%s", i)
	// }
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	// output.FormatCaller = func(i interface{}) string {
	// 	return fmt.Sprintf("%s", i)
	// }
	log := zerolog.New(output).With().Timestamp().Caller().Logger()
	return log
}
