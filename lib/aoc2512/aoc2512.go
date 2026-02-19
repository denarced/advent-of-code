package aoc2512

import (
	"fmt"
	"strconv"
	"strings"
)

func CountFittingRegions(lines []string) (int, error) {
	var fittingCount int
	presents, regions, err := parseLines(lines)
	if err != nil {
		return 0, err
	}
	for _, each := range regions {
		availableSpots := each.width * each.height
		for i, mul := range each.counts {
			if mul == 0 {
				continue
			}
			availableSpots -= mul * presents[i].spots
		}
		if availableSpots >= 0 {
			fittingCount++
		}
	}
	return fittingCount, nil
}

type present struct {
	table []string
	spots int
}
type region struct {
	width, height int
	counts        []int
}

func parseLines(lines []string) (presents []present, regions []region, err error) {
	var blocks [][]string
	var current []string
	for _, each := range lines {
		if current != nil && each == "" {
			blocks = append(blocks, current)
			current = nil
			continue
		}
		if each != "" {
			current = append(current, each)
			continue
		}
	}
	if current != nil {
		blocks = append(blocks, current)
	}
	for i := 0; i < len(blocks)-1; i++ {
		presents = append(presents, parsePresentLines(blocks[i]))
	}
	for _, each := range blocks[len(blocks)-1] {
		reg, regErr := parseRegionLine(each)
		if regErr != nil {
			err = regErr
			return
		}
		regions = append(regions, reg)
	}
	return
}

func parsePresentLines(lines []string) present {
	table := lines[1:]
	var spots int
	for _, each := range table {
		for _, c := range each {
			if c == '#' {
				spots++
			}
		}
	}
	return present{table: table, spots: spots}
}

func parseRegionLine(line string) (region, error) {
	pieces := strings.SplitN(line, ":", 2)
	dimPieces := strings.SplitN(strings.TrimSpace(pieces[0]), "x", 2)
	width, err := strconv.Atoi(dimPieces[0])
	if err != nil {
		return region{}, fmt.Errorf("failed to convert region width - %w", err)
	}
	height, err := strconv.Atoi(dimPieces[1])
	if err != nil {
		return region{}, fmt.Errorf("failed to convert region height - %s", err)
	}
	fields := strings.Fields(pieces[1])
	counts := make([]int, len(fields))
	for i, each := range fields {
		n, err := strconv.Atoi(each)
		if err != nil {
			return region{}, fmt.Errorf("failed to convert region index - %w", err)
		}
		counts[i] = n
	}
	return region{
		width:  width,
		height: height,
		counts: counts,
	}, nil
}
