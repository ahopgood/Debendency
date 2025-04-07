package pkg_test

import (
	"bytes"
	"com/alexander/debendency/pkg"
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"log/slog"
	"strings"
)

var _ = Describe("ConfigureLogger", func() {

	When("logger built", func() {
		var buf bytes.Buffer

		BeforeEach(func() {
			buf.Reset()
			slog.SetLogLoggerLevel(slog.LevelInfo)
			log.SetOutput(&buf)
		})

		It("should log dates", func() {
			pkg.ConfigureLogger(pkg.Config{})
			//2024/09/22 14:16:52 DEBUG test message
			slog.Info("test message")
			actualDate := strings.Fields(buf.String())[0]

			Expect(actualDate).To(MatchRegexp("[[:digit:]]{4}/[[:digit:]]{2}/[[:digit:]]{2}"))
		})

		It("should log timestamps in microseconds", func() {
			pkg.ConfigureLogger(pkg.Config{})

			//2024/09/24 20:14:43.812102 INFO test message
			slog.Info("test message")
			fmt.Println(buf.String())
			actualTimeStamp := strings.Fields(buf.String())[1]
			Expect(actualTimeStamp).To(MatchRegexp("^[[:digit:]]{2}:[[:digit:]]{2}:[[:digit:]]{2}\\.[[:digit:]]{6}$"))
		})

		It("should log long code locations", func() {
			pkg.ConfigureLogger(pkg.Config{})

			//2024/09/28 12:25:10.480752 C:/Users/Alex/IdeaProjects/Debendency/pkg/logger_test.go:49: INFO test message
			slog.Info("test message")
			fmt.Println(buf.String())
			actualCodeLocation := strings.Fields(buf.String())[2]
			Expect(actualCodeLocation).To(MatchRegexp("/pkg/logger_test\\.go:[[:digit:]]+:"))
		})

		It("should not be structured and in json", func() {
			pkg.ConfigureLogger(pkg.Config{})

			//2024/09/22 14:16:52 DEBUG test message
			slog.Info("test message")
			fmt.Println(buf.String())
			actualLogMessage := strings.Fields(buf.String())[4:6]
			Expect(actualLogMessage).To(ContainElements("test", "message"))
		})

		It("should not be structured and in text", func() {
			pkg.ConfigureLogger(pkg.Config{})

			//2024/09/22 14:16:52 DEBUG test message
			slog.Info("test message")
			fmt.Println(buf.String())
			actualLogMessage := strings.Fields(buf.String())[4:6]
			Expect(actualLogMessage).To(ContainElements("test", "message"))
		})

		It("log at debug if Config.Verbose is true", func() {
			pkg.ConfigureLogger(pkg.Config{Verbose: true})

			//2024/09/22 14:16:52 DEBUG test message
			slog.Debug("test message")
			fmt.Println(buf.String())
			actualLogLevel := strings.Fields(buf.String())[3]
			Expect(actualLogLevel).To(ContainSubstring(slog.LevelDebug.String()))
		})

		DescribeTable("should default to Info level",
			func(level slog.Level, expectedLevel string, shouldLog bool) {
				pkg.ConfigureLogger(pkg.Config{})

				slog.Log(context.Background(), level, "test message")
				//outputs:
				//2024/09/22 14:16:52 DEBUG test message
				fmt.Println(buf.String())
				if shouldLog {
					actualLogLevel := strings.Fields(buf.String())[3]
					Expect(actualLogLevel).To(ContainSubstring(expectedLevel))
				} else {
					Expect(buf.String()).To(BeEmpty())
				}
			},
			Entry(
				fmt.Sprintf("Log level set to %d", slog.LevelDebug), slog.LevelDebug, "DEBUG", false),
			Entry(
				fmt.Sprintf("Log level set to %d", slog.LevelInfo), slog.LevelInfo, "INFO", true),
			Entry(
				fmt.Sprintf("Log level set to %d", slog.LevelWarn), slog.LevelWarn, "WARN", true),
			Entry(
				fmt.Sprintf("Log level set to %d", slog.LevelError), slog.LevelError, "ERROR", true),
		)
	})
})
