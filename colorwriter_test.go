package main

import (
	"reflect"
	"testing"

	"github.com/maxzender/jsonexplorer/formatter"
	"github.com/maxzender/jsonexplorer/treemodel"
	"github.com/nsf/termbox-go"
)

var testColorMap = map[formatter.TokenType]termbox.Attribute{
	formatter.DelimiterType: termbox.ColorWhite,
	formatter.BoolType:      termbox.ColorBlue,
	formatter.StringType:    termbox.ColorRed,
	formatter.NumberType:    termbox.ColorYellow,
	formatter.NullType:      termbox.ColorCyan,
}

var defaultColor = termbox.ColorDefault

func TestWrite(t *testing.T) {
	writer := NewColorWriter(testColorMap, defaultColor)
	writer.Write(`{`, formatter.DelimiterType)
	writer.Newline()
	writer.Write(`    `, formatter.WhiteSpaceType)
	writer.Write(`"test"`, formatter.StringType)
	writer.Write(`:`, formatter.DelimiterType)
	writer.Write(` `, formatter.WhiteSpaceType)
	writer.Write(`4`, formatter.NumberType)
	writer.Newline()
	writer.Write(`}`, formatter.DelimiterType)

	expected := []treemodel.Line{
		treemodel.Line{{'{', termbox.ColorWhite}},
		treemodel.Line{
			{' ', 0},
			{' ', 0},
			{' ', 0},
			{' ', 0},
			{'"', termbox.ColorRed},
			{'t', termbox.ColorRed},
			{'e', termbox.ColorRed},
			{'s', termbox.ColorRed},
			{'t', termbox.ColorRed},
			{'"', termbox.ColorRed},
			{':', termbox.ColorWhite},
			{' ', 0},
			{'4', termbox.ColorYellow},
		},
		treemodel.Line{{'}', termbox.ColorWhite}},
	}
	actual := writer.Lines

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected:\n%v but received:\n%v", expected, actual)
	}
}
