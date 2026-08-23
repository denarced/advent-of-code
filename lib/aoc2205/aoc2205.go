package aoc2205

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveTopCrates(lines []string) string {
	stor, instructions := parseLines(lines)
	shared.Logger.Info(
		"Parsed.",
		"store count", len(stor.stacks),
		"instruction count", len(instructions))
	for _, instr := range instructions {
		for range instr.count {
			pop(stor, instr)
		}
	}
	var cream []rune
	for i := 1; i <= len(stor.stacks); i++ {
		st := stor.stacks[i]
		cream = append(cream, st[0])
	}
	result := string(cream)
	shared.Logger.Info("Top cretes derived.", "top", result)
	return result
}

func parseLines(lines []string) (storage, []instruction) {
	storageLines, instructionLines := partitionLines(lines)
	shared.Logger.Info(
		"Lines parsed.",
		"storage", len(storageLines),
		"instructions", len(instructionLines))
	return parseStorageLines(storageLines), toInstructions(instructionLines)
}

type storage struct {
	stacks map[int][]rune
}

type instruction struct {
	count int
	from  int
	to    int
}

func partitionLines(lines []string) (storLines, instrLines []string) {
	var parsingInstructions bool
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			parsingInstructions = true
			continue
		}
		if parsingInstructions {
			instrLines = append(instrLines, each)
		} else {
			storLines = append(storLines, each)
		}
	}
	return
}

func toInstructions(lines []string) []instruction {
	var instructions []instruction
	for _, each := range lines {
		var count, from, to int
		parsedCount, err := fmt.Sscanf(each, "move %d from %d to %d", &count, &from, &to)
		if err != nil {
			shared.Logger.Error("Failed to parse instruction line.", "err", err, "line", each)
			panic(err)
		}
		if parsedCount != 3 {
			shared.Logger.Error("Unexpected parsed count.", "expected", parsedCount)
			panic("parsed count should be 3")
		}
		instructions = append(instructions, instruction{
			count: count,
			from:  from,
			to:    to,
		})
	}
	return instructions
}

func parseStorageLines(lines []string) storage {
	lastLine := lines[len(lines)-1]
	indexToID := map[int]int{}
	for i, c := range lastLine {
		if '0' <= c && c <= '9' {
			indexToID[i] = int(c - '0')
		}
	}

	stor := storage{stacks: map[int][]rune{}}
	for i := 0; i < len(lines)-1; i++ {
		current := []rune(lines[i])
		for colIndex, ID := range indexToID {
			if colIndex < len(current) {
				c := current[colIndex]
				if 'A' <= c && c <= 'Z' {
					stor.stacks[ID] = append(stor.stacks[ID], c)
				}
			}
		}
	}
	return stor
}

func pop(stor storage, instr instruction) {
	from := stor.stacks[instr.from]
	to := stor.stacks[instr.to]
	to = append([]rune{from[0]}, to...)
	from = from[1:]
	stor.stacks[instr.from] = from
	stor.stacks[instr.to] = to
}
