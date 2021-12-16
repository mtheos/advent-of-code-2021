package packetDecoder

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type packet struct {
	version, typeId int
	literal         uint64
	packets         []packet
}

func readInput(fileName string) packet {
	var packets []packet
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		scanner.Scan()
		hex := scanner.Text()
		binary := hex2bin(hex)
		packets = parse(&binary, -1)
	})
	return packets[0]
}

func parse(binary *string, subPacketCount int) []packet {
	var packets []packet
	for {
		if subPacketCount != -1 {
			subPacketCount--
		}
		// arbitrary choice
		if len(*binary) == 0 || (len(*binary) < 16 && bin2Int(*binary, len(*binary)) == 0) {
			break
		}
		version := bin2Int((*binary)[:3], 3)
		typeId := bin2Int((*binary)[3:6], 3)
		*binary = (*binary)[6:]
		var p packet
		switch typeId {
		case LITERAL:
			p = parseLiteralPacket(int(version), int(typeId), binary)
		default:
			p = parseOperatorPacket(int(version), int(typeId), binary)
		}
		packets = append(packets, p)
		if subPacketCount == 0 {
			break
		}
	}
	return packets
}

func parseOperatorPacket(version int, typeId int, binary *string) packet {
	var packets []packet
	lengthTypeId := (*binary)[0]
	*binary = (*binary)[1:]
	if lengthTypeId == '0' {
		length := bin2Int((*binary)[:15], 15)
		*binary = (*binary)[15:]
		subStr := (*binary)[:length]
		*binary = (*binary)[length:]
		packets = parse(&subStr, -1)
	} else {
		numPackets := bin2Int((*binary)[:11], 11)
		*binary = (*binary)[11:]
		packets = append(packets, parse(binary, int(numPackets))...)
	}
	return packet{version, typeId, 0, packets}
}

func parseLiteralPacket(version int, typeId int, binary *string) packet {
	literalBits := ""
	for {
		sign, chunk := (*binary)[0], (*binary)[1:5]
		*binary = (*binary)[5:]
		literalBits += chunk
		if sign == '0' {
			break
		}
	}
	literal := bin2Int(literalBits, len(literalBits))
	return packet{version, typeId, literal, nil}
}

func hex2bin(raw string) string {
	const nibble = 4
	const maxChunk = 64 / nibble
	sb := strings.Builder{}
	for len(raw) > 0 {
		chunk := ""
		if len(raw) > maxChunk {
			chunk = raw[:maxChunk]
			raw = raw[maxChunk:]
		} else {
			chunk = raw
			raw = ""
		}
		binary, err := strconv.ParseUint(chunk, 16, len(chunk)*nibble)
		MaybePanic(err)
		sb.WriteString(fmt.Sprintf("%0*b", len(chunk)*nibble, binary))
	}
	return sb.String()
}

func bin2Int(s string, size int) uint64 {
	i, err := strconv.ParseUint(s, 2, size)
	MaybePanic(err)
	return i
}

func sumVersions(packets []packet) int {
	var sum int
	for _, p := range packets {
		switch p.typeId {
		case LITERAL:
			sum += p.version
		default:
			sum += p.version + sumVersions(p.packets)
		}
	}
	return sum
}

func ezMode(p packet, ch chan<- int) {
	ch <- sumVersions([]packet{p})
}

func hardMode(p packet, ch chan<- uint64) {
	ch <- eval(p)
}

func Go(fileName string, ch chan string) {
	p := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan uint64)

	go ezMode(p, ezChan)
	go hardMode(p, hardChan)

	ch <- fmt.Sprintln("packet Decoder")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
