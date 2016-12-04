package main

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
			pos:  position{line: 26, col: 1, offset: 849},
			expr: &seqExpr{
				pos: position{line: 27, col: 7, offset: 867},
				exprs: []interface{}{
					&oneOrMoreExpr{
						pos: position{line: 27, col: 7, offset: 867},
						expr: &ruleRefExpr{
							pos:  position{line: 27, col: 7, offset: 867},
							name: "Definition",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 27, col: 19, offset: 879},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 27, col: 21, offset: 881},
						name: "EOF",
					},
				},
			},
		},
		{
			name: "Definition",
			pos:  position{line: 30, col: 1, offset: 892},
			expr: &choiceExpr{
				pos: position{line: 31, col: 7, offset: 912},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 31, col: 7, offset: 912},
						name: "FragmentDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 32, col: 7, offset: 937},
						name: "OperationDefinition",
					},
				},
			},
		},
		{
			name: "OperationDefinition",
			pos:  position{line: 35, col: 1, offset: 964},
			expr: &choiceExpr{
				pos: position{line: 36, col: 7, offset: 993},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 36, col: 7, offset: 993},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 36, col: 7, offset: 993},
								name: "OperationType",
							},
							&ruleRefExpr{
								pos:  position{line: 36, col: 21, offset: 1007},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 36, col: 23, offset: 1009},
								expr: &ruleRefExpr{
									pos:  position{line: 36, col: 23, offset: 1009},
									name: "Name",
								},
							},
							&zeroOrOneExpr{
								pos: position{line: 36, col: 29, offset: 1015},
								expr: &ruleRefExpr{
									pos:  position{line: 36, col: 29, offset: 1015},
									name: "VariableDefinitions",
								},
							},
							&zeroOrOneExpr{
								pos: position{line: 36, col: 50, offset: 1036},
								expr: &ruleRefExpr{
									pos:  position{line: 36, col: 50, offset: 1036},
									name: "Directives",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 36, col: 62, offset: 1048},
								name: "SelectionSet",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 37, col: 7, offset: 1067},
						name: "SelectionSet",
					},
				},
			},
		},
		{
			name: "OperationType",
			pos:  position{line: 40, col: 1, offset: 1087},
			expr: &choiceExpr{
				pos: position{line: 41, col: 7, offset: 1110},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 41, col: 7, offset: 1110},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 41, col: 7, offset: 1110},
								name: "_",
							},
							&litMatcher{
								pos:        position{line: 41, col: 9, offset: 1112},
								val:        "query",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 42, col: 7, offset: 1126},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 42, col: 7, offset: 1126},
								name: "_",
							},
							&litMatcher{
								pos:        position{line: 42, col: 9, offset: 1128},
								val:        "mutation",
								ignoreCase: false,
							},
						},
					},
					&actionExpr{
						pos: position{line: 43, col: 7, offset: 1145},
						run: (*parser).callonOperationType8,
						expr: &seqExpr{
							pos: position{line: 43, col: 7, offset: 1145},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 43, col: 7, offset: 1145},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 43, col: 9, offset: 1147},
									label: "n",
									expr: &ruleRefExpr{
										pos:  position{line: 43, col: 11, offset: 1149},
										name: "Name",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SelectionSet",
			pos:  position{line: 46, col: 1, offset: 1224},
			expr: &seqExpr{
				pos: position{line: 46, col: 17, offset: 1240},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 46, col: 17, offset: 1240},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 46, col: 19, offset: 1242},
						val:        "{",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 46, col: 23, offset: 1246},
						expr: &ruleRefExpr{
							pos:  position{line: 46, col: 23, offset: 1246},
							name: "Selection",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 46, col: 34, offset: 1257},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 46, col: 36, offset: 1259},
						val:        "}",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Selection",
			pos:  position{line: 48, col: 1, offset: 1264},
			expr: &choiceExpr{
				pos: position{line: 49, col: 7, offset: 1283},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 49, col: 7, offset: 1283},
						name: "Field",
					},
					&ruleRefExpr{
						pos:  position{line: 50, col: 7, offset: 1295},
						name: "FragmentSpread",
					},
					&ruleRefExpr{
						pos:  position{line: 51, col: 7, offset: 1316},
						name: "InlineFragment",
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 54, col: 1, offset: 1338},
			expr: &seqExpr{
				pos: position{line: 54, col: 10, offset: 1347},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 54, col: 10, offset: 1347},
						expr: &ruleRefExpr{
							pos:  position{line: 54, col: 10, offset: 1347},
							name: "Alias",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 54, col: 17, offset: 1354},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 54, col: 19, offset: 1356},
						name: "Name",
					},
					&zeroOrOneExpr{
						pos: position{line: 54, col: 24, offset: 1361},
						expr: &ruleRefExpr{
							pos:  position{line: 54, col: 24, offset: 1361},
							name: "Arguments",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 54, col: 35, offset: 1372},
						expr: &ruleRefExpr{
							pos:  position{line: 54, col: 35, offset: 1372},
							name: "Directives",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 54, col: 47, offset: 1384},
						expr: &ruleRefExpr{
							pos:  position{line: 54, col: 47, offset: 1384},
							name: "SelectionSet",
						},
					},
				},
			},
		},
		{
			name: "Arguments",
			pos:  position{line: 56, col: 1, offset: 1399},
			expr: &seqExpr{
				pos: position{line: 56, col: 14, offset: 1412},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 56, col: 14, offset: 1412},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 56, col: 16, offset: 1414},
						val:        "(",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 56, col: 20, offset: 1418},
						expr: &ruleRefExpr{
							pos:  position{line: 56, col: 20, offset: 1418},
							name: "Argument",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 56, col: 30, offset: 1428},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 56, col: 32, offset: 1430},
						val:        ")",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Argument",
			pos:  position{line: 58, col: 1, offset: 1435},
			expr: &seqExpr{
				pos: position{line: 58, col: 13, offset: 1447},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 58, col: 13, offset: 1447},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 58, col: 15, offset: 1449},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 58, col: 20, offset: 1454},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 58, col: 22, offset: 1456},
						val:        ":",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 58, col: 26, offset: 1460},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 58, col: 28, offset: 1462},
						name: "Value",
					},
				},
			},
		},
		{
			name: "Alias",
			pos:  position{line: 60, col: 1, offset: 1469},
			expr: &seqExpr{
				pos: position{line: 60, col: 10, offset: 1478},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 60, col: 10, offset: 1478},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 60, col: 12, offset: 1480},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 60, col: 17, offset: 1485},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 60, col: 19, offset: 1487},
						val:        ":",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "FragmentSpread",
			pos:  position{line: 62, col: 1, offset: 1492},
			expr: &seqExpr{
				pos: position{line: 62, col: 19, offset: 1510},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 62, col: 19, offset: 1510},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 62, col: 21, offset: 1512},
						val:        "...",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 62, col: 27, offset: 1518},
						name: "FragmentName",
					},
					&zeroOrOneExpr{
						pos: position{line: 62, col: 40, offset: 1531},
						expr: &ruleRefExpr{
							pos:  position{line: 62, col: 40, offset: 1531},
							name: "Directives",
						},
					},
				},
			},
		},
		{
			name: "FragmentDefinition",
			pos:  position{line: 64, col: 1, offset: 1544},
			expr: &seqExpr{
				pos: position{line: 64, col: 23, offset: 1566},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 64, col: 23, offset: 1566},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 64, col: 25, offset: 1568},
						val:        "fragment",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 64, col: 36, offset: 1579},
						name: "FragmentName",
					},
					&ruleRefExpr{
						pos:  position{line: 64, col: 49, offset: 1592},
						name: "TypeCondition",
					},
					&zeroOrOneExpr{
						pos: position{line: 64, col: 63, offset: 1606},
						expr: &ruleRefExpr{
							pos:  position{line: 64, col: 63, offset: 1606},
							name: "Directives",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 64, col: 75, offset: 1618},
						name: "SelectionSet",
					},
				},
			},
		},
		{
			name: "FragmentName",
			pos:  position{line: 66, col: 1, offset: 1632},
			expr: &seqExpr{
				pos: position{line: 66, col: 17, offset: 1648},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 66, col: 17, offset: 1648},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 66, col: 19, offset: 1650},
						name: "Name",
					},
				},
			},
		},
		{
			name: "TypeCondition",
			pos:  position{line: 68, col: 1, offset: 1675},
			expr: &seqExpr{
				pos: position{line: 68, col: 18, offset: 1692},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 68, col: 18, offset: 1692},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 68, col: 20, offset: 1694},
						val:        "on",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 25, offset: 1699},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 27, offset: 1701},
						name: "NamedType",
					},
				},
			},
		},
		{
			name: "InlineFragment",
			pos:  position{line: 70, col: 1, offset: 1712},
			expr: &seqExpr{
				pos: position{line: 70, col: 19, offset: 1730},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 70, col: 19, offset: 1730},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 70, col: 21, offset: 1732},
						val:        "...",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 70, col: 27, offset: 1738},
						expr: &ruleRefExpr{
							pos:  position{line: 70, col: 27, offset: 1738},
							name: "TypeCondition",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 70, col: 42, offset: 1753},
						expr: &ruleRefExpr{
							pos:  position{line: 70, col: 42, offset: 1753},
							name: "Directives",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 70, col: 54, offset: 1765},
						name: "SelectionSet",
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 73, col: 1, offset: 1789},
			expr: &choiceExpr{
				pos: position{line: 74, col: 7, offset: 1804},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 74, col: 7, offset: 1804},
						name: "Variable",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 7, offset: 1819},
						name: "FloatValue",
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 7, offset: 1836},
						name: "IntValue",
					},
					&ruleRefExpr{
						pos:  position{line: 77, col: 7, offset: 1851},
						name: "StringValue",
					},
					&ruleRefExpr{
						pos:  position{line: 78, col: 7, offset: 1869},
						name: "BooleanValue",
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 7, offset: 1888},
						name: "NullValue",
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 7, offset: 1904},
						name: "EnumValue",
					},
					&ruleRefExpr{
						pos:  position{line: 81, col: 7, offset: 1987},
						name: "ListValue",
					},
					&ruleRefExpr{
						pos:  position{line: 82, col: 7, offset: 2003},
						name: "ObjectValue",
					},
				},
			},
		},
		{
			name: "IntValue",
			pos:  position{line: 85, col: 1, offset: 2022},
			expr: &actionExpr{
				pos: position{line: 85, col: 13, offset: 2034},
				run: (*parser).callonIntValue1,
				expr: &ruleRefExpr{
					pos:  position{line: 85, col: 13, offset: 2034},
					name: "IntegerPart",
				},
			},
		},
		{
			name: "IntegerPart",
			pos:  position{line: 86, col: 1, offset: 2090},
			expr: &choiceExpr{
				pos: position{line: 87, col: 7, offset: 2111},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 87, col: 7, offset: 2111},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 87, col: 7, offset: 2111},
								expr: &ruleRefExpr{
									pos:  position{line: 87, col: 7, offset: 2111},
									name: "NegativeSign",
								},
							},
							&litMatcher{
								pos:        position{line: 87, col: 21, offset: 2125},
								val:        "0",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 88, col: 7, offset: 2135},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 88, col: 7, offset: 2135},
								expr: &ruleRefExpr{
									pos:  position{line: 88, col: 7, offset: 2135},
									name: "NegativeSign",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 88, col: 21, offset: 2149},
								name: "NonZeroDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 88, col: 34, offset: 2162},
								expr: &ruleRefExpr{
									pos:  position{line: 88, col: 34, offset: 2162},
									name: "Digit",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FloatValue",
			pos:  position{line: 91, col: 1, offset: 2176},
			expr: &actionExpr{
				pos: position{line: 92, col: 5, offset: 2194},
				run: (*parser).callonFloatValue1,
				expr: &choiceExpr{
					pos: position{line: 92, col: 7, offset: 2196},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 92, col: 7, offset: 2196},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 92, col: 7, offset: 2196},
									name: "IntegerPart",
								},
								&ruleRefExpr{
									pos:  position{line: 92, col: 19, offset: 2208},
									name: "FractionalPart",
								},
							},
						},
						&seqExpr{
							pos: position{line: 93, col: 7, offset: 2229},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 93, col: 7, offset: 2229},
									name: "IntegerPart",
								},
								&ruleRefExpr{
									pos:  position{line: 93, col: 19, offset: 2241},
									name: "ExponentPart",
								},
							},
						},
						&seqExpr{
							pos: position{line: 94, col: 7, offset: 2260},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 94, col: 7, offset: 2260},
									name: "IntegerPart",
								},
								&ruleRefExpr{
									pos:  position{line: 94, col: 19, offset: 2272},
									name: "FractionalPart",
								},
								&ruleRefExpr{
									pos:  position{line: 94, col: 34, offset: 2287},
									name: "ExponentPart",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FractionalPart",
			pos:  position{line: 96, col: 1, offset: 2352},
			expr: &seqExpr{
				pos: position{line: 96, col: 19, offset: 2370},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 96, col: 19, offset: 2370},
						val:        ".",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 96, col: 23, offset: 2374},
						expr: &ruleRefExpr{
							pos:  position{line: 96, col: 23, offset: 2374},
							name: "Digit",
						},
					},
				},
			},
		},
		{
			name: "ExponentPart",
			pos:  position{line: 97, col: 1, offset: 2381},
			expr: &seqExpr{
				pos: position{line: 97, col: 17, offset: 2397},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 97, col: 17, offset: 2397},
						name: "ExponentIndicator",
					},
					&zeroOrOneExpr{
						pos: position{line: 97, col: 35, offset: 2415},
						expr: &ruleRefExpr{
							pos:  position{line: 97, col: 35, offset: 2415},
							name: "Sign",
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 97, col: 41, offset: 2421},
						expr: &ruleRefExpr{
							pos:  position{line: 97, col: 41, offset: 2421},
							name: "Digit",
						},
					},
				},
			},
		},
		{
			name: "BooleanValue",
			pos:  position{line: 99, col: 1, offset: 2429},
			expr: &choiceExpr{
				pos: position{line: 100, col: 7, offset: 2451},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 100, col: 7, offset: 2451},
						run: (*parser).callonBooleanValue2,
						expr: &seqExpr{
							pos: position{line: 100, col: 7, offset: 2451},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 100, col: 7, offset: 2451},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 100, col: 9, offset: 2453},
									val:        "true",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 101, col: 7, offset: 2505},
						run: (*parser).callonBooleanValue6,
						expr: &seqExpr{
							pos: position{line: 101, col: 7, offset: 2505},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 101, col: 7, offset: 2505},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 101, col: 9, offset: 2507},
									val:        "false",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "StringValue",
			pos:  position{line: 104, col: 1, offset: 2562},
			expr: &choiceExpr{
				pos: position{line: 105, col: 7, offset: 2583},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 105, col: 7, offset: 2583},
						run: (*parser).callonStringValue2,
						expr: &seqExpr{
							pos: position{line: 105, col: 7, offset: 2583},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 105, col: 7, offset: 2583},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 105, col: 9, offset: 2585},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 105, col: 13, offset: 2589},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 106, col: 7, offset: 2635},
						run: (*parser).callonStringValue7,
						expr: &seqExpr{
							pos: position{line: 106, col: 7, offset: 2635},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 106, col: 7, offset: 2635},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 106, col: 9, offset: 2637},
									val:        "\"",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 106, col: 13, offset: 2641},
									label: "s",
									expr: &ruleRefExpr{
										pos:  position{line: 106, col: 15, offset: 2643},
										name: "StringChars",
									},
								},
								&litMatcher{
									pos:        position{line: 106, col: 27, offset: 2655},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 107, col: 7, offset: 2709},
						run: (*parser).callonStringValue14,
						expr: &seqExpr{
							pos: position{line: 107, col: 7, offset: 2709},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 107, col: 7, offset: 2709},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 107, col: 9, offset: 2711},
									val:        "\"",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 107, col: 13, offset: 2715},
									label: "s",
									expr: &ruleRefExpr{
										pos:  position{line: 107, col: 15, offset: 2717},
										name: "StringChars",
									},
								},
								&notExpr{
									pos: position{line: 107, col: 27, offset: 2729},
									expr: &litMatcher{
										pos:        position{line: 107, col: 28, offset: 2730},
										val:        "\"",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "StringChars",
			pos:  position{line: 112, col: 1, offset: 2851},
			expr: &actionExpr{
				pos: position{line: 112, col: 16, offset: 2866},
				run: (*parser).callonStringChars1,
				expr: &oneOrMoreExpr{
					pos: position{line: 112, col: 16, offset: 2866},
					expr: &ruleRefExpr{
						pos:  position{line: 112, col: 16, offset: 2866},
						name: "StringCharacter",
					},
				},
			},
		},
		{
			name: "StringCharacter",
			pos:  position{line: 114, col: 1, offset: 2915},
			expr: &actionExpr{
				pos: position{line: 116, col: 5, offset: 2996},
				run: (*parser).callonStringCharacter1,
				expr: &choiceExpr{
					pos: position{line: 116, col: 7, offset: 2998},
					alternatives: []interface{}{
						&charClassMatcher{
							pos:        position{line: 116, col: 7, offset: 2998},
							val:        "[\\u0009\\u0020\\u0021\\u0023-\\u005B\\u005D-\\uFFFF]",
							chars:      []rune{'\t', ' ', '!'},
							ranges:     []rune{'#', '[', ']', '\uffff'},
							ignoreCase: false,
							inverted:   false,
						},
						&seqExpr{
							pos: position{line: 117, col: 7, offset: 3051},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 117, col: 7, offset: 3051},
									val:        "\\",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 117, col: 12, offset: 3056},
									val:        "u",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 117, col: 16, offset: 3060},
									name: "EscapedUnicode",
								},
							},
						},
						&seqExpr{
							pos: position{line: 118, col: 7, offset: 3081},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 118, col: 7, offset: 3081},
									val:        "\\",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 118, col: 12, offset: 3086},
									name: "EscapedCharacter",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NullValue",
			pos:  position{line: 121, col: 1, offset: 3133},
			expr: &actionExpr{
				pos: position{line: 121, col: 14, offset: 3146},
				run: (*parser).callonNullValue1,
				expr: &seqExpr{
					pos: position{line: 121, col: 14, offset: 3146},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 121, col: 14, offset: 3146},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 121, col: 16, offset: 3148},
							val:        "null",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 123, col: 1, offset: 3191},
			expr: &actionExpr{
				pos: position{line: 123, col: 14, offset: 3204},
				run: (*parser).callonEnumValue1,
				expr: &ruleRefExpr{
					pos:  position{line: 123, col: 14, offset: 3204},
					name: "Name",
				},
			},
		},
		{
			name: "ListValue",
			pos:  position{line: 126, col: 1, offset: 3258},
			expr: &choiceExpr{
				pos: position{line: 127, col: 7, offset: 3277},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 127, col: 7, offset: 3277},
						run: (*parser).callonListValue2,
						expr: &seqExpr{
							pos: position{line: 127, col: 7, offset: 3277},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 127, col: 7, offset: 3277},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 127, col: 9, offset: 3279},
									val:        "[",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 127, col: 13, offset: 3283},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 127, col: 15, offset: 3285},
									val:        "]",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 128, col: 7, offset: 3330},
						run: (*parser).callonListValue8,
						expr: &seqExpr{
							pos: position{line: 128, col: 7, offset: 3330},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 128, col: 7, offset: 3330},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 128, col: 9, offset: 3332},
									val:        "[",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 128, col: 13, offset: 3336},
									label: "va",
									expr: &oneOrMoreExpr{
										pos: position{line: 128, col: 16, offset: 3339},
										expr: &seqExpr{
											pos: position{line: 128, col: 17, offset: 3340},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 128, col: 17, offset: 3340},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 128, col: 19, offset: 3342},
													name: "Value",
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 128, col: 27, offset: 3350},
									val:        "]",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 136, col: 7, offset: 3595},
						run: (*parser).callonListValue18,
						expr: &seqExpr{
							pos: position{line: 136, col: 7, offset: 3595},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 136, col: 7, offset: 3595},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 136, col: 9, offset: 3597},
									val:        "[",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 136, col: 13, offset: 3601},
									expr: &seqExpr{
										pos: position{line: 136, col: 14, offset: 3602},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 136, col: 14, offset: 3602},
												name: "_",
											},
											&ruleRefExpr{
												pos:  position{line: 136, col: 16, offset: 3604},
												name: "Value",
											},
										},
									},
								},
								&notExpr{
									pos: position{line: 136, col: 24, offset: 3612},
									expr: &litMatcher{
										pos:        position{line: 136, col: 25, offset: 3613},
										val:        "]",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ObjectValue",
			pos:  position{line: 141, col: 1, offset: 3736},
			expr: &choiceExpr{
				pos: position{line: 142, col: 7, offset: 3757},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 142, col: 7, offset: 3757},
						run: (*parser).callonObjectValue2,
						expr: &seqExpr{
							pos: position{line: 142, col: 7, offset: 3757},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 142, col: 7, offset: 3757},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 142, col: 9, offset: 3759},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 142, col: 13, offset: 3763},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 142, col: 15, offset: 3765},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 143, col: 7, offset: 3817},
						run: (*parser).callonObjectValue8,
						expr: &seqExpr{
							pos: position{line: 143, col: 7, offset: 3817},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 143, col: 7, offset: 3817},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 143, col: 9, offset: 3819},
									val:        "{",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 143, col: 13, offset: 3823},
									label: "oa",
									expr: &oneOrMoreExpr{
										pos: position{line: 143, col: 16, offset: 3826},
										expr: &ruleRefExpr{
											pos:  position{line: 143, col: 16, offset: 3826},
											name: "ObjectField",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 143, col: 29, offset: 3839},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 143, col: 31, offset: 3841},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 153, col: 7, offset: 4221},
						run: (*parser).callonObjectValue17,
						expr: &seqExpr{
							pos: position{line: 153, col: 7, offset: 4221},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 153, col: 7, offset: 4221},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 153, col: 9, offset: 4223},
									val:        "{",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 153, col: 13, offset: 4227},
									label: "oa",
									expr: &oneOrMoreExpr{
										pos: position{line: 153, col: 16, offset: 4230},
										expr: &ruleRefExpr{
											pos:  position{line: 153, col: 16, offset: 4230},
											name: "ObjectField",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 153, col: 29, offset: 4243},
									name: "_",
								},
								&notExpr{
									pos: position{line: 153, col: 31, offset: 4245},
									expr: &litMatcher{
										pos:        position{line: 153, col: 32, offset: 4246},
										val:        "}",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ObjectField",
			pos:  position{line: 158, col: 1, offset: 4375},
			expr: &actionExpr{
				pos: position{line: 158, col: 16, offset: 4390},
				run: (*parser).callonObjectField1,
				expr: &seqExpr{
					pos: position{line: 158, col: 16, offset: 4390},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 158, col: 16, offset: 4390},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 158, col: 18, offset: 4392},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 158, col: 20, offset: 4394},
								name: "Name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 158, col: 25, offset: 4399},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 158, col: 27, offset: 4401},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 158, col: 31, offset: 4405},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 158, col: 33, offset: 4407},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 158, col: 35, offset: 4409},
								name: "Value",
							},
						},
					},
				},
			},
		},
		{
			name: "Variable",
			pos:  position{line: 162, col: 1, offset: 4471},
			expr: &actionExpr{
				pos: position{line: 162, col: 13, offset: 4483},
				run: (*parser).callonVariable1,
				expr: &seqExpr{
					pos: position{line: 162, col: 13, offset: 4483},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 162, col: 13, offset: 4483},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 162, col: 15, offset: 4485},
							val:        "$",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 162, col: 19, offset: 4489},
							name: "Name",
						},
					},
				},
			},
		},
		{
			name: "VariableDefinitions",
			pos:  position{line: 164, col: 1, offset: 4538},
			expr: &seqExpr{
				pos: position{line: 164, col: 24, offset: 4561},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 164, col: 24, offset: 4561},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 164, col: 26, offset: 4563},
						val:        "(",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 164, col: 30, offset: 4567},
						name: "_",
					},
					&oneOrMoreExpr{
						pos: position{line: 164, col: 32, offset: 4569},
						expr: &ruleRefExpr{
							pos:  position{line: 164, col: 32, offset: 4569},
							name: "VariableDefinition",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 164, col: 52, offset: 4589},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 164, col: 54, offset: 4591},
						val:        ")",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "VariableDefinition",
			pos:  position{line: 166, col: 1, offset: 4596},
			expr: &actionExpr{
				pos: position{line: 166, col: 23, offset: 4618},
				run: (*parser).callonVariableDefinition1,
				expr: &seqExpr{
					pos: position{line: 166, col: 23, offset: 4618},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 166, col: 23, offset: 4618},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 166, col: 25, offset: 4620},
							label: "vi",
							expr: &ruleRefExpr{
								pos:  position{line: 166, col: 28, offset: 4623},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 166, col: 37, offset: 4632},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 166, col: 39, offset: 4634},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 166, col: 43, offset: 4638},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 166, col: 45, offset: 4640},
							label: "ti",
							expr: &ruleRefExpr{
								pos:  position{line: 166, col: 48, offset: 4643},
								name: "Type",
							},
						},
						&labeledExpr{
							pos:   position{line: 166, col: 53, offset: 4648},
							label: "dvi",
							expr: &zeroOrOneExpr{
								pos: position{line: 166, col: 57, offset: 4652},
								expr: &ruleRefExpr{
									pos:  position{line: 166, col: 57, offset: 4652},
									name: "DefaultValue",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DefaultValue",
			pos:  position{line: 176, col: 1, offset: 4810},
			expr: &seqExpr{
				pos: position{line: 176, col: 17, offset: 4826},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 176, col: 17, offset: 4826},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 176, col: 19, offset: 4828},
						val:        "=",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 176, col: 23, offset: 4832},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 176, col: 25, offset: 4834},
						name: "Value",
					},
				},
			},
		},
		{
			name: "Type",
			pos:  position{line: 179, col: 1, offset: 4850},
			expr: &choiceExpr{
				pos: position{line: 180, col: 7, offset: 4864},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 180, col: 7, offset: 4864},
						name: "NonNullType",
					},
					&ruleRefExpr{
						pos:  position{line: 181, col: 7, offset: 4882},
						name: "ListType",
					},
					&ruleRefExpr{
						pos:  position{line: 182, col: 7, offset: 4897},
						name: "NamedType",
					},
				},
			},
		},
		{
			name: "NamedType",
			pos:  position{line: 185, col: 1, offset: 4914},
			expr: &actionExpr{
				pos: position{line: 185, col: 14, offset: 4927},
				run: (*parser).callonNamedType1,
				expr: &labeledExpr{
					pos:   position{line: 185, col: 14, offset: 4927},
					label: "n",
					expr: &ruleRefExpr{
						pos:  position{line: 185, col: 16, offset: 4929},
						name: "Name",
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 187, col: 1, offset: 4970},
			expr: &actionExpr{
				pos: position{line: 187, col: 13, offset: 4982},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 187, col: 13, offset: 4982},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 187, col: 13, offset: 4982},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 187, col: 15, offset: 4984},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 187, col: 19, offset: 4988},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 187, col: 21, offset: 4990},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 187, col: 23, offset: 4992},
								name: "Type",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 187, col: 28, offset: 4997},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 187, col: 30, offset: 4999},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "NonNullType",
			pos:  position{line: 189, col: 1, offset: 5048},
			expr: &choiceExpr{
				pos: position{line: 190, col: 7, offset: 5069},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 190, col: 7, offset: 5069},
						run: (*parser).callonNonNullType2,
						expr: &seqExpr{
							pos: position{line: 190, col: 7, offset: 5069},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 190, col: 7, offset: 5069},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 190, col: 9, offset: 5071},
										name: "NamedType",
									},
								},
								&litMatcher{
									pos:        position{line: 190, col: 19, offset: 5081},
									val:        "!",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 191, col: 7, offset: 5138},
						run: (*parser).callonNonNullType7,
						expr: &seqExpr{
							pos: position{line: 191, col: 7, offset: 5138},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 191, col: 7, offset: 5138},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 191, col: 9, offset: 5140},
										name: "ListType",
									},
								},
								&litMatcher{
									pos:        position{line: 191, col: 18, offset: 5149},
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
			name: "Directives",
			pos:  position{line: 195, col: 1, offset: 5221},
			expr: &oneOrMoreExpr{
				pos: position{line: 195, col: 15, offset: 5235},
				expr: &ruleRefExpr{
					pos:  position{line: 195, col: 15, offset: 5235},
					name: "Directive",
				},
			},
		},
		{
			name: "Directive",
			pos:  position{line: 197, col: 1, offset: 5247},
			expr: &seqExpr{
				pos: position{line: 197, col: 14, offset: 5260},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 197, col: 14, offset: 5260},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 197, col: 16, offset: 5262},
						val:        "@",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 20, offset: 5266},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 25, offset: 5271},
						name: "_",
					},
					&zeroOrOneExpr{
						pos: position{line: 197, col: 27, offset: 5273},
						expr: &ruleRefExpr{
							pos:  position{line: 197, col: 27, offset: 5273},
							name: "Arguments",
						},
					},
				},
			},
		},
		{
			name: "Token",
			pos:  position{line: 201, col: 1, offset: 5305},
			expr: &choiceExpr{
				pos: position{line: 202, col: 7, offset: 5320},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 202, col: 7, offset: 5320},
						name: "Punctuator",
					},
					&ruleRefExpr{
						pos:  position{line: 203, col: 7, offset: 5337},
						name: "Name",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 7, offset: 5348},
						name: "IntValue",
					},
					&ruleRefExpr{
						pos:  position{line: 205, col: 7, offset: 5363},
						name: "FloatValue",
					},
					&ruleRefExpr{
						pos:  position{line: 206, col: 7, offset: 5380},
						name: "StringValue",
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 210, col: 1, offset: 5417},
			expr: &zeroOrMoreExpr{
				pos: position{line: 210, col: 6, offset: 5422},
				expr: &ruleRefExpr{
					pos:  position{line: 210, col: 6, offset: 5422},
					name: "I",
				},
			},
		},
		{
			name:        "I",
			displayName: "\"Ignored\"",
			pos:         position{line: 211, col: 1, offset: 5425},
			expr: &choiceExpr{
				pos: position{line: 212, col: 6, offset: 5445},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 212, col: 6, offset: 5445},
						name: "UnicodeBOM",
					},
					&ruleRefExpr{
						pos:  position{line: 213, col: 6, offset: 5461},
						name: "WhiteSpace",
					},
					&ruleRefExpr{
						pos:  position{line: 214, col: 6, offset: 5477},
						name: "LineTerminator",
					},
					&ruleRefExpr{
						pos:  position{line: 215, col: 6, offset: 5497},
						name: "Comment",
					},
					&ruleRefExpr{
						pos:  position{line: 216, col: 6, offset: 5510},
						name: "Comma",
					},
				},
			},
		},
		{
			name: "SourceCharacter",
			pos:  position{line: 221, col: 1, offset: 5557},
			expr: &charClassMatcher{
				pos:        position{line: 222, col: 5, offset: 5580},
				val:        "[\\u0009\\n\\u000D\\u0020-\\uFFFF]",
				chars:      []rune{'\t', '\n', '\r'},
				ranges:     []rune{' ', '\uffff'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name:        "UnicodeBOM",
			displayName: "\"ByteOrderMark\"",
			pos:         position{line: 224, col: 1, offset: 5611},
			expr: &litMatcher{
				pos:        position{line: 225, col: 5, offset: 5645},
				val:        "\ufeff",
				ignoreCase: false,
			},
		},
		{
			name:        "LineTerminator",
			displayName: "\"EOL\"",
			pos:         position{line: 227, col: 1, offset: 5655},
			expr: &choiceExpr{
				pos: position{line: 228, col: 7, offset: 5685},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 228, col: 7, offset: 5685},
						val:        "\n",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 229, col: 7, offset: 5696},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 229, col: 7, offset: 5696},
								val:        "\r",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 229, col: 12, offset: 5701},
								expr: &litMatcher{
									pos:        position{line: 229, col: 13, offset: 5702},
									val:        "\n",
									ignoreCase: false,
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 230, col: 7, offset: 5713},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 230, col: 7, offset: 5713},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 230, col: 12, offset: 5718},
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
			pos:  position{line: 233, col: 1, offset: 5730},
			expr: &charClassMatcher{
				pos:        position{line: 233, col: 15, offset: 5744},
				val:        "[ \\t]",
				chars:      []rune{' ', '\t'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 234, col: 1, offset: 5750},
			expr: &seqExpr{
				pos: position{line: 234, col: 12, offset: 5761},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 234, col: 12, offset: 5761},
						val:        "#",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 234, col: 16, offset: 5765},
						expr: &anyMatcher{
							line: 234, col: 16, offset: 5765,
						},
					},
					&notExpr{
						pos: position{line: 234, col: 19, offset: 5768},
						expr: &ruleRefExpr{
							pos:  position{line: 234, col: 20, offset: 5769},
							name: "LineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Comma",
			pos:  position{line: 235, col: 1, offset: 5784},
			expr: &litMatcher{
				pos:        position{line: 235, col: 10, offset: 5793},
				val:        ",",
				ignoreCase: false,
			},
		},
		{
			name: "NegativeSign",
			pos:  position{line: 236, col: 1, offset: 5797},
			expr: &litMatcher{
				pos:        position{line: 236, col: 17, offset: 5813},
				val:        "-",
				ignoreCase: false,
			},
		},
		{
			name: "Digit",
			pos:  position{line: 237, col: 1, offset: 5817},
			expr: &charClassMatcher{
				pos:        position{line: 237, col: 10, offset: 5826},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDigit",
			pos:  position{line: 238, col: 1, offset: 5832},
			expr: &charClassMatcher{
				pos:        position{line: 238, col: 17, offset: 5848},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "ExponentIndicator",
			pos:  position{line: 239, col: 1, offset: 5854},
			expr: &charClassMatcher{
				pos:        position{line: 239, col: 22, offset: 5875},
				val:        "[eE]",
				chars:      []rune{'e', 'E'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Sign",
			pos:  position{line: 240, col: 1, offset: 5880},
			expr: &charClassMatcher{
				pos:        position{line: 240, col: 9, offset: 5888},
				val:        "[+-]",
				chars:      []rune{'+', '-'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapedUnicode",
			pos:  position{line: 241, col: 1, offset: 5893},
			expr: &charClassMatcher{
				pos:        position{line: 241, col: 19, offset: 5911},
				val:        "[0-9A-Fa-f]",
				ranges:     []rune{'0', '9', 'A', 'F', 'a', 'f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapedCharacter",
			pos:  position{line: 242, col: 1, offset: 5928},
			expr: &charClassMatcher{
				pos:        position{line: 242, col: 21, offset: 5948},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Punctuator",
			pos:  position{line: 244, col: 1, offset: 5961},
			expr: &choiceExpr{
				pos: position{line: 245, col: 7, offset: 5981},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 245, col: 7, offset: 5981},
						val:        "[!$():=@[\\]{|}]",
						chars:      []rune{'!', '$', '(', ')', ':', '=', '@', '[', ']', '{', '|', '}'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 246, col: 7, offset: 6003},
						val:        "...",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Name",
			pos:  position{line: 249, col: 1, offset: 6016},
			expr: &actionExpr{
				pos: position{line: 249, col: 9, offset: 6024},
				run: (*parser).callonName1,
				expr: &seqExpr{
					pos: position{line: 249, col: 9, offset: 6024},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 249, col: 9, offset: 6024},
							val:        "[_A-Za-z]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 249, col: 18, offset: 6033},
							expr: &charClassMatcher{
								pos:        position{line: 249, col: 18, offset: 6033},
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
			pos:  position{line: 251, col: 1, offset: 6079},
			expr: &notExpr{
				pos: position{line: 251, col: 7, offset: 6085},
				expr: &anyMatcher{
					line: 251, col: 8, offset: 6086,
				},
			},
		},
	},
}

func (c *current) onOperationType8(n interface{}) (interface{}, error) {
	return n, c.NewError("Must be query or mutation.", "OT-1")
}

func (p *parser) callonOperationType8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperationType8(stack["n"])
}

func (c *current) onIntValue1() (interface{}, error) {
	return typereg.ParseInt(string(c.text))
}

func (p *parser) callonIntValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntValue1()
}

func (c *current) onFloatValue1() (interface{}, error) {
	return typereg.ParseFloat(string(c.text))
}

func (p *parser) callonFloatValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFloatValue1()
}

func (c *current) onBooleanValue2() (interface{}, error) {
	return typereg.MakeBool(true), nil
}

func (p *parser) callonBooleanValue2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooleanValue2()
}

func (c *current) onBooleanValue6() (interface{}, error) {
	return typereg.MakeBool(false), nil
}

func (p *parser) callonBooleanValue6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooleanValue6()
}

func (c *current) onStringValue2() (interface{}, error) {
	return typereg.MakeStr(""), nil
}

func (p *parser) callonStringValue2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringValue2()
}

func (c *current) onStringValue7(s interface{}) (interface{}, error) {
	return typereg.MakeStr(s.(string)), nil
}

func (p *parser) callonStringValue7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringValue7(stack["s"])
}

func (c *current) onStringValue14(s interface{}) (interface{}, error) {
	return typereg.MakeStr(s.(string)), c.NewError("Missing close quote for string.", "SV-1")

}

func (p *parser) callonStringValue14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringValue14(stack["s"])
}

func (c *current) onStringChars1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonStringChars1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringChars1()
}

func (c *current) onStringCharacter1() (interface{}, error) {
	return c.text, nil
}

func (p *parser) callonStringCharacter1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringCharacter1()
}

