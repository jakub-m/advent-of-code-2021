package main

import (
	"advent"
	"fmt"
	"strconv"
)

func Calc(raw string) (int, error) {

	p := decodePacketRawBin(raw)
	advent.Printf("%v", p)

	return 0, nil
}

type packet struct {
	version int
	value   int
}

const Literal = 4

func getVersionAndType(rawBin string) (int, int, string) {
	version := binToInt(rawBin[0:3])
	rawBin = rawBin[3:]
	typeId := binToInt(rawBin[0:3])
	rawBin = rawBin[3:]
	switch typeId {
	case Literal:
		return version, Literal, rawBin
	default:
		panic(fmt.Sprintf("unexpected type for: %s", rawBin))
	}
}

func decodePacketRawBin(rawBin string) packet {
	_, typeId, _ := getVersionAndType(rawBin)
	if typeId == Literal {
		p, _ := decodeLiteral(rawBin)
		return p
	} else {
		panic(fmt.Sprintf("expected type id 4: %s", rawBin))
	}
}

func decodeLiteral(rawBin string) (packet, string) {
	version, typeId, rawBin := getVersionAndType(rawBin)
	advent.Assertf(typeId == Literal, "unexpected type: %d, %s", typeId, rawBin)
	subBin := ""
	for {
		subBin += rawBin[1:5]
		// fmt.Printf("subbin %s\n", subBin)
		shouldBreak := rawBin[0] == '0'
		rawBin = rawBin[5:]
		if shouldBreak {
			break
		}
	}
	val := binToInt(subBin)
	p := packet{version: version, value: val}
	return p, rawBin
}

func binToInt(bin string) int {
	i, err := strconv.ParseInt(bin, 2, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}
