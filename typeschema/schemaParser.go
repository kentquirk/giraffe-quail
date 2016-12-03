package typeschema

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var TR *TypeRegistry
var VR *ValueRegistry

var g = &grammar{
	rules: []*rule{
		{
			name: "Document",
			pos:  position{line: 24, col: 1, offset: 837},
			expr: &seqExpr{
				pos: position{line: 25, col: 7, offset: 855},
				exprs: []interface{}{
					&oneOrMoreExpr{
						pos: position{line: 25, col: 7, offset: 855},
						expr: &ruleRefExpr{
							pos:  position{line: 25, col: 7, offset: 855},
							name: "Definition",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 19, offset: 867},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 21, offset: 869},
						name: "EOF",
					},
				},
			},
		},
		{
			name: "Definition",
			pos:  position{line: 28, col: 1, offset: 880},
			expr: &choiceExpr{
				pos: position{line: 29, col: 7, offset: 900},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 29, col: 7, offset: 900},
						name: "ObjTypeDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 30, col: 7, offset: 924},
						name: "EnumDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 31, col: 7, offset: 945},
						name: "InterfaceDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 32, col: 7, offset: 971},
						name: "UnionDefinition",
					},
				},
			},
		},
		{
			name: "ObjTypeDefinition",
			pos:  position{line: 35, col: 1, offset: 994},
			expr: &seqExpr{
				pos: position{line: 35, col: 22, offset: 1015},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 35, col: 22, offset: 1015},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 35, col: 24, offset: 1017},
						val:        "type",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 35, col: 31, offset: 1024},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 35, col: 33, offset: 1026},
						name: "Name",
					},
					&zeroOrOneExpr{
						pos: position{line: 35, col: 38, offset: 1031},
						expr: &ruleRefExpr{
							pos:  position{line: 35, col: 38, offset: 1031},
							name: "Implements",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 35, col: 50, offset: 1043},
						name: "FieldSet",
					},
				},
			},
		},
		{
			name: "EnumDefinition",
			pos:  position{line: 37, col: 1, offset: 1053},
			expr: &actionExpr{
				pos: position{line: 37, col: 19, offset: 1071},
				run: (*parser).callonEnumDefinition1,
				expr: &seqExpr{
					pos: position{line: 37, col: 19, offset: 1071},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 37, col: 19, offset: 1071},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 37, col: 21, offset: 1073},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 37, col: 28, offset: 1080},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 37, col: 30, offset: 1082},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 37, col: 32, offset: 1084},
								name: "Name",
							},
						},
						&labeledExpr{
							pos:   position{line: 37, col: 37, offset: 1089},
							label: "vals",
							expr: &ruleRefExpr{
								pos:  position{line: 37, col: 42, offset: 1094},
								name: "EnumValueSet",
							},
						},
					},
				},
			},
		},
		{
			name: "InterfaceDefinition",
			pos:  position{line: 54, col: 1, offset: 1534},
			expr: &seqExpr{
				pos: position{line: 54, col: 24, offset: 1557},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 54, col: 24, offset: 1557},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 54, col: 26, offset: 1559},
						val:        "interface",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 54, col: 38, offset: 1571},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 54, col: 40, offset: 1573},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 54, col: 45, offset: 1578},
						name: "FieldSet",
					},
				},
			},
		},
		{
			name: "UnionDefinition",
			pos:  position{line: 56, col: 1, offset: 1588},
			expr: &seqExpr{
				pos: position{line: 56, col: 20, offset: 1607},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 56, col: 20, offset: 1607},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 56, col: 22, offset: 1609},
						val:        "union",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 56, col: 30, offset: 1617},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 56, col: 32, offset: 1619},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 56, col: 37, offset: 1624},
						name: "UnionSet",
					},
				},
			},
		},
		{
			name: "Implements",
			pos:  position{line: 59, col: 1, offset: 1689},
			expr: &seqExpr{
				pos: position{line: 59, col: 15, offset: 1703},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 59, col: 15, offset: 1703},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 59, col: 17, offset: 1705},
						val:        "implements",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 59, col: 30, offset: 1718},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 59, col: 32, offset: 1720},
						name: "Name",
					},
				},
			},
		},
		{
			name: "EnumValueSet",
			pos:  position{line: 61, col: 1, offset: 1726},
			expr: &actionExpr{
				pos: position{line: 61, col: 17, offset: 1742},
				run: (*parser).callonEnumValueSet1,
				expr: &seqExpr{
					pos: position{line: 61, col: 17, offset: 1742},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 61, col: 17, offset: 1742},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 19, offset: 1744},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 61, col: 23, offset: 1748},
							label: "ev",
							expr: &ruleRefExpr{
								pos:  position{line: 61, col: 26, offset: 1751},
								name: "EnumValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 37, offset: 1762},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 39, offset: 1764},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "EnumValues",
			pos:  position{line: 63, col: 1, offset: 1788},
			expr: &actionExpr{
				pos: position{line: 63, col: 15, offset: 1802},
				run: (*parser).callonEnumValues1,
				expr: &labeledExpr{
					pos:   position{line: 63, col: 15, offset: 1802},
					label: "evs",
					expr: &oneOrMoreExpr{
						pos: position{line: 63, col: 19, offset: 1806},
						expr: &ruleRefExpr{
							pos:  position{line: 63, col: 19, offset: 1806},
							name: "EnumValue",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 71, col: 1, offset: 1951},
			expr: &actionExpr{
				pos: position{line: 71, col: 14, offset: 1964},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 71, col: 14, offset: 1964},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 71, col: 14, offset: 1964},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 71, col: 16, offset: 1966},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 71, col: 18, offset: 1968},
								name: "Name",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionSet",
			pos:  position{line: 73, col: 1, offset: 2001},
			expr: &seqExpr{
				pos: position{line: 73, col: 13, offset: 2013},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 73, col: 13, offset: 2013},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 73, col: 15, offset: 2015},
						val:        "=",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 19, offset: 2019},
						name: "UnionNames",
					},
				},
			},
		},
		{
			name: "UnionNames",
			pos:  position{line: 75, col: 1, offset: 2031},
			expr: &seqExpr{
				pos: position{line: 75, col: 15, offset: 2045},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 75, col: 15, offset: 2045},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 17, offset: 2047},
						name: "Name",
					},
					&zeroOrMoreExpr{
						pos: position{line: 75, col: 22, offset: 2052},
						expr: &seqExpr{
							pos: position{line: 75, col: 23, offset: 2053},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 75, col: 23, offset: 2053},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 75, col: 25, offset: 2055},
									val:        "|",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 75, col: 29, offset: 2059},
									name: "_",
								},
								&ruleRefExpr{
									pos:  position{line: 75, col: 31, offset: 2061},
									name: "Name",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FieldSet",
			pos:  position{line: 77, col: 1, offset: 2069},
			expr: &seqExpr{
				pos: position{line: 77, col: 14, offset: 2082},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 77, col: 14, offset: 2082},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 77, col: 16, offset: 2084},
						val:        "{",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 77, col: 20, offset: 2088},
						expr: &ruleRefExpr{
							pos:  position{line: 77, col: 20, offset: 2088},
							name: "Field",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 77, col: 27, offset: 2095},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 77, col: 29, offset: 2097},
						val:        "}",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 79, col: 1, offset: 2102},
			expr: &seqExpr{
				pos: position{line: 79, col: 10, offset: 2111},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 79, col: 10, offset: 2111},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 12, offset: 2113},
						name: "Name",
					},
					&zeroOrOneExpr{
						pos: position{line: 79, col: 17, offset: 2118},
						expr: &ruleRefExpr{
							pos:  position{line: 79, col: 17, offset: 2118},
							name: "Arguments",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 28, offset: 2129},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 79, col: 30, offset: 2131},
						val:        ":",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 34, offset: 2135},
						name: "Type",
					},
				},
			},
		},
		{
			name: "Arguments",
			pos:  position{line: 81, col: 1, offset: 2141},
			expr: &seqExpr{
				pos: position{line: 81, col: 14, offset: 2154},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 81, col: 14, offset: 2154},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 81, col: 16, offset: 2156},
						val:        "(",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 81, col: 20, offset: 2160},
						expr: &ruleRefExpr{
							pos:  position{line: 81, col: 20, offset: 2160},
							name: "Argument",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 81, col: 30, offset: 2170},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 81, col: 32, offset: 2172},
						val:        ")",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Argument",
			pos:  position{line: 83, col: 1, offset: 2177},
			expr: &seqExpr{
				pos: position{line: 83, col: 13, offset: 2189},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 83, col: 13, offset: 2189},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 83, col: 15, offset: 2191},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 83, col: 20, offset: 2196},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 83, col: 22, offset: 2198},
						val:        ":",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 83, col: 26, offset: 2202},
						name: "Type",
					},
				},
			},
		},
		{
			name: "Type",
			pos:  position{line: 85, col: 1, offset: 2208},
			expr: &choiceExpr{
				pos: position{line: 86, col: 7, offset: 2222},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 86, col: 7, offset: 2222},
						name: "NonNullType",
					},
					&ruleRefExpr{
						pos:  position{line: 87, col: 7, offset: 2240},
						name: "ListType",
					},
					&ruleRefExpr{
						pos:  position{line: 88, col: 7, offset: 2255},
						name: "NamedType",
					},
				},
			},
		},
		{
			name: "NamedType",
			pos:  position{line: 91, col: 1, offset: 2272},
			expr: &actionExpr{
				pos: position{line: 91, col: 14, offset: 2285},
				run: (*parser).callonNamedType1,
				expr: &seqExpr{
					pos: position{line: 91, col: 14, offset: 2285},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 91, col: 14, offset: 2285},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 91, col: 16, offset: 2287},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 91, col: 18, offset: 2289},
								name: "Name",
							},
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 93, col: 1, offset: 2330},
			expr: &actionExpr{
				pos: position{line: 93, col: 13, offset: 2342},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 93, col: 13, offset: 2342},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 93, col: 13, offset: 2342},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 93, col: 15, offset: 2344},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 93, col: 19, offset: 2348},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 93, col: 21, offset: 2350},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 93, col: 23, offset: 2352},
								name: "Type",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 93, col: 28, offset: 2357},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 93, col: 30, offset: 2359},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "NonNullType",
			pos:  position{line: 95, col: 1, offset: 2422},
			expr: &choiceExpr{
				pos: position{line: 96, col: 7, offset: 2443},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 96, col: 7, offset: 2443},
						run: (*parser).callonNonNullType2,
						expr: &seqExpr{
							pos: position{line: 96, col: 7, offset: 2443},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 96, col: 7, offset: 2443},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 96, col: 9, offset: 2445},
										name: "ListType",
									},
								},
								&litMatcher{
									pos:        position{line: 96, col: 18, offset: 2454},
									val:        "!",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 97, col: 7, offset: 2530},
						run: (*parser).callonNonNullType7,
						expr: &seqExpr{
							pos: position{line: 97, col: 7, offset: 2530},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 97, col: 7, offset: 2530},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 97, col: 9, offset: 2532},
										name: "NamedType",
									},
								},
								&litMatcher{
									pos:        position{line: 97, col: 19, offset: 2542},
									val:        "!",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 102, col: 1, offset: 2637},
			expr: &zeroOrMoreExpr{
				pos: position{line: 102, col: 6, offset: 2642},
				expr: &ruleRefExpr{
					pos:  position{line: 102, col: 6, offset: 2642},
					name: "I",
				},
			},
		},
		{
			name:        "I",
			displayName: "\"Ignored\"",
			pos:         position{line: 103, col: 1, offset: 2645},
			expr: &choiceExpr{
				pos: position{line: 104, col: 6, offset: 2665},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 104, col: 6, offset: 2665},
						name: "UnicodeBOM",
					},
					&ruleRefExpr{
						pos:  position{line: 105, col: 6, offset: 2681},
						name: "WhiteSpace",
					},
					&ruleRefExpr{
						pos:  position{line: 106, col: 6, offset: 2697},
						name: "LineTerminator",
					},
					&ruleRefExpr{
						pos:  position{line: 107, col: 6, offset: 2717},
						name: "Comment",
					},
					&ruleRefExpr{
						pos:  position{line: 108, col: 6, offset: 2730},
						name: "Comma",
					},
				},
			},
		},
		{
			name:        "UnicodeBOM",
			displayName: "\"ByteOrderMark\"",
			pos:         position{line: 113, col: 1, offset: 2777},
			expr: &litMatcher{
				pos:        position{line: 114, col: 5, offset: 2811},
				val:        "\ufeff",
				ignoreCase: false,
			},
		},
		{
			name:        "LineTerminator",
			displayName: "\"EOL\"",
			pos:         position{line: 116, col: 1, offset: 2821},
			expr: &choiceExpr{
				pos: position{line: 117, col: 7, offset: 2851},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 117, col: 7, offset: 2851},
						val:        "\n",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 118, col: 7, offset: 2862},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 118, col: 7, offset: 2862},
								val:        "\r",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 118, col: 12, offset: 2867},
								expr: &litMatcher{
									pos:        position{line: 118, col: 13, offset: 2868},
									val:        "\n",
									ignoreCase: false,
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 119, col: 7, offset: 2879},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 119, col: 7, offset: 2879},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 119, col: 12, offset: 2884},
								val:        "\n",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "WhiteSpace",
			pos:  position{line: 122, col: 1, offset: 2896},
			expr: &charClassMatcher{
				pos:        position{line: 122, col: 15, offset: 2910},
				val:        "[ \\t]",
				chars:      []rune{' ', '\t'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 123, col: 1, offset: 2916},
			expr: &seqExpr{
				pos: position{line: 123, col: 12, offset: 2927},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 123, col: 12, offset: 2927},
						val:        "#",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 123, col: 16, offset: 2931},
						expr: &anyMatcher{
							line: 123, col: 16, offset: 2931,
						},
					},
					&notExpr{
						pos: position{line: 123, col: 19, offset: 2934},
						expr: &ruleRefExpr{
							pos:  position{line: 123, col: 20, offset: 2935},
							name: "LineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Comma",
			pos:  position{line: 124, col: 1, offset: 2950},
			expr: &litMatcher{
				pos:        position{line: 124, col: 10, offset: 2959},
				val:        ",",
				ignoreCase: false,
			},
		},
		{
			name: "Name",
			pos:  position{line: 125, col: 1, offset: 2963},
			expr: &actionExpr{
				pos: position{line: 125, col: 9, offset: 2971},
				run: (*parser).callonName1,
				expr: &seqExpr{
					pos: position{line: 125, col: 9, offset: 2971},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 125, col: 9, offset: 2971},
							val:        "[_A-Za-z]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 125, col: 18, offset: 2980},
							expr: &charClassMatcher{
								pos:        position{line: 125, col: 18, offset: 2980},
								val:        "[_0-9A-Za-z]",
								chars:      []rune{'_'},
								ranges:     []rune{'0', '9', 'A', 'Z', 'a', 'z'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 127, col: 1, offset: 3026},
			expr: &notExpr{
				pos: position{line: 127, col: 7, offset: 3032},
				expr: &anyMatcher{
					line: 127, col: 8, offset: 3033,
				},
			},
		},
	},
}

func (c *current) onEnumDefinition1(n, vals interface{}) (interface{}, error) {
	enumname := n.(string)
	enumtype, err := TR.Register(enumname, Enum)
	if err != nil {
		return enumtype, err
	}
	for _, v := range vals.([]string) {
		enumvalue := v
		_, err := VR.Register(enumvalue, Value{T: enumtype, V: enumvalue})
		if err != nil {
			return enumtype, err
		}
	}
	return enumtype, nil

}

func (p *parser) callonEnumDefinition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumDefinition1(stack["n"], stack["vals"])
}