func (c *current) onNullValue1() (interface{}, error) {
	return typereg.MakeNull(), nil
}

func (p *parser) callonNullValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNullValue1()
}

func (c *current) onEnumValue1() (interface{}, error) {
	return typereg.FindEnum(string(c.text))
}

func (p *parser) callonEnumValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValue1()
}

func (c *current) onListValue2() (interface{}, error) {
	return typereg.MakeList(), nil
}

func (p *parser) callonListValue2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListValue2()
}

func (c *current) onListValue8(va interface{}) (interface{}, error) {
	values := va.([]interface{})
	list := typereg.MakeListOf(values[0].(Value).T)
	for _, v := range values {
		list.Append(v.(Value))
	}
	return list, nil

}

func (p *parser) callonListValue8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListValue8(stack["va"])
}

func (c *current) onListValue18() (interface{}, error) {
	return typereg.MakeList(), c.NewError("Missing close square bracket for list.", "LV-1")

}

func (p *parser) callonListValue18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListValue18()
}

func (c *current) onObjectValue2() (interface{}, error) {
	return typereg.MakeNamelessObj(), nil
}

func (p *parser) callonObjectValue2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObjectValue2()
}

func (c *current) onObjectValue8(oa interface{}) (interface{}, error) {
	obj := typereg.MakeNamelessObj()
	fields := oa.([]interface{})
	for _, v := range fields {
		// the ObjectField is a Value where the N field is the name
		// of the field, and the V field is a nested Value object
		obj.SetField(v.(Value))
	}
	return obj, nil

}

