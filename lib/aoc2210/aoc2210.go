package aoc2210

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

const (
	addx instruction = iota
	noop
)

func SumSignalStrengths(lines []string) (int, []string) {
	registerX := 1
	nextCheckpoint := 20
	checkpointInterval := 40
	frames := parseLines(lines)
	crtWidth, crtHeight := 40, 6
	crt := make([][]rune, crtHeight)
	for i := range crt {
		crt[i] = make([]rune, crtWidth)
	}
	var frameIndex int
	var strentghSum int
	nextActiveCycle := 1
	stack := map[int]int{}
	for cycle := 1; ; cycle++ {
		crt[(cycle-1)/crtWidth][(cycle-1)%crtWidth] = deriveCrtPixel(
			cycle,
			crtWidth,
			registerX,
		)

		if cycle == nextCheckpoint {
			strength := cycle * registerX
			strentghSum += strength
			shared.Logger.Info(
				"Register strength.",
				"cycle", nextCheckpoint,
				"register x", registerX,
				"strength", strength,
				"sum", strentghSum)
			nextCheckpoint += checkpointInterval
		}
		if cycle == nextActiveCycle && frameIndex < len(frames) {
			each := frames[frameIndex]
			frameIndex++

			if each.instr == addx {
				stack[cycle+1] = each.value
				nextActiveCycle = cycle + 2
			} else {
				nextActiveCycle = cycle + 1
			}
		}
		if value, ok := stack[cycle]; ok {
			registerX += value
			delete(stack, cycle)
		}
		if frameIndex >= len(frames) && len(stack) == 0 {
			break
		}
	}
	if len(stack) > 0 {
		shared.Logger.Error("Stack should be empty.", "stack", stack)
		panic("non-empty stack")
	}
	shared.Logger.Info("Signal strength sum counted.", "sum", strentghSum)
	return strentghSum, toLines(crt)
}

func toLines(content [][]rune) []string {
	lines := make([]string, len(content))
	for i, each := range content {
		lines[i] = string(each)
	}
	return lines
}

type instruction int

type frame struct {
	instr instruction
	value int
}

func (v frame) String() string {
	var valueStr string
	if v.instr != noop {
		valueStr = " " + strconv.Itoa(v.value)
	}
	return fmt.Sprintf("{%v%s}", v.instr, valueStr)
}

func (v instruction) String() string {
	switch v {
	case addx:
		return "addx"
	case noop:
		return "noop"
	default:
		return "unknown"
	}
}

func toInstruction(s string) instruction {
	switch s {
	case "addx":
		return addx
	case "noop":
		return noop
	default:
		panic("invalid instruction: \"" + s + "\"")
	}
}

func parseLines(lines []string) []frame {
	var frames []frame
	for _, each := range lines {
		pieces := strings.Fields(each)
		if len(pieces) < 1 {
			shared.Logger.Warn("Invalid line, skipping.", "line", each)
			continue
		}
		instr := toInstruction(pieces[0])
		var value int
		if instr != noop {
			var err error
			value, err = strconv.Atoi(pieces[1])
			if err != nil {
				shared.Logger.Warn(
					"Failed to parse int, skipping line.",
					"line", each,
					"err", err)
				continue
			}
		}
		frames = append(frames, frame{instr: instr, value: value})
	}
	return frames
}

func deriveCrtPixel(cycle, crtWidth, registerX int) rune {
	crtPixelIndex := (cycle - 1) % crtWidth
	lowCrt := registerX - 1
	highCrt := registerX + 1
	if lowCrt <= crtPixelIndex && crtPixelIndex <= highCrt {
		return '#'
	}
	return '.'
}
