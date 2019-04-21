package main

import (
	"encoding/csv"
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/cpu"
	"io"
	"os"
	"regexp"
	"strconv"
)

func main() {
	file, err := os.Open("assets/6502-opcode-matrix.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fout, err := os.Create("assets/6502-opcode-matrix.go.txt")
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	r, _ := regexp.Compile(`^(\w+)(?:\s+([a-z]+))?(?:\s+(\d+))?(\*)?`)

	reader := csv.NewReader(file)
	row := 0
	for rec, err := reader.Read(); rec != nil; rec, err = reader.Read() {
		if err != nil && err != io.EOF {
			panic(err)
		}
		//fmt.Printf("got row %x: %v\n", row, rec)
		if row == 0 {
			row++
			continue
		}
		opcode := (row - 1) << 4
		for col, val := range rec {
			if col == 0 {
				continue
			}
			opcode = ((row - 1) << 4) | (col - 1)
			fmt.Printf("got instruction %02x at %x, %x: %s \n", opcode, row-1, col-1, val)
			m := r.FindStringSubmatch(val)
			inst := &cpu.InstructionInfo{}
			inst.OpCode = byte(opcode)
			inst.Nemonics = m[1]
			switch m[2] {
			case "":
				inst.AddressingMode = cpu.IMP
			case "imm":
				inst.AddressingMode = cpu.IMM
			case "izx":
				inst.AddressingMode = cpu.IZX
			case "izy":
				inst.AddressingMode = cpu.IZY
			case "zp":
				inst.AddressingMode = cpu.ZP
			case "zpx":
				inst.AddressingMode = cpu.ZPX
			case "zpy":
				inst.AddressingMode = cpu.ZPY
			case "abs":
				inst.AddressingMode = cpu.ABS
			case "abx":
				inst.AddressingMode = cpu.ABX
			case "aby":
				inst.AddressingMode = cpu.ABY
			case "rel":
				inst.AddressingMode = cpu.REL
			case "ind":
				inst.AddressingMode = cpu.IND
			default:
				panic(fmt.Errorf("uknown addressing mode %s", m[2]))
			}
			inst.Cycles, _ = strconv.Atoi(m[3])
			inst.VariableCycles = m[4] != ""

			fmt.Printf("opcode=%02x: %v\n", opcode, inst)
			fmt.Fprintf(fout, "0x%02x: {0x%02x, \"%s\", %s, %d, %v},\n",
				inst.OpCode, inst.OpCode, inst.Nemonics, inst.AddressingMode.String(), inst.Cycles, inst.VariableCycles)
		}
		row++
	}
}
