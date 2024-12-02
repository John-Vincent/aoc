package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
)

var fileName = "input/main.txt"

func main() {
	var wg sync.WaitGroup

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
	reports := strings.Split(text, "\n")

	wg.Add(1)
	//go safe(reports, &wg)
	go safeWithForgiveness(reports, &wg)
	wg.Wait()
}

func safe(reports []string, wg *sync.WaitGroup) {
	defer wg.Done()

	var result int = 0
	for _, element := range reports {
		digits := strings.Split(element, " ")
		if listSafe(digits) == -1 {
			result++
		}
	}

	slog.Info("", slog.Int("safe_reports", result))
}

/*
performance of this is horrible and memory usage as well

but it was faster to call the existing function with clones
than to modify it to test by skipping indexes
*/
func safeWithForgiveness(reports []string, wg *sync.WaitGroup) {
	defer wg.Done()

	var result int = 0
	for _, element := range reports {
		digits := strings.Split(element, " ")
		index := listSafe(digits)
		if index == -1 {
			result++
			continue
		}

		//most common case
		clone := make([]string, len(digits))
		copy(clone, digits)
		removeIndex := append(clone[:index], clone[index+1:]...)
		if listSafe(removeIndex) == -1 {
			result++
			continue
		}

		if index < 4 {
			slog.Debug("digits after append", slog.Any("digits", digits))
			removeIndex = digits[1:]
			if listSafe(removeIndex) == -1 {
				result++
				continue
			}
			removeIndex = append(digits[:1], digits[2:]...)
			if listSafe(removeIndex) == -1 {
				result++
			}
		}
	}

	slog.Info("with forgiveness", slog.Int("safe_reports", result))
}

func listSafe(digits []string) int {
	var ascending bool
	listLog := slog.Any("list", digits)

	last, err := strconv.Atoi(digits[0])
	if err != nil {
		slog.Error("failed to parse int", err)
		return -2
	}
	for i := 1; i < len(digits); i++ {
		current, err := strconv.Atoi(digits[i])
		if err != nil {
			slog.Error("failed to parse int", err)
			return -2
		}
		magnitude := last - current
		if magnitude < 0 {
			magnitude = -magnitude
		}

		if i == 1 {
			ascending = last > current
		} else if ascending != (last > current) {
			slog.Debug("direction is wrong",
				listLog,
				slog.Int("last", last),
				slog.Int("current", current),
			)
			return i
		}

		if magnitude == 0 || magnitude > 3 {
			slog.Debug("magnitude is wrong",
				listLog,
				slog.Int("last", last),
				slog.Int("current", current),
			)
			return i
		}
		last = current
	}

	slog.Debug("passed",
		listLog,
	)
	return -1
}
