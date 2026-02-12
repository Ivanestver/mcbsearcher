package main

import (
	"cycles/cycles_alg"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	lammps_parser "github.com/Ivanestver/lammps-file-parser/parser"
)

func main() {
	infilePtr := flag.String("infile", "", "Specifies the input file")
	flag.Parse()
	if len(*infilePtr) == 0 {
		fmt.Println("Wrong usage of the infile parameter")
		return
	}
	content, err := os.ReadFile(*infilePtr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lammpsStruct, err := lammps_parser.Parse(
		string(content),
		*infilePtr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonObj, err := json.Marshal(*lammpsStruct)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cycles, err := cycles_alg.CalculateCycles(string(jsonObj))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	printCycles(cycles)
}

func printCycles(cycles []cycles_alg.Cycle) {
	for i, cycle := range cycles {
		builder := strings.Builder{}
		builder.WriteString("C")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(": ")
		for _, edge := range cycle {
			builder.WriteString(fmt.Sprintf("%v, ", edge))
		}
		builder.WriteString("\n\n")
		fmt.Println(builder.String())
	}
}
