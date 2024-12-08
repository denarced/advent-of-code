package aoc2404

import "github.com/denarced/advent-of-code/shared"

func CountInTable(table []string, word string) int {
	if len(word) == 0 {
		return 0
	}
	count := 0
	for r := 0; r < len(table); r++ {
		for c := 0; c < len(table[r]); c++ {
			if table[r][c] == word[0] {
				count += countWordsAt(table, word, r, c)
			}
		}
	}
	return count
}

func countWordsAt(table []string, word string, row, col int) int {
	directions := []shared.Direction{
		{X: 1, Y: 0},
		{X: 1, Y: -1},
		{X: 0, Y: -1},
		{X: -1, Y: -1},
		{X: -1, Y: 0},
		{X: -1, Y: 1},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}
	count := 0
	for _, each := range directions {
		if readTableAt(table, row, col, len(word), each) == word {
			count++
		}
	}
	return count
}

func CountWordCrosses(table []string, word string) int {
	locations := findWordLocations(table, word)
	counts := map[shared.Location]int{}
	total := 0
	for _, each := range locations {
		if count, ok := counts[each]; ok {
			if count > 1 {
				panic("wtf")
			}
			counts[each]++
			total++
			shared.Logger.Info("Found word cross.", "location", each)
		} else {
			shared.Logger.Debug("Found half of word cross.", "location", each)
			counts[each] = 1
		}
	}
	return total
}

func findWordLocations(table []string, word string) []shared.Location {
	if len(word)%2 != 1 {
		panic("only works with odd length words: 3, 5, 7, ...")
	}
	var locations []shared.Location
	mid := len(word) / 2
	for r := 0; r < len(table); r++ {
		for c := 0; c < len(table[r]); c++ {
			for _, each := range shared.Directions {
				if each.X == 0 || each.Y == 0 {
					continue
				}
				if readTableAt(table, r, c, len(word), each) == word {
					x := r + each.X*mid
					y := c + each.Y*mid
					locations = append(locations, shared.Location{X: x, Y: y})
				}
			}
		}
	}
	return locations
}

func readTableAt(table []string, row, col, count int, dir shared.Direction) string {
	result := ""
	for range count {
		if row < 0 || row >= len(table) {
			break
		}
		line := table[row]
		if col < 0 || col >= len(line) {
			break
		}
		result += line[col : col+1]
		row += dir.X
		col += dir.Y
	}
	return result
}
