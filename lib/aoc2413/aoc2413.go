package aoc2413

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

const conversionFix = 10_000_000_000_000

func DeriveFewestTokens(lines []string, fixConversion bool) int {
	shared.Logger.Debug("Parse machines.", "conversion fix", fixConversion)
	machines := parseMachines(lines, fixConversion)
	shared.Logger.Info("Machines parsed.", "count", len(machines))
	tokens := 0
	for _, each := range machines {
		cheapest := deriveCheapest(each.a, each.b, each.prize)
		logger := shared.Logger.With("machine", each)
		if cheapest < 0 {
			logger.Info("No possible solutions.")
			continue
		}
		logger.Info("Found cheapest.", "cheapest", cheapest)
		tokens += cheapest
	}
	return tokens
}

type button struct {
	name string
	loc  shared.Loc
}

type machine struct {
	a     button
	b     button
	prize shared.Loc
}

func parseMachines(lines []string, fixConversion bool) []machine {
	block := []string{}
	machines := []machine{}
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		block = append(block, trimmed)
		if len(block) >= 3 {
			machines = append(machines, parseMachine(block, fixConversion))
			block = []string{}
		}
	}
	return machines
}

func parseMachine(lines []string, fixConversion bool) machine {
	mac := machine{}
	for _, each := range lines {
		if mac.a.name == "" && strings.HasPrefix(each, "Button A") {
			mac.a = parseButton(each)
			continue
		}
		if mac.b.name == "" && strings.HasPrefix(each, "Button B") {
			mac.b = parseButton(each)
			continue
		}
		if strings.HasPrefix(each, "Prize:") {
			mac.prize = parsePrize(each, fixConversion)
		}
	}
	return mac
}

func parseButton(s string) button {
	prepped := strings.ReplaceAll(strings.ReplaceAll(s, ":", ""), ",", "")
	pieces := strings.Fields(prepped)
	name := pieces[1]
	x := parseCoordinateValue(pieces[2])
	y := parseCoordinateValue(pieces[3])
	return button{
		name: name,
		loc: shared.Loc{
			X: x,
			Y: y,
		},
	}
}

func parseCoordinateValue(s string) int {
	i, err := strconv.Atoi(s[1:])
	if err != nil {
		shared.Logger.Error("Failed to parse coordinate.", "s", s, "err", err)
		panic("Failed to parse coordinate.")
	}
	return i
}

func parsePrize(line string, fixConversion bool) shared.Loc {
	prepped := strings.ReplaceAll(strings.ReplaceAll(line, ":", ""), ",", "")
	pieces := strings.Fields(prepped)
	loc := shared.Loc{X: -1, Y: -1}
	for _, each := range pieces[1:] {
		if strings.HasPrefix(each, "X=") {
			i, err := strconv.Atoi(each[2:])
			if err != nil {
				shared.Logger.Error(
					"Failed to parse X coordinate from prize.",
					"piece",
					each,
					"error",
					err,
				)
				panic("Failed to parse X coordinate from prize.")
			}
			loc.X = i
		}
		if strings.HasPrefix(each, "Y=") {
			i, err := strconv.Atoi(each[2:])
			if err != nil {
				shared.Logger.Error(
					"Failed to parse Y coordinate from prize.",
					"piece",
					each,
					"error",
					err,
				)
				panic("Failed to parse Y coordinate from prize.")
			}
			loc.Y = i
		}
	}
	if !fixConversion {
		return loc
	}
	return shared.Loc{
		X: loc.X + conversionFix,
		Y: loc.Y + conversionFix,
	}
}

func deriveCheapest(a, b button, prize shared.Loc) int {
	newRat := func(n, d int) *big.Rat {
		return big.NewRat(int64(n), int64(d))
	}
	// The following simultaneous equations are at play here:
	// axi + bxj = px
	// ayi + byi = px
	// When solved for i, they become:
	// i = (px * -(by/bx) + py) / (ay - (by/bx) * ax)
	px := newRat(prize.X, 1)
	byBx := newRat(-b.loc.Y, b.loc.X)
	py := newRat(prize.Y, 1)
	ay := newRat(a.loc.Y, 1)
	ax := newRat(a.loc.X, 1)

	nom := big.NewRat(1, 1).Mul(px, byBx)
	nom.Add(nom, py)

	den := big.NewRat(1, 1).Mul(byBx, ax)
	den.Add(ay, den)

	i := big.NewRat(1, 1).Mul(nom, big.NewRat(1, 1).Inv(den))
	if !i.IsInt() {
		return -1
	}

	// And once we know i, solving j is easy:
	// axi + bxj = px
	// becomes
	// j = (px - axi) / bx
	above := newRat(1, 1).Sub(px, newRat(1, 1).Mul(ax, i))
	j := above.Mul(above, newRat(1, 1).Inv(newRat(b.loc.X, 1)))
	if !j.IsInt() {
		return -1
	}

	result := newRat(1, 1).Mul(newRat(3, 1), i)
	result.Add(result, j)
	iNum := result.Num()
	iDem := result.Denom()
	return int(iNum.Div(iNum, iDem).Int64())
}
