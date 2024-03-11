package parser

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/oktavarium/sgs/internal/internal/vector"
)

func ParseSequenceNumber(data string) (int, error) {
	var seqNumber int
	rgxp, err := regexp.Compile("(?<seqNumber>\\d+) c\\d+t")
	if err != nil {
		return seqNumber, fmt.Errorf("compile regexp: %w", err)
	}
	matches := rgxp.FindStringSubmatch(data)
	if matches == nil {
		return seqNumber, fmt.Errorf("sequence number not found")
	}
	seqNumber, err = strconv.Atoi(matches[rgxp.SubexpIndex("seqNumber")])
	if err != nil {
		return seqNumber, fmt.Errorf("strconv atoi: %w", err)
	}

	return seqNumber, nil
}

func ParseID(data string) (string, error) {
	rgxp, err := regexp.Compile("(?<id>c\\d+t)")
	if err != nil {
		return "", fmt.Errorf("compile regexp: %w", err)
	}
	matches := rgxp.FindStringSubmatch(data)
	if matches == nil {
		return "", fmt.Errorf("sequence number not found")
	}

	return matches[rgxp.SubexpIndex("id")], nil
}

func ParseInitialPosition(data string) (vector.Vector, error) {
	var pos vector.Vector
	rgxp, err := regexp.Compile("n (?<x>-?([0-9]*[.])?[0-9]+) (?<y>-?([0-9]*[.])?[0-9]+) (?<z>-?([0-9]*[.])?[0-9]+)")
	if err != nil {
		return pos, fmt.Errorf("compile regexp: %w", err)
	}
	matches := rgxp.FindStringSubmatch(data)
	if matches == nil {
		return pos, fmt.Errorf("sequence number not found")
	}
	pos, err = ParsePosition(
		matches[rgxp.SubexpIndex("x")],
		matches[rgxp.SubexpIndex("y")],
		matches[rgxp.SubexpIndex("z")],
	)
	if err != nil {
		return pos, fmt.Errorf("parse position: %w", err)
	}

	return pos, nil
}

func ParsePosition(x, y, z string) (vector.Vector, error) {
	var vec vector.Vector
	var err error
	vec.X, err = strconv.ParseFloat(x, 64)
	if err != nil {
		return vec, fmt.Errorf("parse x: %w", err)
	}
	vec.Y, err = strconv.ParseFloat(y, 64)
	if err != nil {
		return vec, fmt.Errorf("parse y: %w", err)
	}
	vec.Z, err = strconv.ParseFloat(z, 64)
	if err != nil {
		return vec, fmt.Errorf("parse z: %w", err)
	}

	return vec, nil
}

func ParseInput(data string) (string, error) {
	rgxp, err := regexp.Compile("c\\d+t (?<input>[adws])")
	if err != nil {
		return "", fmt.Errorf("compile regexp: %w", err)
	}
	matches := rgxp.FindStringSubmatch(data)
	if matches == nil {
		return "", fmt.Errorf("sequence number not found")
	}

	return matches[rgxp.SubexpIndex("input")], nil
}
