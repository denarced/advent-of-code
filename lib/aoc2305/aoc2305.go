package aoc2305

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

type intRange struct {
	start, end int
}

func toRanges(ints []int, useRange bool) []intRange {
	var result []intRange
	if !useRange {
		for _, each := range ints {
			result = append(result, intRange{start: each, end: each})
		}
	} else {
		for i := 0; i < len(ints); i += 2 {
			result = append(result, intRange{start: ints[i], end: ints[i] + ints[i+1] - 1})
		}
	}
	return result
}

func DeriveLowestLocation(lines []string, useRange bool) int {
	shared.Logger.Info("Derive lowest location.", "range", useRange)
	seeds, corrs := parseLines(lines)
	lowest := math.MaxInt
	for _, aRange := range toRanges(seeds, useRange) {
		for i := aRange.start; i <= aRange.end; i++ {
			seed := i
			for _, aCorr := range corrs {
				seed = aCorr.translate(seed)
			}
			candidate := min(lowest, seed)
			if candidate < lowest {
				shared.Logger.Info("New lowest found.", "seed", seed)
			}
			lowest = candidate
		}
	}
	shared.Logger.Info("Lowest found.", "lowest", lowest)
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
	blocks := shared.SplitToBlocks(lines[1:])
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
	seeds := gent.Map(strings.Fields(pieces[1]), func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			shared.Logger.Error("Failed to convert seed to number", "err", err)
			panic(err)
		}
		return i
	})
	return seeds
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