func (p *parser) callonObjectValue8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObjectValue8(stack["oa"])
}

func (c *current) onObjectValue17(oa interface{}) (interface{}, error) {
	return typereg.MakeNamelessObj(), c.NewError("Missing close curly brace for object.", "OV-1")

}

func (p *parser) callonObjectValue17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObjectValue17(stack["oa"])
}

func (c *current) onObjectField1(n, v interface{}) (interface{}, error) {
	return Value{N: n.(string), V: v}, nil
}

func (p *parser) callonObjectField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObjectField1(stack["n"], stack["v"])
}

func (c *current) onVariable1() (interface{}, error) {
	return Variable{N: string(c.text)}, nil
}

func (p *parser) callonVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVariable1()
}

func (c *current) onVariableDefinition1(vi, ti, dvi interface{}) (interface{}, error) {
	v := vi.(Variable)
	t := ti.(Type)
	dv := dvi.(Value)
	v.T = t
	v.V = dv
	return v, nil

}

func (p *parser) callonVariableDefinition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVariableDefinition1(stack["vi"], stack["ti"], stack["dvi"])
}

func (c *current) onNamedType1(n interface{}) (interface{}, error) {
	return typereg.Get(n.(string))
}

func (p *parser) callonNamedType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamedType1(stack["n"])
}

func (c *current) onListType1(t interface{}) (interface{}, error) {
	return ListType{Contains: t.(Type)}, nil
}

func (p *parser) callonListType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListType1(stack["t"])
}

func (c *current) onNonNullType2(t interface{}) (interface{}, error) {
	return NonNullType{Contains: t.(Type)}, nil
}

func (p *parser) callonNonNullType2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonNullType2(stack["t"])
}

func (c *current) onNonNullType7(t interface{}) (interface{}, error) {
	return NonNullType{Contains: t.(Type)}, nil
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
