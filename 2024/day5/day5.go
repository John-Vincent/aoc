package main

import (
	"log/slog"
	"os"
	"strings"
)

var fileName = "input/main.txt"

type input struct {
	rules map[int][]int
	edits [][]int
}

func main() {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	file := strings.ToLower(os.Getenv("TEST_FILE"))

	if len(file) > 0 {
		fileName = "input/test.txt"
	}

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("failed to read file", slog.Any("error", err))
		os.Exit(-1)
	}

	text := string(bytes)

	slog.Debug("\n" + text)
	input := parseInput(bytes)
	slog.Debug("parsed input", slog.Any("value", *input))
	testEdits(input)
}

func testEdits(input *input) {
	slog.Debug("test edits", slog.Any("edits", input.edits))

	for edit := range input.edits {
		for page := range edit {

		}
	}
}

func parseInput(data []byte) *input {
	var result = input{rules: make(map[int][]int), edits: make([][]int, 0, 1000)}
	var first, second, i int
	pipe, doubleNewline := false, false

	for i = 0; i < len(data) && !doubleNewline; i++ {
		switch char := data[i]; char {
		case '|':
			pipe = true
		case '\n':
			if pipe == false {
				doubleNewline = true
				continue
			}
			result.rules[second] = append(result.rules[second], first)
			pipe = false
			first = 0
			second = 0
		default:
			if pipe {
				second = second*10 + int(char-'0')
			} else {
				first = first*10 + int(char-'0')
			}
		}
	}

	result.edits = append(result.edits, make([]int, 0, 100))
	for first, second = 0, 0; i < len(data); i++ {
		switch data[i] {
		case ',':
			result.edits[second] = append(result.edits[second], first)
			first = 0
		case '\n':
			result.edits[second] = append(result.edits[second], first)
			second++
			result.edits = append(result.edits, make([]int, 0, 100))
			first = 0
		default:
			first = first*10 + int(data[i]-'0')
		}
	}

	result.edits = result.edits[:len(result.edits)-1]
	return &result
}
