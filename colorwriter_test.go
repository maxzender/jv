package main

import (
	"reflect"
	"testing"

	"github.com/maxzender/jsonexplorer/treemodel"
	"github.com/nsf/termbox-go"
)

var testColorMap = map[jsonfmt.TokenType]termbox.Attribute{
	jsonfmt.DelimiterType: termbox.ColorWhite,
	jsonfmt.BoolType:      termbox.ColorBlue,
	jsonfmt.StringType:    termbox.ColorRed,
	jsonfmt.NumberType:    termbox.ColorYellow,
	jsonfmt.NullType:      termbox.ColorCyan,
}

var defaultColor = termbox.ColorDefault

func TestWrite(t *testing.T) {
	writer := NewColorWriter(testColorMap, defaultColor)
	writer.Write(`{`, jsonfmt.DelimiterType)
	writer.Newline()
	writer.Write(`    `, jsonfmt.WhiteSpaceType)
	writer.Write(`"test"`, jsonfmt.StringType)
	writer.Write(`:`, jsonfmt.DelimiterType)
	writer.Write(` `, jsonfmt.WhiteSpaceType)
	writer.Write(`4`, jsonfmt.NumberType)
	writer.Newline()
	writer.Write(`}`, jsonfmt.DelimiterType)

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
