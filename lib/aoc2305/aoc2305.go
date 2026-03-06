package aoc2305

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func DeriveLowestLocation(lines []string) int {
	seeds, corrs := parseLines(lines)
	lowest := math.MaxInt
	for _, each := range seeds {
		for _, fire := range corrs {
			each = fire.translate(each)
		}
		lowest = min(lowest, each)
	}
	return lowest
}

type spec struct {
	src  int
	dst  int
	size int
}

type corr struct {
	specs []spec
}

func (v corr) translate(value int) int {
	for _, each := range v.specs {
		if each.src <= value && value < each.src+each.size {
			return each.dst + (value - each.src)
		}
	}
	return value
}

func trimTitle(s string) string {
	return strings.Fields(s)[0]
}

func parseLines(lines []string) ([]int, []corr) {
	seeds := parseSeeds(lines[0])
	blocks := splitToBlocks(lines[1:])
	packs := make([]corr, 7)
	for _, each := range blocks {
		switch trimTitle(each[0]) {
		case "seed-to-soil":
			packs[0] = parseMap(each)
		case "soil-to-fertilizer":
			packs[1] = parseMap(each)
		case "fertilizer-to-water":
			packs[2] = parseMap(each)
		case "water-to-light":
			packs[3] = parseMap(each)
		case "light-to-temperature":
			packs[4] = parseMap(each)
		case "temperature-to-humidity":
			packs[5] = parseMap(each)
		case "humidity-to-location":
			packs[6] = parseMap(each)
		default:
			panic(fmt.Sprintf("no such title: %s", each[0]))
		}
	}
	return seeds, packs
}

func parseSeeds(s string) []int {
	pieces := strings.SplitN(s, ":", 2)
	if strings.TrimSpace(pieces[0]) != "seeds" {
		panic("first line should start with \"seeds\"")
	}
	return gent.Map(strings.Fields(pieces[1]), func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			shared.Logger.Error("Failed to convert seed to number", "err", err)
			panic(err)
		}
		return i
	})
}

func splitToBlocks(lines []string) [][]string {
	var start int
	// Skip empty lines.
	for i := range lines {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" {
			start = i
			break
		}
	}
	var blocks [][]string
	var current []string
	for i := start; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" {
			current = append(current, trimmed)
			continue
		}
		if current != nil {
			blocks = append(blocks, current)
			current = nil
		}
	}
	if current != nil {
		blocks = append(blocks, current)
	}
	return blocks
}

func parseMap(lines []string) corr {
	specs := []spec{}
	for _, each := range lines[1:] {
		fields := strings.Fields(each)
		if len(fields) != 3 {
			shared.Logger.Error(
				"Invalid range pack.",
				"length",
				len(fields),
				"line",
				each,
				"lines",
				lines,
			)
			panic(
				fmt.Sprintf(
					"invalid range pack, invalid length (%d): %s, %+v",
					len(fields),
					each,
					lines,
				),
			)
		}
		values := gent.Map(fields, func(s string) int {
			i, err := strconv.Atoi(s)
			if err != nil {
				shared.Logger.Error("Failed to convert range number.", "s", s, "err", err)
				panic(err)
			}
			return i
		})
		specs = append(specs, spec{
			dst:  values[0],
			src:  values[1],
			size: values[2],
		})
	}
	return corr{specs: specs}
}
