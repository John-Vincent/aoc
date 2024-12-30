package main

import (
	"log/slog"
	"os"
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

	text, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("failed to read file", slog.Any("error", err))
		os.Exit(-1)
	}

	slog.Debug("\n" + string(text))
	basicCount := countXmas(text)
	x_masCount := countX_Mas(text)

	slog.Info("found", slog.Int("xmas's", basicCount))
	slog.Info("found", slog.Int("x_mas's", x_masCount))
}

func findNewline(board []byte) int {
	var newLine = byte('\n')
	for i := 0; i < len(board); i++ {
		if board[i] == newLine {
			return i
		}
	}
	return -1
}

func countX_Mas(board []byte) int {
	var count = 0
	characters := [3]byte{'M', 'A', 'S'}
	corners := [4]int{}
	lineLength := findNewline(board)

outer:
	for i := 0; i < len(board); i++ {
		if board[i] == characters[1] {
			corners[0] = i - (lineLength + 2)
			corners[1] = i + (lineLength + 2)
			corners[2] = i - lineLength
			corners[3] = i + lineLength
			for j := 0; j < 4; j++ {
				if corners[j] < 0 || corners[j] >= len(board) {
					continue outer
				}
				if board[corners[j]] != characters[0] && board[corners[j]] != characters[2] {
					continue outer
				}
			}
			if board[corners[0]] != board[corners[1]] && board[corners[2]] != board[corners[3]] {
				slog.Debug("found X_MAS",
					slog.Any("middle", i),
					slog.Any("corners", corners),
				)
				count++
			}
		}
	}
	return count
}

func countXmas(board []byte) int {
	var count = 0
	var increments = [8]int{1, -1}
	characters := []byte{'X', 'M', 'A', 'S'}
	var lineLength = findNewline(board)

	increments[2] = lineLength
	increments[3] = lineLength + 1
	increments[4] = lineLength + 2
	increments[5] = -increments[2]
	increments[6] = -increments[3]
	increments[7] = -increments[4]

	slog.Debug("stats",
		slog.Any("increments", increments),
		slog.Any("4,0", string(board[3])),
		slog.Any("4,1", string(board[3+increments[3]])),
		slog.Any("4,2", string(board[3+2*increments[3]])),
		slog.Any("4,3", string(board[3+3*increments[3]])),
	)

	for i := 0; i < len(board); i++ {
		if board[i] == characters[0] {
			for j := 0; j < len(increments); j++ {
				if test(board, i, increments[j], characters) {
					slog.Debug("found xmas",
						slog.Int("start", i),
						slog.Int("direction", increments[j]),
					)
					count++
				}
			}
		}
	}

	return count
}

func test(board []byte, index int, increment int, characters []byte) bool {
	var current = index
	for i := 1; i < len(characters); i++ {
		current += increment
		if current < 0 || current >= len(board) || board[current] != characters[i] {
			return false
		}
	}
	return true
}
