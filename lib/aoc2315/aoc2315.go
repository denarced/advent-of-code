package aoc2315

import (
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

const boxCount = 256

const (
	opRemove opType = iota
	opAdd
)

type opType int

func SumHashes(lines []string) int {
	shared.Logger.Info("Sum hashes.")
	var sum int
	parseLines(lines, func(each string) {
		sum += hash(each)
	})
	shared.Logger.Info("Hashes summed.", "sum", sum)
	return sum
}

func hash(s string) (h int) {
	for _, r := range s {
		h = ((h + int(r)) * 17) % boxCount
	}
	return
}

func parseLines(lines []string, cb func(string)) {
	for _, line := range lines {
		for _, each := range strings.Split(line, ",") {
			if each != "" {
				cb(each)
			}
		}
	}
}

type lens struct {
	label string
	focal int
}

type lensBox struct {
	lenses []lens
}

func (v *lensBox) find(label string) (int, bool) {
	for i, each := range v.lenses {
		if each.label == label {
			return i, true
		}
	}
	return 0, false
}

func (v *lensBox) add(label string, focal int) {
	index, exists := v.find(label)
	if exists {
		v.lenses[index].focal = focal
		return
	}
	v.lenses = append(v.lenses, lens{
		label: label,
		focal: focal,
	})
}

func (v *lensBox) remove(label string) {
	index, exists := v.find(label)
	if !exists {
		return
	}
	shortened := make([]lens, 0, len(v.lenses)-1)
	if index > 0 {
		shortened = v.lenses[:index]
	}
	if len(v.lenses) > index+1 {
		shortened = append(shortened, v.lenses[index+1:]...)
	}
	if len(shortened) == 0 {
		shortened = nil
	}
	v.lenses = shortened
}

func getBox(boxes []*lensBox, index int) *lensBox {
	box := boxes[index]
	if box == nil {
		box = new(lensBox)
		boxes[index] = box
	}
	return box
}

func DeriveFocusingPower(lines []string) int {
	shared.Logger.Info("Derive focusing power.")
	boxes := make([]*lensBox, boxCount)
	parseLines(lines, func(each string) {
		shared.Logger.Debug("Process command.", "each", each)
		label, kind, focal := splitCommand(each)
		boxIndex := hash(label)
		box := getBox(boxes, boxIndex)
		switch kind {
		case opAdd:
			box.add(label, focal)
		case opRemove:
			box.remove(label)
		default:
			panic("unknown opType")
		}
	})

	shared.Logger.Info("Sum lenses to derive focusing power.")
	var power int
	for i, box := range boxes {
		if box == nil {
			continue
		}
		for j, each := range box.lenses {
			pwr := (i + 1) * (j + 1) * each.focal
			power += pwr
		}
	}
	shared.Logger.Info("Focusing power derived.", "power", power)
	return power
}

func splitCommand(cmd string) (label string, kind opType, focal int) {
	pieces := strings.Split(cmd, "=")
	if len(pieces) == 2 {
		kind = opAdd
		var err error
		focal, err = strconv.Atoi(pieces[1])
		if err != nil {
			shared.Logger.Error("Failed to convert focal length.",
				"from", pieces[1],
				"cmd", cmd,
				"err", err)
			panic(err)
		}
	} else {
		kind = opRemove
		pieces = strings.Split(cmd, "-")
		if len(pieces) != 2 {
			shared.Logger.Error("Failed to split command.", "cmd", cmd)
			panic("failed to split command")
		}
	}
	label = pieces[0]
	return
}
