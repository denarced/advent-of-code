package aoc2506

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

type operator int

const (
	opAdd operator = iota
	opMul
)

type expr struct {
	values []int64
	op     operator
}

func Calculate(lines []string, standardMath bool) int64 {
	var expressions []expr
	var err error
	if standardMath {
		expressions, err = parseLines(lines)
	} else {
		expressions, err = parseLinesInColumns(lines)
	}
	if err != nil {
		shared.Logger.Error(
			"Failed to parse lines.",
			"err", err,
			"line count", len(lines),
			"stanard math", standardMath)
		panic("failed to parse lines")
	}
	var grand int64
	for _, each := range expressions {
		var total int64
		if each.op == opMul {
			total = 1
		}
		for _, factor := range each.values {
			switch each.op {
			case opAdd:
				total += factor
			case opMul:
				total *= factor
			default:
				panic(fmt.Sprintf("invalid operator: %d", each.op))
			}
		}
		grand += total
	}
	return grand
}

func parseLines(lines []string) ([]expr, error) {
	var expressions []expr
	var parsed [][]string
	for _, line := range lines {
		parsed = append(parsed, strings.Fields(line))
	}
	count := len(parsed[0])
	// Verify that there's equal number of fields.
	for i := 1; i < len(parsed); i++ {
		if len(parsed[i]) != count {
			return nil, fmt.Errorf(
				"line has %d fields, should be %d, line: %v",
				len(parsed[i]),
				count,
				parsed[i],
			)
		}
	}

	// Build columns.
	var columns [][]string
	for i := range count {
		var column []string
		for _, each := range parsed {
			column = append(column, each[i])
		}
		columns = append(columns, column)
	}

	// Build expressions.
	for i := range count {
		column := columns[i]
		expression := expr{op: parseOp(column[len(column)-1])}
		for j := 0; j < len(column)-1; j++ {
			value, err := strconv.ParseInt(column[j], 10, 64)
			if err != nil {
				return nil, err
			}
			expression.values = append(expression.values, value)
		}
		expressions = append(expressions, expression)
	}
	return expressions, nil
}

func parseOp(s string) operator {
	switch s {
	case "+":
		return opAdd
	case "*":
		return opMul
	default:
		panic(fmt.Sprintf("invalid operator: %s", s))
	}
}

func parseLinesInColumns(lines []string) ([]expr, error) {
	if err := verifyRectangle(lines); err != nil {
		return nil, err
	}
	table := splitByEmptyColumns(lines)
	var expressions []expr
	for _, each := range table {
		op := parseOp(strings.TrimSpace(each[len(each)-1]))
		var values []int64
		for _, s := range pivot(each[0 : len(each)-1]) {
			value, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
		}
		expression := expr{values: values, op: op}
		expressions = append(expressions, expression)
	}
	return expressions, nil
}

func verifyRectangle(lines []string) error {
	count := -1
	for _, each := range lines {
		if count < 0 {
			count = len(each)
			continue
		}
		if len(each) != count {
			return fmt.Errorf(
				"lines are not of equal length, length: %d, expected: %d",
				len(each),
				count)
		}
	}
	return nil
}

func splitByEmptyColumns(lines []string) [][]string {
	if len(lines) == 0 {
		return nil
	}
	length := len(lines[0])
	var empties []int
mainLoop:
	for i := range length {
		for _, each := range lines {
			if each[i] != ' ' {
				continue mainLoop
			}
		}
		empties = append(empties, i)
	}
	table := make([][]string, len(empties)+1)
	for _, each := range lines {
		for i, l := range splitWithIndexes(each, empties) {
			table[i] = append(table[i], l)
		}
	}
	return table
}

func splitWithIndexes(s string, indexes []int) []string {
	var current []rune
	var runes [][]rune
	for i, r := range s {
		if slices.Contains(indexes, i) {
			if len(current) > 0 {
				runes = append(runes, current)
				current = nil
			}
			continue
		}
		current = append(current, r)
	}
	if len(current) > 0 {
		runes = append(runes, current)
	}
	var result []string
	for _, each := range runes {
		result = append(result, string(each))
	}
	return result
}

func pivot(table []string) []string {
	if len(table) == 0 {
		return table
	}
	l := len(table[0])
	var result []string
	for i := l - 1; i >= 0; i-- {
		var line []rune
		for _, each := range table {
			line = append(line, rune(each[i]))
		}
		result = append(result, string(line))
	}
	return result
}
