package aoc2409

import (
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

type deque struct {
	queue []int
	left  int
	right int
}

func newDeque(s string) *deque {
	ints, err := shared.ToInts(strings.Split(s, ""))
	shared.Die(err, "newDeque -> ToInts")
	right := len(ints) - 1
	if right%2 == 1 {
		right--
	}
	return &deque{
		queue: ints,
		left:  0,
		right: right,
	}
}

func (v *deque) hasMore() bool {
	return v.left < len(v.queue)
}

func (v *deque) decStart() (id int, file bool, ok bool) {
	var width int
	for i := v.left; i < len(v.queue); i++ {
		width = v.queue[i]
		if width == 0 {
			continue
		}
		ok = true
		v.left = i
		break
	}
	id = v.left / 2
	file = v.left%2 == 0
	v.queue[v.left] = width - 1
	shared.Logger.Debug(
		"Dec start.",
		"id",
		id,
		"file",
		file,
		"ok",
		ok,
		"new width",
		v.queue[v.left],
	)
	return
}

func (v *deque) decEnd() (id int, ok bool) {
	var width int
	for i := v.right; i > v.left; i -= 2 {
		width = v.queue[i]
		if width == 0 {
			continue
		}
		ok = true
		v.right = i
		break
	}
	if !ok {
		return
	}
	v.queue[v.right] = width - 1
	id = v.right / 2
	shared.Logger.Debug("Dec end.", "id", id, "ok", ok, "new width", v.queue[v.right])
	return
}

func CountChecksum(s string) int {
	shared.Logger.Info("Count checksum.", "disk length", len(s))
	deq := newDeque(s)
	shared.Logger.Debug("Deque created.", "deque", deq)
	checksum := 0
	pos := 0
	summed := []int{}
	for deq.hasMore() {
		id, file, ok := deq.decStart()
		if !ok {
			shared.Logger.Debug("Break: start != ok.")
			break
		}
		if file {
			checksum += pos * id
			summed = append(summed, id)
			shared.Logger.Debug("Checksum updated (left).", "pos", pos, "id", id, "new", checksum)
		} else {
			if !deq.hasMore() {
				shared.Logger.Debug("Break: no more.")
				break
			}
			rightID, ok := deq.decEnd()
			if ok {
				checksum += pos * rightID
				summed = append(summed, rightID)
				shared.Logger.Debug(
					"Checksum updated (right).",
					"pos",
					pos,
					"id",
					rightID,
					"new",
					checksum)
			}
		}
		pos++
	}
	shared.Logger.Debug("Summed.", "values", summed)
	return checksum
}

func CountDefragmentedChecksum(s string) int {
	shared.Logger.Info("Count defragmented checksum.", "disk length", len(s))
	if len(s) < 3 {
		return 0
	}
	shared.Logger.Info("Defrag.")
	org := defrag(toAtoms(s))

	pos := 0
	checksum := 0
	summed := []int{}
	shared.Logger.Debug("Count checksum.", "org", org)
	for _, each := range org {
		for range each.width {
			if each.file && each.width > 0 {
				checksum += pos * each.id
				summed = append(summed, each.id)
			} else {
				summed = append(summed, 0)
			}
			pos++
		}
	}
	shared.Logger.Debug("Summed.", "summed", summed)
	return checksum
}

type atom struct {
	width int
	id    int
	file  bool
}

func slipIn(org []atom, piece atom, i int) []atom {
	evolved := append([]atom{}, org[0:i]...)
	evolved = append(evolved, piece)
	evolved = append(evolved, org[i:]...)
	return evolved
}

func visualize(org []atom) string {
	v := ""
	for _, each := range org {
		if each.file {
			v += strings.Repeat(strconv.FormatInt(int64(each.id), 10), each.width)
		} else {
			v += strings.Repeat(".", each.width)
		}
	}
	return v
}

func toAtoms(s string) []atom {
	ints := shared.OrPanic2(shared.ToInts(strings.Split(s, "")))("to ints")
	org := []atom{}
	for i, each := range ints {
		id := 0
		file := i%2 == 0
		if file {
			id = i / 2
		}
		org = append(org, atom{width: each, id: id, file: file})
	}
	return org
}

func defrag(org []atom) []atom {
	j := len(org) - 1
	if !org[j].file {
		j--
	}
	for ; j > 1; j-- {
		toMove := org[j]
		if !toMove.file {
			shared.Logger.Debug("Skip space.")
			continue
		}
		if toMove.width <= 0 {
			shared.Logger.Debug("Skip empty file.")
			continue
		}
		for i := 1; i < j; i++ {
			slot := org[i]
			if slot.file {
				continue
			}
			leftOverCount := slot.width - toMove.width
			if leftOverCount < 0 {
				continue
			}

			shared.Logger.Debug("Moved file found.", "i", i, "moved", toMove, "slot", slot)
			org[i] = toMove
			org[j] = atom{id: 0, width: toMove.width, file: false}
			if leftOverCount > 0 {
				leftOver := atom{width: leftOverCount, id: 0, file: false}
				org = slipIn(org, leftOver, i+1)
				j++
			}
			break
		}
	}
	return org
}