func (c *current) onEnumValueSet1(ev interface{}) (interface{}, error) {
	return ev, nil
}

func (p *parser) callonEnumValueSet1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValueSet1(stack["ev"])
}

func (c *current) onEnumValues1(evs interface{}) (interface{}, error) {
	r := make([]string, 0)
	for _, n := range evs.([]interface{}) {
		r = append(r, n.(string))
	}
	return r, nil
}

func (p *parser) callonEnumValues1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValues1(stack["evs"])
}

func (c *current) onEnumValue1(n interface{}) (interface{}, error) {
	return n.(string), nil
}

func (p *parser) callonEnumValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValue1(stack["n"])
}

func (c *current) onNamedType1(n interface{}) (interface{}, error) {
	return TR.MaybeGet(n.(string))
}

func (p *parser) callonNamedType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamedType1(stack["n"])
}

func (c *current) onListType1(t interface{}) (interface{}, error) {
	return TR.MaybeGet(TypeNameFor(List, t.(Type).Key()))
}

func (p *parser) callonListType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListType1(stack["t"])
}

func (c *current) onNonNullType2(t interface{}) (interface{}, error) {
	return TR.MaybeGet(TypeNameFor(NonNullable, t.(Type).Key()))
}

func (p *parser) callonNonNullType2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonNullType2(stack["t"])
}

func (c *current) onNonNullType7(t interface{}) (interface{}, error) {
	return TR.MaybeGet(TypeNameFor(NonNullable, t.(Type).Key()))
}

func (p *parser) callonNonNullType7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonNullType7(stack["t"])
}

func (c *current) onName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onName1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
