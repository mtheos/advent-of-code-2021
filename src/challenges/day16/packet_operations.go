package packetDecoder

import "math"

const (
	SUM     = 0
	PRODUCT = 1
	MIN     = 2
	MAX     = 3
	LITERAL = 4
	GT      = 5
	LT      = 6
	EQ      = 7
)

var ops map[int]func([]packet) uint64

func init() {
	ops = map[int]func([]packet) uint64{
		SUM:     sumOp,
		PRODUCT: productOp,
		MIN:     minOp,
		MAX:     maxOp,
		LITERAL: panicOp,
		GT:      gtOp,
		LT:      ltOp,
		EQ:      eqOp,
	}
}

func eval(p packet) uint64 {
	switch p.typeId {
	case LITERAL:
		return p.literal
	default:
		return ops[p.typeId](p.packets)
	}
}

// SUM - their value is the sum of the values of their sub-packets.
// If they only have a single sub-packet, their value is the value of the sub-packet.
func sumOp(packets []packet) uint64 {
	var acc uint64
	for _, p := range packets {
		acc += eval(p)
	}
	return acc
}

// PRODUCT - their value is the result of multiplying together the values of their sub-packets.
// If they only have a single sub-packet, their value is the value of the sub-packet.
func productOp(packets []packet) uint64 {
	var acc uint64 = 1
	for _, p := range packets {
		acc *= eval(p)
	}
	return acc
}

// MIN - their value is the minimum of the values of their sub-packets.
func minOp(packets []packet) uint64 {
	var acc uint64 = math.MaxUint64
	for _, p := range packets {
		acc = uint64(math.Min(float64(acc), float64(eval(p))))
	}
	return acc
}

// MAX - their value is the maximum of the values of their sub-packets.
func maxOp(packets []packet) uint64 {
	var acc uint64
	for _, p := range packets {
		acc = uint64(math.Max(float64(acc), float64(eval(p))))
	}
	return acc
}

// LITERAL - represents a literal value, not an operation.
func panicOp([]packet) uint64 {
	panic("I should not be called")
}

// GT - their value is 1 if the value of the first sub-packet is greater than
// the value of the second sub-packet; otherwise, their value is 0.
// These packets always have exactly two sub-packets.
func gtOp(packets []packet) uint64 {
	first, second := eval(packets[0]), eval(packets[1])
	if first > second {
		return 1
	}
	return 0
}

// LT - their value is 1 if the value of the first sub-packet is less than
// the value of the second sub-packet; otherwise, their value is 0.
// These packets always have exactly two sub-packets.
func ltOp(packets []packet) uint64 {
	first, second := eval(packets[0]), eval(packets[1])
	if first < second {
		return 1
	}
	return 0
}

// EQ - their value is 1 if the value of the first sub-packet is equal to
// the value of the second sub-packet; otherwise, their value is 0.
// These packets always have exactly two sub-packets.
func eqOp(packets []packet) uint64 {
	first, second := eval(packets[0]), eval(packets[1])
	if first == second {
		return 1
	}
	return 0
}
