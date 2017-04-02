package contentview

type ContentView struct {
	OriginalContent []string
	Segments        map[int]int
	ExpandedLines   map[int]struct{}
	LineMap         map[int]int
}

func New(content []string) *ContentView {
	return &ContentView{
		OriginalContent: content,
		Segments:        parseSegments(content),
		ExpandedLines:   map[int]struct{}{0: struct{}{}},
		LineMap:         make(map[int]int),
	}
}

func (v *ContentView) Content() []string {
	var result []string
	skipTill := 0
	virtualLine := 0
	for line, val := range v.OriginalContent {
		if line < skipTill {
			continue
		}

		_, lineExpanded := v.ExpandedLines[line]
		if !lineExpanded && line > skipTill {
			skipTill = v.Segments[line]
		}

		v.LineMap[virtualLine] = line
		virtualLine++
		result = append(result, val)
	}
	return result
}

func (v *ContentView) ToggleLine(virtualLine int) {
	actualLine := v.LineMap[virtualLine]
	_, isExpanded := v.ExpandedLines[actualLine]
	if isExpanded {
		delete(v.ExpandedLines, actualLine)
	} else {
		v.ExpandedLines[actualLine] = struct{}{}
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
