package parser

import (
	"fmt"
	"strings"
)

const INDENT = "  "

func MulIndent(mul int) string {
	var b strings.Builder
	for range mul {
		b.WriteString(INDENT)
	}
	return b.String()
}

func LineWithDepth(line string, depth int) string {
	return MulIndent(depth) + line
}
func LinesWithDepth(lines []string, depth int) []string {
	var stringsWithNewDepth []string
	for _, line := range lines {
		stringsWithNewDepth = append(stringsWithNewDepth, MulIndent(depth)+line)
	}
	return stringsWithNewDepth
}

func JoinStringRest(strings ...string) string {
	return JoinWithSep(strings, " ")
}
func JoinBuilderG[V fmt.Stringer](ss []V) string {
	var b strings.Builder
	for _, v := range ss {
		b.WriteString(v.String())
	}
	return b.String()
}
func JoinBuilder(ss []string) string {
	var b strings.Builder
	for _, v := range ss {
		b.WriteString(v)
	}
	return b.String()
}

func JoinWithSepG[V fmt.Stringer](ss []V, sep string) string {
	var b strings.Builder
	for i, v := range ss {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(v.String())
	}
	return b.String()
}

func JoinWithSep(ss []string, sep string) string {
	var b strings.Builder
	for i, v := range ss {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(v)
	}
	return b.String()
}

func JoinLinesG[V fmt.Stringer](ss []V) string {
	var b strings.Builder
	for _, v := range ss {
		b.WriteString(v.String())
		b.WriteByte('\n')
	}
	return b.String()
}
func JoinLines(ss []string) string {
	var b strings.Builder
	for _, v := range ss {
		b.WriteString(v)
		b.WriteByte('\n')
	}
	return b.String()
}

func JoinLinesWithSepG[V fmt.Stringer](ss []V, sep string) string {
	var b strings.Builder
	for _, v := range ss {
		n := sep + v.String()
		b.WriteString(n)
		b.WriteByte('\n')
	}
	return b.String()
}

const Indent string = "	"

func JoinLineIndentG[V fmt.Stringer](ss []V) string {
	return JoinLinesWithSepG(ss, Indent)
}
func JoinLinesWithSep(ss []string, sep string) string {
	var b strings.Builder
	for _, str := range ss {
		line := sep + str
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func JoinLinesWithIndent(ss []string) string {
	return JoinLinesWithSep(ss, Indent)
}
