package aoc2417

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

const (
	instAdv = 0
	instBxl = 1
	instBst = 2
	instJnz = 3
	instBxc = 4
	instOut = 5
	instBdv = 6
	instCdv = 7
)

func DeriveOutput(lines []string) string {
	shared.Logger.Info("Derive program output.")
	if len(lines) == 0 {
		shared.Logger.Info("Empty lines so quitting.")
		return ""
	}

	cpu := newProcessor(lines)
	runBatch(cpu)
	shared.Logger.Info("Ints in output.", "ints", cpu.output)
	return strings.Join(
		shared.MapValues(
			cpu.output,
			func(i int) string {
				return strconv.Itoa(i)
			},
		),
		",")
}

func runBatch(cpu *processor) {
	ongoing := true
	for ongoing {
		ongoing = cpu.process()
	}
}

type processor struct {
	registers []int
	index     int
	feed      []int
	output    []int
}

func newProcessor(lines []string) *processor {
	registers := []int{}
	for i := 0; i < 3; i++ {
		registers = append(registers, shared.ParseIntOrDie(strings.Fields(lines[i])[2]))
	}

	for _, each := range lines {
		if strings.HasPrefix(each, "Program:") {
			pieces := strings.Split(strings.Fields(each)[1], ",")
			ints, err := shared.ToInts(pieces)
			shared.Die(err, "Failed to parse program feed ints.")
			return &processor{
				registers: registers,
				index:     0,
				feed:      ints,
			}
		}
	}

	panic("Failed to create processor.")
}

func (v *processor) process() bool {
	remaining := len(v.feed) - v.index
	if remaining < 2 {
		if remaining == 1 {
			shared.Logger.Error("Invalid feed state.", "processor", v)
			panic("Remaining should be divisable by 2.")
		}
		return false
	}

	opcode, operand := v.popFeed()
	increment := true
	switch opcode {
	case instAdv:
		v.adv(operand)
	case instBxl:
		v.bxl(operand)
	case instBst:
		v.bst(operand)
	case instJnz:
		increment = v.jnz(operand)
	case instBxc:
		v.bxc()
	case instOut:
		v.out(operand)
	case instBdv:
		v.bdv(operand)
	case instCdv:
		v.cdv(operand)
	default:
		panic(fmt.Sprintf("Unknown opcode: %d.", opcode))
	}

	if increment {
		v.index += 2
	}
	return true
}

func (v *processor) popFeed() (opcode int, operand int) {
	opcode = v.feed[v.index]
	operand = v.feed[v.index+1]
	return
}

func (v *processor) adv(operand int) {
	v.registers[0] = v.calculateDv(operand)
}

func (v *processor) bxl(operand int) {
	result := v.registers[1] ^ operand
	v.registers[1] = result
}

func (v *processor) bst(operand int) {
	value := v.deriveOperand(operand) % 8
	v.registers[1] = value
}

func (v *processor) jnz(operand int) bool {
	if v.registers[0] == 0 {
		return true
	}
	v.index = operand

	return false
}

func (v *processor) bxc() {
	v.registers[1] = v.registers[1] ^ v.registers[2]
}

func (v *processor) out(operand int) {
	value := v.deriveOperand(operand) % 8
	v.output = append(v.output, value)
}

func (v *processor) bdv(operand int) {
	v.registers[1] = v.calculateDv(operand)
}

func (v *processor) cdv(operand int) {
	v.registers[2] = v.calculateDv(operand)
}

func (v *processor) calculateDv(operand int) int {
	value := v.deriveOperand(operand)
	num := v.registers[0]
	den := shared.Pow(2, value)
	return num / den
}

func (v *processor) deriveOperand(i int) int {
	if 0 <= i && i <= 3 {
		return i
	}
	if i >= 7 {
		shared.Logger.Error("Invalid combo operand.", "operand", i)
		panic("Invalid combo operand.")
	}
	return v.registers[i-4]
}
