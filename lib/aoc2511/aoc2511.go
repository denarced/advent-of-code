package aoc2511

import "strings"

func CountPaths(lines []string) int {
	shortcuts := parseLines(lines)
	stack := []node{shortcuts["you"]}
	var count int
	for len(stack) > 0 {
		var nextStack []node
		for _, each := range stack {
			for name, nod := range each {
				if name == "out" {
					count++
				} else {
					nextStack = append(nextStack, nod)
				}
			}
		}
		stack = nextStack
	}
	return count
}

type node map[string]node

func parseLines(lines []string) node {
	shortcuts := map[string]node{}
	for _, each := range lines {
		from, to := parseLine(each)
		var nod node
		if existing, ok := shortcuts[from]; ok {
			nod = existing
		} else {
			nod = map[string]node{}
			shortcuts[from] = nod
		}
		for _, one := range to {
			if existing, ok := shortcuts[one]; ok {
				nod[one] = existing
			} else {
				nose := map[string]node{}
				shortcuts[one] = nose
				nod[one] = shortcuts[one]
			}
		}
	}
	return shortcuts
}

func parseLine(line string) (string, []string) {
	pieces := strings.Split(strings.TrimSpace(line), ":")
	return strings.TrimSpace(pieces[0]), strings.Fields(pieces[1])
}
