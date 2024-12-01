package main

import (
	"container/heap"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var fileName = "input/main.txt"

//var fileName = "input/test.txt"

func main() {
	var wg sync.WaitGroup

	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	b, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("failed to read file", err)
		os.Exit(-1)
	}

	text := string(b)
	exp, err := regexp.Compile("[0-9]+")
	if err != nil {
		slog.Error("failed to compile regex", err)
	}

	matches := exp.FindAllStringSubmatch(text, -1)
	wg.Add(2)
	go distance(matches, &wg)
	go similarity(matches, &wg)
	wg.Wait()
}

func similarity(matches [][]string, wg *sync.WaitGroup) {
	defer wg.Done()

	counts := [2]map[int]int{{}, {}}
	for i := 0; i < len(matches); i++ {
		number, err := strconv.Atoi(matches[i][0])
		if err != nil {
			slog.Error("failed to parse int", err)
			return
		}
		if val, ok := counts[i%2][number]; ok {
			counts[i%2][number] = val + 1
		} else {
			counts[i%2][number] = 1
		}
	}
	slog.Debug("similarity maps",
		slog.Any("left", counts[0]),
		slog.Any("right", counts[1]),
	)

	result := 0
	for k, v1 := range counts[0] {
		if v2, ok := counts[1][k]; ok {
			slog.Debug("adding product",
				slog.Int("number", k),
				slog.Int("left", v1),
				slog.Int("right", v2),
			)
			result += k * v1 * v2
		}
	}
	slog.Info("Similarity", slog.Int("result", result))
}

func distance(matches [][]string, wg *sync.WaitGroup) {
	defer wg.Done()
	heaps := [2]*IntHeap{{}, {}}

	heap.Init(heaps[0])
	heap.Init(heaps[1])
	for i := 0; i < len(matches); i++ {
		number, err := strconv.Atoi(matches[i][0])
		if err != nil {
			slog.Error("failed to parse int", err)
			return
		}
		heap.Push(heaps[i%2], number)
	}

	ans := 0
	for heaps[0].Len() > 0 {
		left, right := heap.Pop(heaps[0]).(int), heap.Pop(heaps[1]).(int)
		if left > right {
			ans += left - right
		} else {
			ans += right - left
		}
	}
	slog.Info("distance", slog.Int("result", ans))
}
