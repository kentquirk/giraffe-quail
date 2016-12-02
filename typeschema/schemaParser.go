package parser

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

var typereg TypeRegistry

func init() {
	typereg = CreateTypeRegistry()
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Document",
			pos:  position{line: 26, col: 1, offset: 874},
			expr: &seqExpr{
				pos: position{line: 27, col: 7, offset: 892},
				exprs: []interface{}{
					&oneOrMoreExpr{
						pos: position{line: 27, col: 7, offset: 892},
						expr: &ruleRefExpr{
							pos:  position{line: 27, col: 7, offset: 892},
							name: "Definition",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 27, col: 19, offset: 904},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 27, col: 21, offset: 906},
						name: "EOF",
					},
				},
			},
		},
		{
			name: "Definition",
			pos:  position{line: 30, col: 1, offset: 917},
			expr: &choiceExpr{
				pos: position{line: 31, col: 7, offset: 937},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 31, col: 7, offset: 937},
						name: "ObjTypeDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 32, col: 7, offset: 961},
						name: "EnumDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 7, offset: 982},
						name: "InterfaceDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 34, col: 7, offset: 1008},
						name: "UnionDefinition",
					},
				},
			},
		},
		{
			name: "ObjTypeDefinition",
			pos:  position{line: 37, col: 1, offset: 1031},
			expr: &seqExpr{
				pos: position{line: 37, col: 22, offset: 1052},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 37, col: 22, offset: 1052},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 37, col: 24, offset: 1054},
						val:        "type",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 37, col: 31, offset: 1061},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 37, col: 33, offset: 1063},
						name: "Name",
					},
					&zeroOrOneExpr{
						pos: position{line: 37, col: 38, offset: 1068},
						expr: &ruleRefExpr{
							pos:  position{line: 37, col: 38, offset: 1068},
							name: "Implements",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 37, col: 50, offset: 1080},
						name: "FieldSet",
					},
				},
			},
		},
		{
			name: "EnumDefinition",
			pos:  position{line: 39, col: 1, offset: 1090},
			expr: &seqExpr{
				pos: position{line: 39, col: 19, offset: 1108},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 39, col: 19, offset: 1108},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 39, col: 21, offset: 1110},
						val:        "enum",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 39, col: 28, offset: 1117},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 39, col: 30, offset: 1119},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 39, col: 35, offset: 1124},
						name: "EnumValueSet",
					},
				},
			},
		},
		{
			name: "InterfaceDefinition",
			pos:  position{line: 41, col: 1, offset: 1138},
			expr: &seqExpr{
				pos: position{line: 41, col: 24, offset: 1161},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 41, col: 24, offset: 1161},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 41, col: 26, offset: 1163},
						val:        "interface",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 41, col: 38, offset: 1175},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 41, col: 40, offset: 1177},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 41, col: 45, offset: 1182},
						name: "FieldSet",
					},
				},
			},
		},
		{
			name: "UnionDefinition",
			pos:  position{line: 43, col: 1, offset: 1192},
			expr: &seqExpr{
				pos: position{line: 43, col: 20, offset: 1211},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 43, col: 20, offset: 1211},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 43, col: 22, offset: 1213},
						val:        "union",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 43, col: 30, offset: 1221},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 43, col: 32, offset: 1223},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 43, col: 37, offset: 1228},
						name: "UnionSet",
					},
				},
			},
		},
		{
			name: "Implements",
			pos:  position{line: 46, col: 1, offset: 1293},
			expr: &seqExpr{
				pos: position{line: 46, col: 15, offset: 1307},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 46, col: 15, offset: 1307},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 46, col: 17, offset: 1309},
						val:        "implements",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 46, col: 30, offset: 1322},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 46, col: 32, offset: 1324},
						name: "Name",
					},
				},
			},
		},
		{
			name: "EnumValueSet",
			pos:  position{line: 48, col: 1, offset: 1330},
			expr: &seqExpr{
				pos: position{line: 48, col: 17, offset: 1346},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 48, col: 17, offset: 1346},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 48, col: 19, offset: 1348},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 48, col: 23, offset: 1352},
						name: "EnumValues",
					},
					&ruleRefExpr{
						pos:  position{line: 48, col: 34, offset: 1363},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 48, col: 36, offset: 1365},
						val:        "}",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EnumValues",
			pos:  position{line: 50, col: 1, offset: 1370},
			expr: &choiceExpr{
				pos: position{line: 51, col: 7, offset: 1390},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 51, col: 7, offset: 1390},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 51, col: 7, offset: 1390},
								name: "_",
							},
							&ruleRefExpr{
								pos:  position{line: 51, col: 9, offset: 1392},
								name: "Name",
							},
						},
					},
					&seqExpr{
						pos: position{line: 52, col: 7, offset: 1403},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 52, col: 7, offset: 1403},
								name: "_",
							},
							&ruleRefExpr{
								pos:  position{line: 52, col: 9, offset: 1405},
								name: "Name",
							},
							&ruleRefExpr{
								pos:  position{line: 52, col: 14, offset: 1410},
								name: "_",
							},
							&litMatcher{
								pos:        position{line: 52, col: 16, offset: 1412},
								val:        ",",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 52, col: 20, offset: 1416},
								name: "_",
							},
							&ruleRefExpr{
								pos:  position{line: 52, col: 22, offset: 1418},
								name: "EnumValues",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionSet",
			pos:  position{line: 55, col: 1, offset: 1436},
			expr: &seqExpr{
				pos: position{line: 55, col: 13, offset: 1448},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 55, col: 13, offset: 1448},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 55, col: 15, offset: 1450},
						val:        "=",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 55, col: 19, offset: 1454},
						name: "UnionNames",
					},
				},
			},
		},
		{
			name: "UnionNames",
			pos:  position{line: 57, col: 1, offset: 1466},
			expr: &choiceExpr{
				pos: position{line: 58, col: 7, offset: 1486},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 58, col: 7, offset: 1486},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 58, col: 7, offset: 1486},
								name: "_",
							},
							&ruleRefExpr{
								pos:  position{line: 58, col: 9, offset: 1488},
								name: "Name",
							},
						},
					},
					&seqExpr{
						pos: position{line: 59, col: 7, offset: 1499},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 59, col: 7, offset: 1499},
								name: "_",
							},
							&ruleRefExpr{
								pos:  position{line: 59, col: 9, offset: 1501},
								name: "Name",
							},
							&litMatcher{
								pos:        position{line: 59, col: 14, offset: 1506},
								val:        "|",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 59, col: 18, offset: 1510},
								name: "UnionNames",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldSet",
			pos:  position{line: 62, col: 1, offset: 1528},
			expr: &seqExpr{
				pos: position{line: 62, col: 14, offset: 1541},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 62, col: 14, offset: 1541},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 62, col: 16, offset: 1543},
						val:        "{",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 62, col: 20, offset: 1547},
						expr: &ruleRefExpr{
							pos:  position{line: 62, col: 20, offset: 1547},
							name: "Field",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 62, col: 27, offset: 1554},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 62, col: 29, offset: 1556},
						val:        "}",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 64, col: 1, offset: 1561},
			expr: &seqExpr{
				pos: position{line: 64, col: 10, offset: 1570},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 64, col: 10, offset: 1570},
						name: "Name",
					},
					&zeroOrOneExpr{
						pos: position{line: 64, col: 15, offset: 1575},
						expr: &ruleRefExpr{
							pos:  position{line: 64, col: 15, offset: 1575},
							name: "Arguments",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 64, col: 26, offset: 1586},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 64, col: 28, offset: 1588},
						val:        ":",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 64, col: 32, offset: 1592},
						name: "Type",
					},
				},
			},
		},
		{
			name: "Arguments",
			pos:  position{line: 66, col: 1, offset: 1598},
			expr: &seqExpr{
				pos: position{line: 66, col: 14, offset: 1611},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 66, col: 14, offset: 1611},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 66, col: 16, offset: 1613},
						val:        "(",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 66, col: 20, offset: 1617},
						expr: &ruleRefExpr{
							pos:  position{line: 66, col: 20, offset: 1617},
							name: "Argument",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 66, col: 30, offset: 1627},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 66, col: 32, offset: 1629},
						val:        ")",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Argument",
			pos:  position{line: 68, col: 1, offset: 1634},
			expr: &seqExpr{
				pos: position{line: 68, col: 13, offset: 1646},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 68, col: 13, offset: 1646},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 15, offset: 1648},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 20, offset: 1653},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 68, col: 22, offset: 1655},
						val:        ":",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 26, offset: 1659},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 28, offset: 1661},
						name: "Type",
					},
				},
			},
		},
		{
			name: "Type",
			pos:  position{line: 70, col: 1, offset: 1667},
			expr: &choiceExpr{
				pos: position{line: 71, col: 7, offset: 1681},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 71, col: 7, offset: 1681},
						name: "NonNullType",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 7, offset: 1699},
						name: "ListType",
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 7, offset: 1714},
						name: "ScalarType",
					},
					&ruleRefExpr{
						pos:  position{line: 74, col: 7, offset: 1731},
						name: "NamedType",
					},
				},
			},
		},
		{
			name: "ScalarType",
			pos:  position{line: 78, col: 1, offset: 1828},
			expr: &choiceExpr{
				pos: position{line: 79, col: 7, offset: 1848},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 79, col: 7, offset: 1848},
						val:        "Int",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 80, col: 7, offset: 1860},
						val:        "Float",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 81, col: 7, offset: 1874},
						val:        "String",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 82, col: 7, offset: 1889},
						val:        "Boolean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 83, col: 7, offset: 1905},
						val:        "ID",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "NamedType",
			pos:  position{line: 86, col: 1, offset: 1917},
			expr: &seqExpr{
				pos: position{line: 86, col: 14, offset: 1930},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 86, col: 14, offset: 1930},
						name: "_",
					},
					&labeledExpr{
						pos:   position{line: 86, col: 16, offset: 1932},
						label: "n",
						expr: &ruleRefExpr{
							pos:  position{line: 86, col: 18, offset: 1934},
							name: "Name",
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 88, col: 1, offset: 1940},
			expr: &seqExpr{
				pos: position{line: 88, col: 13, offset: 1952},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 88, col: 13, offset: 1952},
						val:        "[",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 88, col: 17, offset: 1956},
						name: "_",
					},
					&labeledExpr{
						pos:   position{line: 88, col: 19, offset: 1958},
						label: "t",
						expr: &ruleRefExpr{
							pos:  position{line: 88, col: 21, offset: 1960},
							name: "Type",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 88, col: 26, offset: 1965},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 88, col: 28, offset: 1967},
						val:        "]",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "NonNullType",
			pos:  position{line: 90, col: 1, offset: 1972},
			expr: &choiceExpr{
				pos: position{line: 91, col: 7, offset: 1993},
				alternatives: []interface{}{
					&labeledExpr{
						pos:   position{line: 91, col: 7, offset: 1993},
						label: "t",
						expr: &ruleRefExpr{
							pos:  position{line: 91, col: 9, offset: 1995},
							name: "ListType",
						},
					},
					&seqExpr{
						pos: position{line: 92, col: 7, offset: 2010},
						exprs: []interface{}{
							&labeledExpr{
								pos:   position{line: 92, col: 7, offset: 2010},
								label: "t",
								expr: &ruleRefExpr{
									pos:  position{line: 92, col: 9, offset: 2012},
									name: "ScalarType",
								},
							},
							&litMatcher{
								pos:        position{line: 92, col: 20, offset: 2023},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 93, col: 7, offset: 2033},
						exprs: []interface{}{
							&labeledExpr{
								pos:   position{line: 93, col: 7, offset: 2033},
								label: "t",
								expr: &ruleRefExpr{
									pos:  position{line: 93, col: 9, offset: 2035},
									name: "NamedType",
								},
							},
							&litMatcher{
								pos:        position{line: 93, col: 19, offset: 2045},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 98, col: 1, offset: 2075},
			expr: &zeroOrMoreExpr{
				pos: position{line: 98, col: 6, offset: 2080},
				expr: &ruleRefExpr{
					pos:  position{line: 98, col: 6, offset: 2080},
					name: "I",
				},
			},
		},
		{
			name:        "I",
			displayName: "\"Ignored\"",
			pos:         position{line: 99, col: 1, offset: 2083},
			expr: &choiceExpr{
				pos: position{line: 100, col: 6, offset: 2103},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 100, col: 6, offset: 2103},
						name: "UnicodeBOM",
					},
					&ruleRefExpr{
						pos:  position{line: 101, col: 6, offset: 2119},
						name: "WhiteSpace",
					},
					&ruleRefExpr{
						pos:  position{line: 102, col: 6, offset: 2135},
						name: "LineTerminator",
					},
					&ruleRefExpr{
						pos:  position{line: 103, col: 6, offset: 2155},
						name: "Comment",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 6, offset: 2168},
						name: "Comma",
					},
				},
			},
		},
		{
			name:        "UnicodeBOM",
			displayName: "\"ByteOrderMark\"",
			pos:         position{line: 109, col: 1, offset: 2215},
			expr: &litMatcher{
				pos:        position{line: 110, col: 5, offset: 2249},
				val:        "\ufeff",
				ignoreCase: false,
			},
		},
		{
			name:        "LineTerminator",
			displayName: "\"EOL\"",
			pos:         position{line: 112, col: 1, offset: 2259},
			expr: &choiceExpr{
				pos: position{line: 113, col: 7, offset: 2289},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 113, col: 7, offset: 2289},
						val:        "\n",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 114, col: 7, offset: 2300},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 114, col: 7, offset: 2300},
								val:        "\r",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 114, col: 12, offset: 2305},
								expr: &litMatcher{
									pos:        position{line: 114, col: 13, offset: 2306},
									val:        "\n",
									ignoreCase: false,
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 115, col: 7, offset: 2317},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 115, col: 7, offset: 2317},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 115, col: 12, offset: 2322},
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
			pos:  position{line: 118, col: 1, offset: 2334},
			expr: &charClassMatcher{
				pos:        position{line: 118, col: 15, offset: 2348},
				val:        "[ \\t]",
				chars:      []rune{' ', '\t'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 119, col: 1, offset: 2354},
			expr: &seqExpr{
				pos: position{line: 119, col: 12, offset: 2365},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 119, col: 12, offset: 2365},
						val:        "#",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 119, col: 16, offset: 2369},
						expr: &anyMatcher{
							line: 119, col: 16, offset: 2369,
						},
					},
					&notExpr{
						pos: position{line: 119, col: 19, offset: 2372},
						expr: &ruleRefExpr{
							pos:  position{line: 119, col: 20, offset: 2373},
							name: "LineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Comma",
			pos:  position{line: 120, col: 1, offset: 2388},
			expr: &litMatcher{
				pos:        position{line: 120, col: 10, offset: 2397},
				val:        ",",
				ignoreCase: false,
			},
		},
		{
			name: "Name",
			pos:  position{line: 121, col: 1, offset: 2401},
			expr: &actionExpr{
				pos: position{line: 121, col: 9, offset: 2409},
				run: (*parser).callonName1,
				expr: &seqExpr{
					pos: position{line: 121, col: 9, offset: 2409},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 121, col: 9, offset: 2409},
							val:        "[_A-Za-z]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 121, col: 18, offset: 2418},
							expr: &charClassMatcher{
								pos:        position{line: 121, col: 18, offset: 2418},
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
			pos:  position{line: 123, col: 1, offset: 2464},
			expr: &notExpr{
				pos: position{line: 123, col: 7, offset: 2470},
				expr: &anyMatcher{
					line: 123, col: 8, offset: 2471,
				},
			},
		},
	},
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
