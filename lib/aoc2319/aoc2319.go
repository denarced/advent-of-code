package aoc2319

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func SumRatings(lines []string) int {
	shared.Logger.Info("Sum ratings.", "line count", len(lines))
	nameToFlow, parts := parseLines(lines)
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Got situation.", "worksflows", nameToFlow, "parts", parts)
	}
	steps := make([]step, len(parts))
	for i, each := range parts {
		steps[i] = step{workflowName: "in", aPart: each}
	}
	var sum int
	var aStep step
	var count int
	for len(steps) > 0 {
		count++
		aStep, steps = steps[0], steps[1:]
		nextStep := processStep(aStep, nameToFlow[aStep.workflowName])
		switch nextStep.workflowName {
		case "R":
			shared.Logger.Debug("Rejected.", "part", nextStep.aPart)
			continue
		case "A":
			sub := nextStep.aPart.sum()
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Add to the sum.", "sub", sub)
			}
			sum += sub
		default:
			steps = append(steps, nextStep)
		}
	}
	shared.Logger.Info("Got ratings sum.", "sum", sum, "loop count", count)
	return sum
}

func parseLines(lines []string) (map[string]workflow, []part) {
	blocks := shared.SplitToBlocks(lines)
	if len(blocks) != 2 {
		shared.Logger.Error("Invalid block count.", "count", len(blocks), "blocks", blocks)
		panic("invalid data - bad block count")
	}
	workflows := parseWorkflows(blocks[0])
	parts := parseParts(blocks[1])
	return workflows, parts
}

type step struct {
	workflowName string
	aPart        part
}

type workflow struct {
	name  string
	specs []spec
}

type spec struct {
	attr        string
	dest        string
	value       int
	less        bool
	endComplete bool
}

type part map[string]int

func (v part) sum() int {
	var sum int
	for _, v := range v {
		sum += v
	}
	return sum
}

func (v part) pick(name string) int {
	val, ok := v[name]
	if !ok {
		panic("no such part attribute: " + name)
	}
	return val
}

func parseWorkflows(lines []string) map[string]workflow {
	m := map[string]workflow{}
	for _, each := range lines {
		flow := parseWorkflow(each)
		m[flow.name] = flow
	}
	return m
}

func parseWorkflow(s string) workflow {
	openIndex := strings.Index(s, "{")
	closeIndex := strings.Index(s, "}")
	name := s[:openIndex]
	flow := workflow{name: name}
	for _, each := range strings.Split(s[openIndex+1:closeIndex], ",") {
		flow.specs = append(flow.specs, parseWorkflowSpec(each))
	}
	return flow
}

func parseWorkflowSpec(s string) spec {
	pieces := strings.Split(s, ":")
	if len(pieces) == 1 {
		return spec{dest: s, endComplete: true}
	}
	opPieces := strings.Split(pieces[0], "<")
	less := true
	if len(opPieces) == 1 {
		opPieces = strings.Split(pieces[0], ">")
		less = false
	}
	val := parseIntOrDie(opPieces[1])
	return spec{
		attr:  opPieces[0],
		dest:  pieces[1],
		value: val,
		less:  less,
	}
}

func parseParts(lines []string) []part {
	var parts []part
	for _, each := range lines {
		parts = append(parts, parsePart(each))
	}
	return parts
}

func parsePart(s string) part {
	openIndex := strings.Index(s, "{")
	closeIndex := strings.Index(s, "}")
	content := s[openIndex+1 : closeIndex]
	aPart := map[string]int{}
	for _, each := range strings.Split(content, ",") {
		pieces := strings.Split(each, "=")
		num := pieces[1]
		aPart[pieces[0]] = parseIntOrDie(num)
	}

	return aPart
}

func parseIntOrDie(num string) int {
	val, err := strconv.Atoi(num)
	if err != nil {
		shared.Logger.Error("Failed to parse int.", "s", num, "err", err)
		panic(err)
	}
	return val
}

func processStep(aStep step, flow workflow) step {
	for _, each := range flow.specs {
		if each.endComplete {
			return step{workflowName: each.dest, aPart: aStep.aPart}
		}
		value := aStep.aPart.pick(each.attr)
		if each.less {
			if value < each.value {
				return step{workflowName: each.dest, aPart: aStep.aPart}
			}
		} else if value > each.value {
			return step{workflowName: each.dest, aPart: aStep.aPart}
		}
	}
	panic("my ass")
}

type restriction struct {
	low  int
	high int
}

func (v restriction) String() string {
	return fmt.Sprintf("(%d-%d)", v.low, v.high)
}

type policy map[string]restriction

func (v policy) intersection(other policy) (result policy) {
	result = policy{}
	for key, this := range v {
		if this.high < other[key].low || other[key].high < this.low {
			result = nil
			return
		}
		result[key] = restriction{
			low:  max(this.low, other[key].low),
			high: min(this.high, other[key].high),
		}
	}
	return
}

