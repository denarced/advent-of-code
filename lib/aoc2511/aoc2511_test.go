package aoc2511

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountPaths(t *testing.T) {
	readTestData := func(req *require.Assertions, filen string) []string {
		lines, err := inr.ReadPath("testdata/" + filen)
		req.NoError(err, "read test data")
		return lines
	}

	t.Run("you", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		lines := readTestData(req, "in.txt")
		count := CountPaths(lines, "you")
		req.Equal(5, count)
	})

	t.Run("svr", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		lines := readTestData(req, "in2.txt")
		count := CountPaths(lines, "svr")
		req.Equal(2, count)
	})
}

func TestSpider(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in2.txt")
	req.NoError(err, "failed to read test data")
	req.NotEmpty(lines, "lines are empty")
	points := parseTree(lines)
	webIt(points)
	for _, each := range []struct {
		name    string
		asshole bool
		devil   bool
	}{
		{"aaa", false, true},
		{"bbb", false, true},
		{"ccc", true, true},
		{"dac", true, false},
		{"ddd", true, false},
		{"eee", true, true},
		{"fff", true, false},
		{"fft", false, true},
		{"ggg", true, false},
		{"hhh", true, false},
		{"hub", true, false},
		{"out", true, false},
		{"tty", false, true},
	} {
		bobo := points[each.name]
		req.Equalf(each.asshole, bobo.foundByAsshole, "%s asshole", each.name)
		req.Equalf(each.devil, bobo.foundByDevil, "%s devil", each.name)
	}
}

func TestSvr(t *testing.T) {
	t.Run("mess", func(t *testing.T) {
		shared.InitTestLogging(t)
		lines := []string{
			"aaa: svr eee jjj fft",
			"eee: out",
			"svr: fft bbb ddd hhh kkk",
			"kkk: ddd",
			"hhh: fft",
			"ddd: out",
			"bbb: out",
			"fft: dac ccc fff",
			"fff: dac lll",
			"lll: out",
			"ccc: out",
			"dac: out ggg",
			"ggg: iii",
			"iii: out",
			"out: mmm nnn",
			"mmm: ooo ppp",
			"nnn: qqq rrr",
		}
		count := CountPaths(lines, "svr")
		req := require.New(t)
		// svr -> fft: 2 ways
		//     svr -> fft
		//     svr -> hhh -> fft
		// fft -> dac: 2 ways
		//     fft -> dac
		//     fft -> fff -> dac
		// dac -> out: 2 ways
		//     dac -> out
		//     dac -> ggg -> iii -> out
		// Total: 2 * 2 * 2
		req.Equal(8, count)
	})

	t.Run("single", func(t *testing.T) {
		lines := []string{
			"aaa: svr",
			"svr: fft",
			"fft: dac",
			"dac: out",
			"out: zzz",
		}
		count := CountPaths(lines, "svr")
		req := require.New(t)
		req.Equal(1, count)
	})

	t.Run("single with dead ends", func(t *testing.T) {
		lines := []string{
			"aaa: svr",
			"svr: fft bbb bbc",
			"fft: dac ccc ccd",
			"dac: out ddd dde",
			"out: zzz",
		}
		count := CountPaths(lines, "svr")
		req := require.New(t)
		req.Equal(1, count)
	})

	t.Run("double", func(t *testing.T) {
		lines := []string{
			"svr: aaa bbb",
			"aaa: fft",
			"bbb: fft",

			"fft: ccc ddd",
			"ccc: dac",
			"ddd: dac",

			"dac: eee fff",
			"eee: out",
			"fff: out",
		}
		count := CountPaths(lines, "svr")
		req := require.New(t)
		req.Equal(8, count)
	})

	t.Run("100", func(t *testing.T) {
		lines := []string{"aaa: svr"}
		// Counted connections.
		lines = generateLines(lines, "svr", "fft", "bbb", 100)
		lines = generateLines(lines, "fft", "dac", "ccc", 100)
		lines = generateLines(lines, "dac", "out", "ddd", 100)
		lines = append(lines, "out: zzz")

		// Jump over all key points.
		lines = generateLines(lines, "aaa", "bbb", "eee", 10)
		lines = generateLines(lines, "bbb", "ccc", "efe", 10)
		lines = generateLines(lines, "ccc", "ddd", "ege", 10)
		lines = generateLines(lines, "ddd", "zzz", "ehe", 10)

		// Jump from key point to middle points.
		lines = generateLines(lines, "svr", "ccc", "fff", 10)
		lines = generateLines(lines, "fft", "ddd", "fgf", 10)
		lines = generateLines(lines, "dac", "zzz", "fhf", 10)

		// From key point to key point, jumping over key point in between.
		lines = generateLines(lines, "svr", "dac", "ggg", 10)
		lines = generateLines(lines, "fft", "out", "ghg", 10)

		// Jump over everything, from key point to key point.
		lines = generateLines(lines, "aaa", "out", "hhh", 10)
		lines = generateLines(lines, "svr", "out", "hih", 10)

		count := CountPaths(lines, "svr")
		req := require.New(t)
		req.Equal(100*100*100, count)
	})

	t.Run("zig-zag", func(t *testing.T) {
		req := require.New(t)
		lines := []string{
			"aaa: bbb bbc",
			"bbb: svr",
			"bbc: svr",
			"svr: ccc ccd",
			"ccc: fft",
			"ccd: fft",
			"fft: ddd dde",
			"ddd: dac",
			"dde: dac",
			"dac: eee eef",
			"eee: out",
			"eef: out",
		}
		count := CountPaths(lines, "svr")
		req.Equal(8, count)
	})
}

