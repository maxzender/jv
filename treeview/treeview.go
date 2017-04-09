package treeview

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

type TreeView struct {
	content       []string
	segments      map[int]int
	expandedLines map[int]struct{}
	lineMap       map[int]int
}

func New(r io.Reader) (*TreeView, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")

	return &TreeView{
		content:       lines,
		segments:      parseSegments(lines),
		expandedLines: map[int]struct{}{0: struct{}{}},
		lineMap:       make(map[int]int),
	}, nil
}

func (v *TreeView) Render() io.Reader {
	var buf bytes.Buffer
	skipTill := 0
	virtualLn := 0
	for actualLn, val := range v.content {
		if actualLn < skipTill {
			continue
		}

		_, lineExpanded := v.expandedLines[actualLn]
		r, _ := utf8.DecodeLastRuneInString(val)
		if !lineExpanded && (r == '{' || r == '[') {
			closingBraceLine := v.segments[actualLn]
			skipTill = closingBraceLine + 1
			closingBrace := strings.TrimSpace(v.content[closingBraceLine])
			buf.WriteString(fmt.Sprintf("%s...%s\n", val, closingBrace))
		} else {
			buf.WriteString(fmt.Sprintf("%s\n", val))
		}

		v.lineMap[virtualLn] = actualLn
		virtualLn++
	}

	return &buf
}

func (v *TreeView) ToggleLine(virtualLn int) {
	actualLn := v.lineMap[virtualLn]
	_, isExpanded := v.expandedLines[actualLn]
	if isExpanded {
		delete(v.expandedLines, actualLn)
	} else {
		v.expandedLines[actualLn] = struct{}{}
	}
}

func parseSegments(lines []string) map[int]int {
	resultSegments := make(map[int]int)
	bracketBalances := make(map[int]int)
	var bal int
	for num, line := range lines {
		for _, c := range line {
			switch c {
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