func (v policy) String() (res string) {
	for key, val := range v {
		s := fmt.Sprintf("%s%s", key, val.String())
		if res != "" {
			res += " "
		}
		res += s
	}
	res = "[" + res + "]"
	return
}

func (v policy) multiply() int {
	var prod int
	for _, each := range v {
		length := each.high - each.low + 1
		if prod == 0 {
			prod = length
		} else {
			prod *= length
		}
	}
	return prod
}

func Negotiate(lines []string, genDefaultPolicy func() policy) int {
	workflows, _ := parseLines(lines)
	var policies []policy
	pierce(
		workflows,
		"in",
		nil,
		func(link *shared.Link[comparison]) {
			pol := derivePolicy(link, genDefaultPolicy)
			if pol != nil {
				policies = append(policies, pol)
			}
		})
	shared.Logger.Info("Start to count result.", "policy count", len(policies))
	land := countLand(policies)
	shared.Logger.Info("Land derived.", "land", land)
	return land
}

func derivePolicy(link *shared.Link[comparison], genDefaultPolicy func() policy) policy {
	var pol policy
	if genDefaultPolicy == nil {
		pol = policy{
			"x": restriction{low: 1, high: 4000},
			"a": restriction{low: 1, high: 4000},
			"m": restriction{low: 1, high: 4000},
			"s": restriction{low: 1, high: 4000},
		}
	} else {
		pol = genDefaultPolicy()
	}
	for link != nil {
		rest, ok := pol[link.Item.attr]
		if !ok {
			panic("no such restriction: " + link.Item.attr)
		}
		rest, ok = cutRestriction(rest, link.Item)
		if !ok {
			return nil
		}
		pol[link.Item.attr] = rest
		link = link.Parent
	}
	return pol
}

func cutRestriction(rest restriction, comp comparison) (restriction, bool) {
	if comp.less {
		next := comp.val - 1
		if next < rest.high {
			rest.high = next
		}
	} else {
		next := comp.val + 1
		if next > rest.low {
			rest.low = next
		}
	}
	return rest, rest.low <= rest.high
}

type comparison struct {
	attr string
	val  int
	less bool
}

func (v comparison) String() string {
	op := "<"
	if !v.less {
		op = ">"
	}
	return fmt.Sprintf("%s%s%d", v.attr, op, v.val)
}

func (v comparison) reverse() comparison {
	var val int
	if v.less {
		val = v.val - 1
	} else {
		val = v.val + 1
	}
	return comparison{
		attr: v.attr,
		val:  val,
		less: !v.less,
	}
}

func pierce(
	workflows map[string]workflow,
	name string,
	comparisons *shared.Link[comparison],
	cb func(*shared.Link[comparison]),
) {
	flow := workflows[name]
	for _, each := range flow.specs {
		if each.endComplete {
			if each.dest == "A" {
				cb(comparisons)
			} else if each.dest != "R" {
				pierce(workflows, each.dest, comparisons, cb)
			}
			continue
		}
		c := comparison{attr: each.attr, val: each.value, less: each.less}
		if each.dest != "R" {
			next := shared.AddLink(comparisons, c)
			if each.dest == "A" {
				cb(next)
			} else {
				pierce(workflows, each.dest, next, cb)
			}
		}
		comparisons = shared.AddLink(comparisons, c.reverse())
	}
}

func countLand(policies []policy) int {
	keys := collectSortedPolicyKeys(policies)
	if len(keys) == 0 {
		return 0
	}
	sortPolicies(policies, keys)
	primary := keys[0]
	var overlap int
	for i := 0; i < len(policies)-1; i++ {
		first := policies[i]
		for j := i + 1; j < len(policies); j++ {
			second := policies[j]
			if first[primary].high < second[primary].low {
				break
			}
			if cross := first.intersection(second); cross != nil {
				overlap += cross.multiply()
			}
		}
	}
	var total int
	for _, each := range policies {
		total += each.multiply()
	}
	return total - overlap
}

func collectSortedPolicyKeys(policies []policy) []string {
	if len(policies) == 0 {
		return nil
	}
	m := map[string]int{}
	for _, each := range policies {
		for k := range each {
			m[k] = 0
		}
	}
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, len(m))
	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortPolicies(policies []policy, keys []string) {
	getLow := func(pol policy, key string) int {
		value, ok := pol[key]
		if ok {
			return value.low
		}
		panic("policy doesn't have dimension " + key)
	}
	slices.SortFunc(policies, func(a, b policy) int {
		for _, key := range keys {
			x := getLow(a, key)
			y := getLow(b, key)
			diff := x - y
			if diff != 0 {
				return diff
			}
		}
		return 0
	})
}
