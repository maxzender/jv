package treemodel

import (
	"unicode"

	"github.com/nsf/termbox-go"
)

type TreeModel struct {
	lines         []Line
	expandedLines map[int]struct{}
	lineMap       map[int]int
	segments      map[int]int
}

type Char struct {
	Val   rune
	Color termbox.Attribute
}

type Line []Char

func New(lines []Line) *TreeModel {
	model := &TreeModel{
		lines:         lines,
		expandedLines: map[int]struct{}{},
		lineMap:       make(map[int]int),
		segments:      parseSegments(lines),
	}
	model.recalculateLineMap()
	model.ToggleLine(0)

	return model
}

func (v *TreeModel) ToggleLine(virtualLn int) {
	actualLn := v.lineMap[virtualLn]
	if !v.isBeginningOfSegment(actualLn) {
		return
	}

	if v.isExpanded(actualLn) {
		delete(v.expandedLines, actualLn)
	} else {
		v.expandedLines[actualLn] = struct{}{}
	}
	v.recalculateLineMap()
}

func (v *TreeModel) Line(virtualLn int) Line {
	actualLn, ok := v.lineMap[virtualLn]
	if ok {
		ln := v.lines[actualLn]
		if v.isBeginningOfSegment(actualLn) && !v.isExpanded(actualLn) {
			ln = v.lineWithDots(actualLn)
		}
		return ln
	} else {
		return nil
	}
}

func (v *TreeModel) lineWithDots(actualLn int) Line {
	ln := v.lines[actualLn]

	lastChar := ln[len(ln)-1]
	ln = append(ln, Char{'…', lastChar.Color})

	matchingBrace := v.lines[v.segments[actualLn]]
	for _, c := range matchingBrace {
		if !unicode.IsSpace(c.Val) {
			ln = append(ln, c)
		}
	}

	return ln
}

func (v *TreeModel) isExpanded(actualLn int) bool {
	_, isExpanded := v.expandedLines[actualLn]
	return isExpanded
}

func (v *TreeModel) isBeginningOfSegment(actualLn int) bool {
	_, ok := v.segments[actualLn]
	return ok
}

func (v *TreeModel) recalculateLineMap() {
	v.lineMap = make(map[int]int)
	skipTill := 0
	virtualLn := 0
	for actualLn, _ := range v.lines {
		if actualLn < skipTill {
			continue
		}

		if v.isBeginningOfSegment(actualLn) && !v.isExpanded(actualLn) {
			skipTill = v.segments[actualLn] + 1
		}

		v.lineMap[virtualLn] = actualLn
		virtualLn++
	}
}

func parseSegments(lines []Line) map[int]int {
	resultSegments := make(map[int]int)
	bracketBalances := make(map[int]int)
	var bal int
	for num, line := range lines {
		for _, c := range line {
			switch c.Val {
			case '{', '[':
				bal++
				bracketBalances[bal] = num
			case '}', ']':
				resultSegments[bracketBalances[bal]] = num
				bal--
			}
		}
	}

	return resultSegments
}
