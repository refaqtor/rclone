// +build go1.10

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/artpar/rclone/lib/encoder"
)

const (
	edgeLeft = iota
	edgeRight
)

type mapping struct {
	mask     uint
	src, dst []rune
}
type stringPair struct {
	a, b string
}

const header = `// Code generated by ./internal/gen/main.go. DO NOT EDIT.

` + `//go:generate go run ./internal/gen/main.go

package encoder

`

var maskBits = []struct {
	mask uint
	name string
}{
	{encoder.EncodeZero, "EncodeZero"},
	{encoder.EncodeWin, "EncodeWin"},
	{encoder.EncodeSlash, "EncodeSlash"},
	{encoder.EncodeBackSlash, "EncodeBackSlash"},
	{encoder.EncodeHashPercent, "EncodeHashPercent"},
	{encoder.EncodeDel, "EncodeDel"},
	{encoder.EncodeCtl, "EncodeCtl"},
	{encoder.EncodeLeftSpace, "EncodeLeftSpace"},
	{encoder.EncodeLeftTilde, "EncodeLeftTilde"},
	{encoder.EncodeRightSpace, "EncodeRightSpace"},
	{encoder.EncodeRightPeriod, "EncodeRightPeriod"},
	{encoder.EncodeInvalidUtf8, "EncodeInvalidUtf8"},
}
var edges = []struct {
	mask    uint
	name    string
	edge    int
	orig    rune
	replace rune
}{
	{encoder.EncodeLeftSpace, "EncodeLeftSpace", edgeLeft, ' ', '␠'},
	{encoder.EncodeLeftTilde, "EncodeLeftTilde", edgeLeft, '~', '～'},
	{encoder.EncodeRightSpace, "EncodeRightSpace", edgeRight, ' ', '␠'},
	{encoder.EncodeRightPeriod, "EncodeRightPeriod", edgeRight, '.', '．'},
}

var allMappings = []mapping{{
	encoder.EncodeZero, []rune{
		0,
	}, []rune{
		'␀',
	}}, {
	encoder.EncodeWin, []rune{
		':', '?', '"', '*', '<', '>', '|',
	}, []rune{
		'：', '？', '＂', '＊', '＜', '＞', '｜',
	}}, {
	encoder.EncodeSlash, []rune{
		'/',
	}, []rune{
		'／',
	}}, {
	encoder.EncodeBackSlash, []rune{
		'\\',
	}, []rune{
		'＼',
	}}, {
	encoder.EncodeHashPercent, []rune{
		'#', '%',
	}, []rune{
		'＃', '％',
	}}, {
	encoder.EncodeDel, []rune{
		0x7F,
	}, []rune{
		'␡',
	}}, {
	encoder.EncodeCtl,
	runeRange(0x01, 0x1F),
	runeRange('␁', '␟'),
}}

var (
	rng = rand.New(rand.NewSource(42))

	printables          = runeRange(0x20, 0x7E)
	fullwidthPrintables = runeRange(0xFF00, 0xFF5E)
	encodables          = collectEncodables(allMappings)
	encoded             = collectEncoded(allMappings)
	greek               = runeRange(0x03B1, 0x03C9)
)