func TestGenerateLines(t *testing.T) {
	var tests = []struct {
		existing        []string
		count           int
		expectedLengths []int
	}{
		{nil, 1, []int{2, 2}},
		{nil, 2, []int{3, 2, 2}},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v+%d", tt.existing, tt.count)
		t.Run(name, func(t *testing.T) {
			ass := assert.New(t)
			lines := generateLines(tt.existing, "aaa", "ccc", "bbb", tt.count)
			for i := range tt.expectedLengths {
				ass.Equalf(
					tt.expectedLengths[i],
					len(strings.Fields(lines[i])),
					"length mismatch: %d, actual: %v",
					i,
					lines[i],
				)
			}
		})
	}

	t.Run("zzz+1", func(t *testing.T) {
		lines := generateLines(nil, "dac", "out", "zzz", 2)
		req := require.New(t)
		req.Equal(
			[]string{
				"dac: zzz zza",
				"zzz: out",
				"zza: out",
			},
			lines)
	})
}

func generateLines(lines []string, from, to, mid string, count int) []string {
	if count < 1 {
		return lines
	}
	current := mid
	inc := func(index int) {
		c := current[index] + 1
		if c > 'z' {
			c = 'a'
		}
		altered := []rune(current)
		altered[index] = rune(c)
		current = string(altered)
	}
	ids := make([]string, 0, count)
	for range count {
		ids = append(ids, current)
		for i := len(current) - 1; i > 0; i-- {
			inc(i)
			if current[i] != mid[i] {
				break
			}
			if i == 1 {
				panic("out of bounds")
			}
		}
	}
	lines = append(lines, fmt.Sprintf("%s: %s", from, strings.Join(ids, " ")))
	for _, each := range ids {
		lines = append(lines, fmt.Sprintf("%s: %s", each, to))
	}
	return lines
}

func ss(s ...string) []string {
	return s
}

