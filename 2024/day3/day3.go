package main

import (
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var fileName = "input/main.txt"

func main() {

	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	file := strings.ToLower(os.Getenv("TEST_FILE"))

	if len(file) > 0 {
		fileName = "input/test.txt"
	}

	b, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("failed to read file", err)
		os.Exit(-1)
	}

	text := strings.TrimSpace(string(b))
	slog.Debug(text)

	exp := regexp.MustCompile(`(?:mul\(([0-9]{0,3}),([0-9]{0,3})\))|(?:do\(\))|(?:don't\(\))`)
	matches := exp.FindAllStringSubmatch(text, -1)

	slog.Debug("matches", slog.Any("matches", matches))
	sum(matches)
}

func sum(matches [][]string) {
	basicSum, sum := 0, 0
	enabled := true
	for _, op := range matches {
		if op[0] == "do()" {
			enabled = true
			continue
		}
		if op[0] == "don't()" {
			enabled = false
			continue
		}
		left, err := strconv.Atoi(op[1])
		if err != nil {
			slog.Error("failed to parse number", err)
		}
		right, err := strconv.Atoi(op[2])
		if err != nil {
			slog.Error("failed to parse number", err)
		}
		basicSum += left * right
		if enabled {
			sum += left * right
		}
	}
	slog.Info("basic sum of products", slog.Int("result", basicSum))
	slog.Info("sensitive sum of products", slog.Int("result", sum))
}