func main() {
	fd, err := os.Create("encoder_cases_test.go")
	fatal(err, "Unable to open encoder_cases_test.go:")
	defer func() {
		fatal(fd.Close(), "Failed to close encoder_cases_test.go:")
	}()
	fatalW(fd.WriteString(header))("Failed to write header:")

	fatalW(fd.WriteString("var testCasesSingle = []testCase{\n\t"))("Write:")
	_i := 0
	i := func() (r int) {
		r, _i = _i, _i+1
		return
	}
	for _, m := range maskBits {
		if len(getMapping(m.mask).src) == 0 {
			continue
		}
		if _i != 0 {
			fatalW(fd.WriteString(" "))("Write:")
		}
		in, out := buildTestString(
			[]mapping{getMapping(m.mask)},                               // pick
			[]mapping{getMapping(0)},                                    // quote
			printables, fullwidthPrintables, encodables, encoded, greek) // fill
		fatalW(fmt.Fprintf(fd, `{ // %d
		mask: %s,
		in:   %s,
		out:  %s,
	},`, i(), m.name, strconv.Quote(in), strconv.Quote(out)))("Error writing test case:")
	}
	fatalW(fd.WriteString(`
}

var testCasesSingleEdge = []testCase{
	`))("Write:")
	_i = 0
	for _, e := range edges {
		if _i != 0 {
			fatalW(fd.WriteString(" "))("Write:")
		}
		fatalW(fmt.Fprintf(fd, `{ // %d
		mask: %s,
		in:   %s,
		out:  %s,
	},`, i(), e.name, strconv.Quote(string(e.orig)), strconv.Quote(string(e.replace))))("Error writing test case:")
		for _, m := range maskBits {
			if len(getMapping(m.mask).src) == 0 {
				continue
			}
			pairs := buildEdgeTestString(
				e.edge, e.orig, e.replace,
				[]mapping{getMapping(0), getMapping(m.mask)}, // quote
				printables, fullwidthPrintables, encodables, encoded, greek) // fill
			for _, p := range pairs {
				fatalW(fmt.Fprintf(fd, ` { // %d
		mask: %s | %s,
		in:   %s,
		out:  %s,
	},`, i(), m.name, e.name, strconv.Quote(p.a), strconv.Quote(p.b)))("Error writing test case:")
			}
		}
	}
	fatalW(fmt.Fprintf(fd, ` { // %d
		mask: EncodeLeftSpace,
		in:   "  ",
		out:  "␠ ",
	}, { // %d
		mask: EncodeLeftTilde,
		in:   "~~",
		out:  "～~",
	}, { // %d
		mask: EncodeRightSpace,
		in:   "  ",
		out:  " ␠",
	}, { // %d
		mask: EncodeRightPeriod,
		in:   "..",
		out:  ".．",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightPeriod,
		in:   " .",
		out:  "␠．",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightSpace,
		in:   " ",
		out:  "␠",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightSpace,
		in:   "  ",
		out:  "␠␠",
	}, { // %d
		mask: EncodeLeftSpace | EncodeRightSpace,
		in:   "   ",
		out:  "␠ ␠",
	},
}
`, i(), i(), i(), i(), i(), i(), i(), i()))("Error writing test case:")
}

func fatal(err error, s ...interface{}) {
	if err != nil {
		log.Fatalln(append(s, err))
	}
}
func fatalW(_ int, err error) func(...interface{}) {
	if err != nil {
		return func(s ...interface{}) {
			log.Fatalln(append(s, err))
		}
	}
	return func(s ...interface{}) {}
}

// construct a slice containing the runes between (l)ow (inclusive) and (h)igh (inclusive)
func runeRange(l, h rune) []rune {
	if h < l {
		panic("invalid range")
	}
	out := make([]rune, h-l+1)
	for i := range out {
		out[i] = l + rune(i)
	}
	return out
}

func getMapping(mask uint) mapping {
	for _, m := range allMappings {
		if m.mask == mask {
			return m
		}
	}
	return mapping{}
}
func collectEncodables(m []mapping) (out []rune) {
	for _, s := range m {
		for _, r := range s.src {
			out = append(out, r)
		}
	}
	return
}
func collectEncoded(m []mapping) (out []rune) {
	for _, s := range m {
		for _, r := range s.dst {
			out = append(out, r)
		}
	}
	return
}

