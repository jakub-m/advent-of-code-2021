package main

import (
	"advent"
	"fmt"
	"strconv"
	"strings"
)

func Calc(rawHex string) (int, error) {
	bin, err := hexToBin(rawHex)
	advent.Printf("hex %s\n", rawHex)
	advent.Printf("binary %s\n", bin)
	if err != nil {
		return 0, err
	}
	p, _ := decodePacketRawBin(bin)
	advent.Printf("%v", p)

	var sumVerRec func(p packet) int

	sumVerRec = func(p packet) int {
		switch p := p.(type) {
		case literalPacket:
			return p.getVersion()
		case operatorPacket:
			v := p.version
			for _, ch := range p.packets {
				v += sumVerRec(ch)
			}
			return v
		default:
			panic(p)
		}
	}

	return sumVerRec(p), nil
}

type packet interface {
	getVersion() int
}

const (
	Literal = 4
)

func getVersionAndType(rawBin string) (int, int, string) {
	version := binToInt(rawBin[0:3])
	rawBin = rawBin[3:]
	typeId := binToInt(rawBin[0:3])
	rawBin = rawBin[3:]
	switch typeId {
	case Literal:
		return version, Literal, rawBin
	default:
		// operator
		return version, typeId, rawBin
	}
}

func decodePacketRawBin(rawBin string) (packet, string) {
	_, typeId, _ := getVersionAndType(rawBin)
	if typeId == Literal {
		return decodeLiteral(rawBin)
	} else {
		return decodeOperator(rawBin)
	}
}

type literalPacket struct {
	version int
	value   int
}

func (p literalPacket) getVersion() int {
	return p.version
}

func (p literalPacket) String() string {
	return fmt.Sprintf("lit:%d,%d", p.version, p.value)
}

func decodeLiteral(rawBin string) (literalPacket, string) {
	advent.Printf("decodeLiteral %s\n", rawBin)
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
	//val := binToInt(subBin)
	val := -1
	p := literalPacket{version: version, value: val}
	return p, rawBin
}

type operatorPacket struct {
	version int
	packets []packet
}

func (p operatorPacket) getVersion() int {
	return p.version
}

func (p operatorPacket) String() string {
	packets := []string{}
	for _, p := range p.packets {
		packets = append(packets, fmt.Sprint(p))
	}
	return fmt.Sprintf("op:%d,[%s]", p.version, strings.Join(packets, ", "))
}

func decodeOperator(rawBin string) (operatorPacket, string) {
	advent.Printf("decodeOperator %s\n", rawBin)
	version, _, rawBin := getVersionAndType(rawBin)
	lengthTypeId, rawBin := rawBin[0], rawBin[1:]

	var finalRawBin string
	packets := []packet{}
	if lengthTypeId == '0' {
		// next 15 are tot length of other subpackets
		s, rawBin := rawBin[0:15], rawBin[15:]
		nbitsOfNextPackets := binToInt(s)
		advent.Printf("nbitsOfNextPackets %s %d\n", s, nbitsOfNextPackets)
		breakAtRawbinLen := len(rawBin) - nbitsOfNextPackets
		advent.Assertf(breakAtRawbinLen > 0, "bullshit expected length: %d", breakAtRawbinLen)
		for {
			p, rawBin2 := decodePacketRawBin(rawBin)
			rawBin = rawBin2
			packets = append(packets, p)
			if len(rawBin) == breakAtRawbinLen {
				finalRawBin = rawBin
				break
			}
			advent.Assertf(len(rawBin) >= breakAtRawbinLen, "Expected to finish at len %d but is already %d", len(rawBin), breakAtRawbinLen)
		}
	} else {
		// 11 bits are num of subpackets
		s, rawBin := rawBin[0:11], rawBin[11:]
		numOfNextPackets := binToInt(s)
		advent.Printf("numOfNextPackets %d\n", numOfNextPackets)
		for {
			p, rawBin2 := decodePacketRawBin(rawBin)
			rawBin = rawBin2
			packets = append(packets, p)
			if len(packets) == numOfNextPackets {
				finalRawBin = rawBin
				break
			}
		}
	}

	return operatorPacket{version: version, packets: packets}, finalRawBin
}

func binToInt(bin string) int {
	i, err := strconv.ParseInt(bin, 2, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}

func hexToBin(hex string) (string, error) {
	b := ""
	for _, h := range hex {
		d, err := strconv.ParseInt(fmt.Sprintf("%c", h), 16, 32)
		if err != nil {
			return "", err
		}
		b += fmt.Sprintf("%04s", strconv.FormatInt(d, 2))
	}
	return b, nil
}