func TestCombDefect(t *testing.T) {
	run := func(name string, lines []string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			gen := func() (map[string]*Point, *Point, *Point) {
				fats := append([]string{"aaa: svr"}, lines...)
				fats = append(fats, "fft: dac", "dac: out", "out: zzz")
				parsed := parseTree(fats)
				return parsed, parsed["svr"], parsed["fft"]
			}
			var counts []int

			_, svr, fft := gen()
			weblessCounter := newCounter()
			// EXERCISE
			comb(fft, svr, false, nil, weblessCounter)
			counts = append(counts, weblessCounter.count)

			_, svr, fft = gen()
			weblessReverseCounter := newCounter()
			// EXERCISE
			comb(svr, fft, true, nil, weblessReverseCounter)
			counts = append(counts, weblessReverseCounter.count)

			logAllowed := func(allowed *gent.Set[*Point]) {
				if shared.IsDebugEnabled() {
					expr := gent.Map(allowed.ToSlice(), func(p *Point) string {
						var parents, kids string
						for each := range p.parents {
							if parents != "" {
								parents += " "
							}
							parents += each
						}
						for each := range p.kids {
							if kids != "" {
								kids += " "
							}
							kids += each
						}
						return "<" + parents + " + " + kids + ">"
					})
					slices.Sort(expr)
					shared.Logger.Debug(
						"About to comb.",
						"case",
						"webful",
						"allowed",
						expr,
					)
				}
			}

			parsed, svr, fft := gen()
			web(svr, fft)
			webfulCounter := newCounter()
			allowed := FilterDoubles(parsed)
			logAllowed(allowed)
			// EXERCISE
			comb(fft, svr, false, allowed, webfulCounter)
			counts = append(counts, webfulCounter.count)

			parsed, svr, fft = gen()
			web(svr, fft)
			webfulReverseCounter := newCounter()
			allowed = FilterDoubles(parsed)
			logAllowed(allowed)
			// EXERCISE
			comb(svr, fft, true, allowed, webfulReverseCounter)
			counts = append(counts, webfulReverseCounter.count)

			// VERIFY
			firstCount := counts[0]
			for i := 1; i < len(counts); i++ {
				req.Equalf(
					firstCount,
					counts[i],
					"count mismatch: %s, %s",
					counts,
					lines,
				)
			}
		})
	}

	run(
		"defect",
		// svr - bbb - ddd - eee - fft
		// svr - bbb - fft
		// svr - bbb - eee - fft
		// svr - bbb - ccc - ddd - eee - fft
		// svr - ddd - eee - fft
		ss(
			"svr: bbb ddd",
			"bbb: ddd fft eee ccc",
			"ccc: ddd",
			"ddd: eee",
			"eee: fft",
		))

	// This was an aid to find a defect in web function.
	gen := &generator{}
	for i := range 500 {
		lines := gen.generate()
		run(fmt.Sprint(i), lines)
	}
}

type generator struct {
	r *rand.Rand
}

func (v *generator) generate() []string {
	if v.r == nil {
		v.r = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	}
	nodeCount := 8
	nodes := generateNodes(nodeCount)
	frontLinks := map[int]*gent.Set[int]{}
	backLinks := map[int]*gent.Set[int]{}
	for i := range nodes {
		frontMax := len(nodes) - i - 1
		if frontMax > 0 {
			kids := v.addFrontLinks(i, frontMax, nodes, frontLinks)
			addBackLinks(i, kids, backLinks)
		}
		// First (svr) can't have backlinks.
		if i == 0 {
			continue
		}
		backs := backLinks[i]
		if backs != nil {
			continue
		}
		backIndex := v.r.IntN(i)
		backLinks[i] = gent.NewSet(backIndex)
		frontLinks[backIndex].Add(i)
	}
	return generateLinesFromLinks(nodes, frontLinks)
}

func generateIntRange(start, end int) []int {
	if end <= start {
		return nil
	}
	res := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		res = append(res, i)
	}
	return res
}

func (v *generator) addFrontLinks(
	i, frontMax int,
	nodes []string,
	frontLinks map[int]*gent.Set[int],
) []int {
	count := v.r.IntN(min(frontMax, 3)) + 1
	after := generateIntRange(i+1, len(nodes))
	kids := make([]int, 0, count)
	for len(kids) < count {
		index := v.r.IntN(len(after))
		kids = append(kids, after[index])
		after = append(after[:index], after[index+1:]...)
	}
	fronts := frontLinks[i]
	if fronts == nil {
		fronts = gent.NewSet(kids...)
		frontLinks[i] = fronts
	} else {
		for _, each := range kids {
			fronts.Add(each)
		}
	}
	return kids
}

func addBackLinks(i int, kids []int, backLinks map[int]*gent.Set[int]) {
	for _, each := range kids {
		backs := backLinks[each]
		if backs == nil {
			backs = gent.NewSet(i)
			backLinks[each] = backs
		} else {
			backs.Add(i)
		}
	}
}

func generateNodes(nodeCount int) []string {
	nodes := make([]string, nodeCount)
	for i := range nodeCount {
		c := rune('b' + i)
		name := string([]rune{c, c, c})
		nodes[i] = name
	}
	nodes = append([]string{"svr"}, nodes...)
	nodes = append(nodes, "fft")
	return nodes
}

func generateLinesFromLinks(nodes []string, frontLinks map[int]*gent.Set[int]) []string {
	var lines []string
	for i, name := range nodes {
		if i == len(nodes)-1 {
			break
		}
		kids := ""
		frontLinks[i].ForEachAll(func(index int) {
			if kids != "" {
				kids += " "
			}
			kids += nodes[index]
		})
		lines = append(lines, fmt.Sprintf("%s: %s", name, kids))
	}
	return lines
}