func buildTestString(mappings, testMappings []mapping, fill ...[]rune) (string, string) {
	combinedMappings := append(mappings, testMappings...)
	var (
		rIn  []rune
		rOut []rune
	)
	for _, m := range mappings {
		if len(m.src) == 0 || len(m.src) != len(m.dst) {
			panic("invalid length")
		}
		rIn = append(rIn, m.src...)
		rOut = append(rOut, m.dst...)
	}
	inL := len(rIn)
	testL := inL * 3
	if testL < 30 {
		testL = 30
	}
	rIn = append(rIn, make([]rune, testL-inL)...)
	rOut = append(rOut, make([]rune, testL-inL)...)
	quoteOut := make([]bool, testL)
	set := func(i int, in, out rune, quote bool) {
		rIn[i] = in
		rOut[i] = out
		quoteOut[i] = quote
	}
	for i, r := range rOut[:inL] {
		set(inL+i, r, r, true)
	}

outer:
	for pos := inL * 2; pos < testL; pos++ {
		m := pos % len(fill)
		i := rng.Intn(len(fill[m]))
		r := fill[m][i]
		for _, m := range combinedMappings {
			if pSrc := runePos(r, m.src); pSrc != -1 {
				set(pos, r, m.dst[pSrc], false)
				continue outer
			} else if pDst := runePos(r, m.dst); pDst != -1 {
				set(pos, r, r, true)
				continue outer
			}
		}
		set(pos, r, r, false)
	}

	rng.Shuffle(testL, func(i, j int) {
		rIn[i], rIn[j] = rIn[j], rIn[i]
		rOut[i], rOut[j] = rOut[j], rOut[i]
		quoteOut[i], quoteOut[j] = quoteOut[j], quoteOut[i]
	})

	var bOut strings.Builder
	bOut.Grow(testL)
	for i, r := range rOut {
		if quoteOut[i] {
			bOut.WriteRune(encoder.QuoteRune)
		}
		bOut.WriteRune(r)
	}
	return string(rIn), bOut.String()
}

func buildEdgeTestString(edge int, orig, replace rune, testMappings []mapping, fill ...[]rune) (out []stringPair) {
	testL := 30
	rIn := make([]rune, testL)
	rOut := make([]rune, testL)
	quoteOut := make([]bool, testL)

	set := func(i int, in, out rune, quote bool) {
		rIn[i] = in
		rOut[i] = out
		quoteOut[i] = quote
	}

outer:
	for pos := 0; pos < testL; pos++ {
		m := pos % len(fill)
		i := rng.Intn(len(fill[m]))
		r := fill[m][i]
		for _, m := range testMappings {
			if pSrc := runePos(r, m.src); pSrc != -1 {
				set(pos, r, m.dst[pSrc], false)
				continue outer
			} else if pDst := runePos(r, m.dst); pDst != -1 {
				set(pos, r, r, true)
				continue outer
			}
		}
		set(pos, r, r, false)
	}

	rng.Shuffle(testL, func(i, j int) {
		rIn[i], rIn[j] = rIn[j], rIn[i]
		rOut[i], rOut[j] = rOut[j], rOut[i]
		quoteOut[i], quoteOut[j] = quoteOut[j], quoteOut[i]
	})
	set(10, orig, orig, false)

	out = append(out, stringPair{string(rIn), quotedToString(rOut, quoteOut)})
	for _, i := range []int{0, 1, testL - 2, testL - 1} {
		for _, j := range []int{1, testL - 2, testL - 1} {
			if j < i {
				continue
			}
			rIn := append([]rune{}, rIn...)
			rOut := append([]rune{}, rOut...)
			quoteOut := append([]bool{}, quoteOut...)

			for _, in := range []rune{orig, replace} {
				expect, quote := in, false
				if i == 0 && edge == edgeLeft ||
					i == testL-1 && edge == edgeRight {
					expect, quote = replace, in == replace
				}
				rIn[i], rOut[i], quoteOut[i] = in, expect, quote

				if i != j {
					for _, in := range []rune{orig, replace} {
						expect, quote = in, false
						if j == testL-1 && edge == edgeRight {
							expect, quote = replace, in == replace
						}
						rIn[j], rOut[j], quoteOut[j] = in, expect, quote
					}
				}
				out = append(out, stringPair{string(rIn), quotedToString(rOut, quoteOut)})
			}
		}
	}
	return
}

func runePos(r rune, s []rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return -1
}

// quotedToString returns a string for the chars slice where a encoder.QuoteRune is
// inserted before a char[i] when quoted[i] is true.
func quotedToString(chars []rune, quoted []bool) string {
	var out strings.Builder
	out.Grow(len(chars))
	for i, r := range chars {
		if quoted[i] {
			out.WriteRune(encoder.QuoteRune)
		}
		out.WriteRune(r)
	}
	return out.String()
}
