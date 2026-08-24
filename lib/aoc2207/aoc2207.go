package aoc2207

import (
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	directory nodeType = iota
	file
)

func SumRecursiveDirSize(lines []string) (sum, sizeToDelete int) {
	shared.Logger.Info("Start calculating recursive dir size.", "line count", len(lines))
	fsys := parseLines(lines)
	requiredAvailableSpace := 30_000_000
	extraSpaceNeeded := requiredAvailableSpace - fsys.available
	shared.Logger.Info(
		"Filesystem ready.",
		"available space", fsys.available,
		"space needed", extraSpaceNeeded)
	dirs := fsys.findDirs(fsys.root)
	shared.Logger.Info("Dirs gathered.", "count", len(dirs))
	var totalSize int
	minimumSize := 2_000_000_000
	minimumName := ""
	for _, each := range dirs {
		size := fsys.countSize(each)
		if size >= extraSpaceNeeded {
			if size < minimumSize {
				minimumSize = size
				minimumName = each.name
				shared.Logger.Info("Dir candidate to be deleted.", "name", each.name, "size", size)
			}
		}
		if size <= 100_000 {
			totalSize += size
		}
	}
	shared.Logger.Info(
		"Recursive dir size calculated.",
		"size", totalSize,
		"dir to delete", minimumName,
		"dir size to delete", minimumSize)
	return totalSize, minimumSize
}

func parseLines(lines []string) *filesystem {
	fsys := newFilesystem()
	for i := 0; i < len(lines); {
		fsys.execute(&i, lines)
	}
	return fsys
}

func newFilesystem() *filesystem {
	return &filesystem{
		cwd: []string{},
		root: &node{
			name: "/",
			kind: directory,
			kids: map[string]*node{},
		},
		available: 70_000_000,
	}
}

type nodeType int

type node struct {
	name string
	kind nodeType
	size int
	kids map[string]*node
}

type filesystem struct {
	cwd       []string
	root      *node
	available int
}

func (v *filesystem) execute(lineIndex *int, lines []string) {
	if lines[*lineIndex][0] != '$' {
		shared.Logger.Error("Command expected.", "index", *lineIndex)
		panic("command expected")
	}
	parts := strings.Fields(lines[*lineIndex])[1:]
	switch parts[0] {
	case "cd":
		v.cd(parts[1])
		*lineIndex++
	case "ls":
		*lineIndex += v.ls(lines[*lineIndex+1:]) + 1
	default:
		panic("unknown command: " + lines[*lineIndex])
	}
}

func (v *filesystem) cd(dirn string) {
	if dirn[0] == '/' {
		v.cwd = []string{dirn}
	} else if dirn == ".." {
		v.cwd = v.cwd[:len(v.cwd)-1]
	} else {
		v.cwd = append(v.cwd, dirn)
	}
}

func (v *filesystem) ls(lines []string) int {
	var listing []string
	for _, each := range lines {
		if each[0] == '$' {
			break
		}
		listing = append(listing, each)
	}
	for _, each := range listing {
		parts := strings.Fields(each)
		if parts[0] == "dir" {
			v.addDir(parts[1])
			continue
		}
		v.addFile(each)
	}
	return len(listing)
}

func (v *filesystem) addDir(dname string) {
	nod := v.getCurrentNode()
	nod.kids[dname] = &node{
		name: dname,
		kind: directory,
		kids: map[string]*node{},
	}
}

func (v *filesystem) getCurrentNode() *node {
	cur := v.root
	for _, each := range v.cwd[1:] {
		cur = cur.kids[each]
	}
	return cur
}

func (v *filesystem) addFile(line string) {
	pieces := strings.Fields(line)
	if len(pieces) != 2 {
		panic("should have 2 pieces on file line")
	}
	size := gent.OrPanic2(strconv.Atoi(pieces[0]))("failed to parse file size")
	v.available -= size
	nod := v.getCurrentNode()
	name := pieces[1]
	nod.kids[name] = &node{
		name: name,
		kind: file,
		size: size,
	}
}

func (v *filesystem) findDirs(nod *node) []*node {
	var nodes []*node
	for _, each := range nod.kids {
		if each.kind == directory {
			nodes = append(nodes, each)
			kids := v.findDirs(each)
			nodes = append(nodes, kids...)
		}
	}
	return nodes
}

func (v *filesystem) countSize(nod *node) int {
	var size int
	for _, each := range nod.kids {
		if each.kind == file {
			size += each.size
		} else {
			size += v.countSize(each)
		}
	}
	return size
}
