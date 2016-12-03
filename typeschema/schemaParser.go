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
				pos: position{line: 24, col: 13, offset: 849},
				exprs: []interface{}{
					&oneOrMoreExpr{
						pos: position{line: 24, col: 13, offset: 849},
						expr: &ruleRefExpr{
							pos:  position{line: 24, col: 13, offset: 849},
							name: "Definition",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 24, col: 25, offset: 861},
						name: "_",
					},
					&ruleRefExpr{
						pos:  position{line: 24, col: 27, offset: 863},
						name: "EOF",
					},
				},
			},
		},
		{
			name: "Definition",
			pos:  position{line: 26, col: 1, offset: 868},
			expr: &choiceExpr{
				pos: position{line: 27, col: 7, offset: 888},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 27, col: 7, offset: 888},
						name: "ObjTypeDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 28, col: 7, offset: 912},
						name: "EnumDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 29, col: 7, offset: 933},
						name: "InterfaceDefinition",
					},
					&ruleRefExpr{
						pos:  position{line: 30, col: 7, offset: 959},
						name: "UnionDefinition",
					},
				},
			},
		},
		{
			name: "ObjTypeDefinition",
			pos:  position{line: 33, col: 1, offset: 982},
			expr: &actionExpr{
				pos: position{line: 33, col: 22, offset: 1003},
				run: (*parser).callonObjTypeDefinition1,
				expr: &seqExpr{
					pos: position{line: 33, col: 22, offset: 1003},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 33, col: 22, offset: 1003},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 33, col: 24, offset: 1005},
							val:        "type",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 33, col: 31, offset: 1012},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 33, col: 33, offset: 1014},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 35, offset: 1016},
								name: "Name",
							},
						},
						&labeledExpr{
							pos:   position{line: 33, col: 40, offset: 1021},
							label: "i",
							expr: &zeroOrOneExpr{
								pos: position{line: 33, col: 42, offset: 1023},
								expr: &ruleRefExpr{
									pos:  position{line: 33, col: 42, offset: 1023},
									name: "Implements",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 33, col: 54, offset: 1035},
							label: "fs",
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 57, offset: 1038},
								name: "FieldSet",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumDefinition",
			pos:  position{line: 41, col: 1, offset: 1243},
			expr: &actionExpr{
				pos: position{line: 41, col: 19, offset: 1261},
				run: (*parser).callonEnumDefinition1,
				expr: &seqExpr{
					pos: position{line: 41, col: 19, offset: 1261},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 41, col: 19, offset: 1261},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 41, col: 21, offset: 1263},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 41, col: 28, offset: 1270},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 41, col: 30, offset: 1272},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 41, col: 32, offset: 1274},
								name: "Name",
							},
						},
						&labeledExpr{
							pos:   position{line: 41, col: 37, offset: 1279},
							label: "vals",
							expr: &ruleRefExpr{
								pos:  position{line: 41, col: 42, offset: 1284},
								name: "EnumValueSet",
							},
						},
					},
				},
			},
		},
		{
			name: "InterfaceDefinition",
			pos:  position{line: 58, col: 1, offset: 1724},
			expr: &actionExpr{
				pos: position{line: 58, col: 24, offset: 1747},
				run: (*parser).callonInterfaceDefinition1,
				expr: &seqExpr{
					pos: position{line: 58, col: 24, offset: 1747},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 58, col: 24, offset: 1747},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 58, col: 26, offset: 1749},
							val:        "interface",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 58, col: 38, offset: 1761},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 58, col: 40, offset: 1763},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 58, col: 42, offset: 1765},
								name: "Name",
							},
						},
						&labeledExpr{
							pos:   position{line: 58, col: 47, offset: 1770},
							label: "fs",
							expr: &ruleRefExpr{
								pos:  position{line: 58, col: 50, offset: 1773},
								name: "FieldSet",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionDefinition",
			pos:  position{line: 66, col: 1, offset: 1984},
			expr: &actionExpr{
				pos: position{line: 66, col: 20, offset: 2003},
				run: (*parser).callonUnionDefinition1,
				expr: &seqExpr{
					pos: position{line: 66, col: 20, offset: 2003},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 66, col: 20, offset: 2003},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 66, col: 22, offset: 2005},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 66, col: 30, offset: 2013},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 66, col: 32, offset: 2015},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 66, col: 34, offset: 2017},
								name: "Name",
							},
						},
						&labeledExpr{
							pos:   position{line: 66, col: 39, offset: 2022},
							label: "us",
							expr: &ruleRefExpr{
								pos:  position{line: 66, col: 42, offset: 2025},
								name: "UnionSet",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionSet",
			pos:  position{line: 72, col: 1, offset: 2165},
			expr: &actionExpr{
				pos: position{line: 72, col: 13, offset: 2177},
				run: (*parser).callonUnionSet1,
				expr: &seqExpr{
					pos: position{line: 72, col: 13, offset: 2177},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 72, col: 13, offset: 2177},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 72, col: 15, offset: 2179},
							val:        "=",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 72, col: 19, offset: 2183},
							label: "un",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 22, offset: 2186},
								name: "UnionNames",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionNames",
			pos:  position{line: 74, col: 1, offset: 2217},
			expr: &choiceExpr{
				pos: position{line: 75, col: 7, offset: 2237},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 75, col: 7, offset: 2237},
						run: (*parser).callonUnionNames2,
						expr: &seqExpr{
							pos: position{line: 75, col: 7, offset: 2237},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 75, col: 7, offset: 2237},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 75, col: 9, offset: 2239},
									label: "fn",
									expr: &ruleRefExpr{
										pos:  position{line: 75, col: 12, offset: 2242},
										name: "Name",
									},
								},
								&labeledExpr{
									pos:   position{line: 75, col: 17, offset: 2247},
									label: "mn",
									expr: &oneOrMoreExpr{
										pos: position{line: 75, col: 20, offset: 2250},
										expr: &ruleRefExpr{
											pos:  position{line: 75, col: 20, offset: 2250},
											name: "AnotherName",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 92, col: 7, offset: 2795},
						run: (*parser).callonUnionNames10,
						expr: &seqExpr{
							pos: position{line: 92, col: 7, offset: 2795},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 92, col: 7, offset: 2795},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 92, col: 9, offset: 2797},
									label: "fn",
									expr: &ruleRefExpr{
										pos:  position{line: 92, col: 12, offset: 2800},
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
			name: "Implements",
			pos:  position{line: 102, col: 1, offset: 3088},
			expr: &actionExpr{
				pos: position{line: 102, col: 15, offset: 3102},
				run: (*parser).callonImplements1,
				expr: &seqExpr{
					pos: position{line: 102, col: 15, offset: 3102},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 102, col: 15, offset: 3102},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 102, col: 17, offset: 3104},
							val:        "implements",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 102, col: 30, offset: 3117},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 102, col: 32, offset: 3119},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 102, col: 34, offset: 3121},
								name: "Name",
							},
						},
					},
				},
			},
		},
		{
			name: "EnumValueSet",
			pos:  position{line: 106, col: 1, offset: 3174},
			expr: &actionExpr{
				pos: position{line: 106, col: 17, offset: 3190},
				run: (*parser).callonEnumValueSet1,
				expr: &seqExpr{
					pos: position{line: 106, col: 17, offset: 3190},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 106, col: 17, offset: 3190},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 106, col: 19, offset: 3192},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 106, col: 23, offset: 3196},
							label: "ev",
							expr: &ruleRefExpr{
								pos:  position{line: 106, col: 26, offset: 3199},
								name: "EnumValues",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 106, col: 37, offset: 3210},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 106, col: 39, offset: 3212},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "EnumValues",
			pos:  position{line: 108, col: 1, offset: 3236},
			expr: &actionExpr{
				pos: position{line: 108, col: 15, offset: 3250},
				run: (*parser).callonEnumValues1,
				expr: &labeledExpr{
					pos:   position{line: 108, col: 15, offset: 3250},
					label: "evs",
					expr: &oneOrMoreExpr{
						pos: position{line: 108, col: 19, offset: 3254},
						expr: &ruleRefExpr{
							pos:  position{line: 108, col: 19, offset: 3254},
							name: "EnumValue",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 116, col: 1, offset: 3399},
			expr: &actionExpr{
				pos: position{line: 116, col: 14, offset: 3412},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 116, col: 14, offset: 3412},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 116, col: 14, offset: 3412},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 116, col: 16, offset: 3414},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 116, col: 18, offset: 3416},
								name: "Name",
							},
						},
					},
				},
			},
		},
		{
			name: "AnotherName",
			pos:  position{line: 118, col: 1, offset: 3449},
			expr: &actionExpr{
				pos: position{line: 118, col: 16, offset: 3464},
				run: (*parser).callonAnotherName1,
				expr: &seqExpr{
					pos: position{line: 118, col: 16, offset: 3464},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 118, col: 16, offset: 3464},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 118, col: 18, offset: 3466},
							val:        "|",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 22, offset: 3470},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 118, col: 24, offset: 3472},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 26, offset: 3474},
								name: "Name",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldSet",
			pos:  position{line: 120, col: 1, offset: 3507},
			expr: &actionExpr{
				pos: position{line: 120, col: 14, offset: 3520},
				run: (*parser).callonFieldSet1,
				expr: &seqExpr{
					pos: position{line: 120, col: 14, offset: 3520},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 120, col: 14, offset: 3520},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 120, col: 16, offset: 3522},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 120, col: 20, offset: 3526},
							label: "fields",
							expr: &oneOrMoreExpr{
								pos: position{line: 120, col: 27, offset: 3533},
								expr: &ruleRefExpr{
									pos:  position{line: 120, col: 27, offset: 3533},
									name: "Field",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 120, col: 34, offset: 3540},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 120, col: 36, offset: 3542},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 128, col: 1, offset: 3709},
			expr: &actionExpr{
				pos: position{line: 128, col: 10, offset: 3718},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 128, col: 10, offset: 3718},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 128, col: 10, offset: 3718},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 128, col: 12, offset: 3720},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 128, col: 14, offset: 3722},
								name: "Name",
							},
						},
						&labeledExpr{
							pos:   position{line: 128, col: 19, offset: 3727},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 128, col: 24, offset: 3732},
								expr: &ruleRefExpr{
									pos:  position{line: 128, col: 24, offset: 3732},
									name: "Arguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 128, col: 35, offset: 3743},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 128, col: 37, offset: 3745},
							val:        ":",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 128, col: 41, offset: 3749},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 128, col: 43, offset: 3751},
								name: "Type",
							},
						},
					},
				},
			},
		},
		{
			name: "Arguments",
			pos:  position{line: 135, col: 1, offset: 3931},
			expr: &actionExpr{
				pos: position{line: 135, col: 14, offset: 3944},
				run: (*parser).callonArguments1,
				expr: &seqExpr{
					pos: position{line: 135, col: 14, offset: 3944},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 135, col: 14, offset: 3944},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 135, col: 16, offset: 3946},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 135, col: 20, offset: 3950},
							label: "args",
							expr: &oneOrMoreExpr{
								pos: position{line: 135, col: 25, offset: 3955},
								expr: &ruleRefExpr{
									pos:  position{line: 135, col: 25, offset: 3955},
									name: "Argument",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 35, offset: 3965},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 135, col: 37, offset: 3967},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Argument",
			pos:  position{line: 143, col: 1, offset: 4148},
			expr: &actionExpr{
				pos: position{line: 143, col: 13, offset: 4160},
				run: (*parser).callonArgument1,
				expr: &seqExpr{
					pos: position{line: 143, col: 13, offset: 4160},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 143, col: 13, offset: 4160},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 143, col: 15, offset: 4162},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 17, offset: 4164},
								name: "Name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 143, col: 22, offset: 4169},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 143, col: 24, offset: 4171},
							val:        ":",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 143, col: 28, offset: 4175},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 30, offset: 4177},
								name: "Type",
							},
						},
					},
				},
			},
		},
		{
			name: "Type",
			pos:  position{line: 150, col: 1, offset: 4286},
			expr: &choiceExpr{
				pos: position{line: 151, col: 7, offset: 4300},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 151, col: 7, offset: 4300},
						name: "NonNullType",
					},
					&ruleRefExpr{
						pos:  position{line: 152, col: 7, offset: 4318},
						name: "ListType",
					},
					&ruleRefExpr{
						pos:  position{line: 153, col: 7, offset: 4333},
						name: "NamedType",
					},
				},
			},
		},
		{
			name: "NamedType",
			pos:  position{line: 156, col: 1, offset: 4350},
			expr: &actionExpr{
				pos: position{line: 156, col: 14, offset: 4363},
				run: (*parser).callonNamedType1,
				expr: &seqExpr{
					pos: position{line: 156, col: 14, offset: 4363},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 156, col: 14, offset: 4363},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 156, col: 16, offset: 4365},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 156, col: 18, offset: 4367},
								name: "Name",
							},
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 158, col: 1, offset: 4408},
			expr: &actionExpr{
				pos: position{line: 158, col: 13, offset: 4420},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 158, col: 13, offset: 4420},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 158, col: 13, offset: 4420},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 158, col: 15, offset: 4422},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 158, col: 19, offset: 4426},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 158, col: 21, offset: 4428},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 158, col: 23, offset: 4430},
								name: "Type",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 158, col: 28, offset: 4435},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 158, col: 30, offset: 4437},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "NonNullType",
			pos:  position{line: 160, col: 1, offset: 4500},
			expr: &choiceExpr{
				pos: position{line: 161, col: 7, offset: 4521},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 161, col: 7, offset: 4521},
						run: (*parser).callonNonNullType2,
						expr: &seqExpr{
							pos: position{line: 161, col: 7, offset: 4521},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 161, col: 7, offset: 4521},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 161, col: 9, offset: 4523},
										name: "ListType",
									},
								},
								&litMatcher{
									pos:        position{line: 161, col: 18, offset: 4532},
									val:        "!",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 162, col: 7, offset: 4608},
						run: (*parser).callonNonNullType7,
						expr: &seqExpr{
							pos: position{line: 162, col: 7, offset: 4608},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 162, col: 7, offset: 4608},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 162, col: 9, offset: 4610},
										name: "NamedType",
									},
								},
								&litMatcher{
									pos:        position{line: 162, col: 19, offset: 4620},
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
			pos:  position{line: 167, col: 1, offset: 4715},
			expr: &zeroOrMoreExpr{
				pos: position{line: 167, col: 6, offset: 4720},
				expr: &ruleRefExpr{
					pos:  position{line: 167, col: 6, offset: 4720},
					name: "I",
				},
			},
		},
		{
			name:        "I",
			displayName: "\"Ignored\"",
			pos:         position{line: 168, col: 1, offset: 4723},
			expr: &choiceExpr{
				pos: position{line: 169, col: 6, offset: 4743},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 169, col: 6, offset: 4743},
						name: "UnicodeBOM",
					},
					&ruleRefExpr{
						pos:  position{line: 170, col: 6, offset: 4759},
						name: "WhiteSpace",
					},
					&ruleRefExpr{
						pos:  position{line: 171, col: 6, offset: 4775},
						name: "LineTerminator",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 6, offset: 4795},
						name: "Comment",
					},
					&ruleRefExpr{
						pos:  position{line: 173, col: 6, offset: 4808},
						name: "Comma",
					},
				},
			},
		},
		{
			name:        "UnicodeBOM",
			displayName: "\"ByteOrderMark\"",
			pos:         position{line: 178, col: 1, offset: 4855},
			expr: &litMatcher{
				pos:        position{line: 179, col: 5, offset: 4889},
				val:        "\ufeff",
				ignoreCase: false,
			},
		},
		{
			name:        "LineTerminator",
			displayName: "\"EOL\"",
			pos:         position{line: 181, col: 1, offset: 4899},
			expr: &choiceExpr{
				pos: position{line: 182, col: 7, offset: 4929},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 182, col: 7, offset: 4929},
						val:        "\n",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 183, col: 7, offset: 4940},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 183, col: 7, offset: 4940},
								val:        "\r",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 183, col: 12, offset: 4945},
								expr: &litMatcher{
									pos:        position{line: 183, col: 13, offset: 4946},
									val:        "\n",
									ignoreCase: false,
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 184, col: 7, offset: 4957},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 184, col: 7, offset: 4957},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 184, col: 12, offset: 4962},
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
			pos:  position{line: 187, col: 1, offset: 4974},
			expr: &charClassMatcher{
				pos:        position{line: 187, col: 15, offset: 4988},
				val:        "[ \\t]",
				chars:      []rune{' ', '\t'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 188, col: 1, offset: 4994},
			expr: &seqExpr{
				pos: position{line: 188, col: 12, offset: 5005},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 188, col: 12, offset: 5005},
						val:        "#",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 188, col: 16, offset: 5009},
						expr: &anyMatcher{
							line: 188, col: 16, offset: 5009,
						},
					},
					&notExpr{
						pos: position{line: 188, col: 19, offset: 5012},
						expr: &ruleRefExpr{
							pos:  position{line: 188, col: 20, offset: 5013},
							name: "LineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Comma",
			pos:  position{line: 189, col: 1, offset: 5028},
			expr: &litMatcher{
				pos:        position{line: 189, col: 10, offset: 5037},
				val:        ",",
				ignoreCase: false,
			},
		},
		{
			name: "Name",
			pos:  position{line: 190, col: 1, offset: 5041},
			expr: &actionExpr{
				pos: position{line: 190, col: 9, offset: 5049},
				run: (*parser).callonName1,
				expr: &seqExpr{
					pos: position{line: 190, col: 9, offset: 5049},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 190, col: 9, offset: 5049},
							val:        "[_A-Za-z]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 190, col: 18, offset: 5058},
							expr: &charClassMatcher{
								pos:        position{line: 190, col: 18, offset: 5058},
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
			pos:  position{line: 192, col: 1, offset: 5104},
			expr: &notExpr{
				pos: position{line: 192, col: 7, offset: 5110},
				expr: &anyMatcher{
					line: 192, col: 8, offset: 5111,
				},
			},
		},
	},
}

func (c *current) onObjTypeDefinition1(n, i, fs interface{}) (interface{}, error) {
	fields := make([]Field, 0)
	for _, f := range fs.([]Field) {
		fields = append(fields, f)
	}
	return TR.RegisterWithFields(n.(string), Obj, fields)

}

func (p *parser) callonObjTypeDefinition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObjTypeDefinition1(stack["n"], stack["i"], stack["fs"])
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

func (c *current) onInterfaceDefinition1(n, fs interface{}) (interface{}, error) {
	fields := make([]Field, 0)
	for _, f := range fs.([]Field) {
		fields = append(fields, f)
	}
	return TR.RegisterWithFields(n.(string), Interface, fields)

}

func (p *parser) callonInterfaceDefinition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInterfaceDefinition1(stack["n"], stack["fs"])
}

func (c *current) onUnionDefinition1(n, us interface{}) (interface{}, error) {
	unionname := n.(string)
	subtypes := us.([]Type)
	return TR.Register(unionname, Union, subtypes...)

}

func (p *parser) callonUnionDefinition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionDefinition1(stack["n"], stack["us"])
}

func (c *current) onUnionSet1(un interface{}) (interface{}, error) {
	return un, nil
}

func (p *parser) callonUnionSet1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionSet1(stack["un"])
}

func (c *current) onUnionNames2(fn, mn interface{}) (interface{}, error) {
	firstname := fn.(string)
	firsttype, err := TR.MaybeGet(firstname)
	types := []Type{firsttype}
	if err != nil {
		return types, err
	}
	morenames := mn.([]interface{})
	for _, n := range morenames {
		t, err := TR.MaybeGet(n.(string))
		if err != nil {
			return types, err
		}
		types = append(types, t)
	}
	return types, nil

}

func (p *parser) callonUnionNames2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionNames2(stack["fn"], stack["mn"])
}

func (c *current) onUnionNames10(fn interface{}) (interface{}, error) {
	firstname := fn.(string)
	firsttype, _ := TR.MaybeGet(firstname)
	types := []Type{firsttype}
	return types, errors.New("A union must specify at least two types.")

}

func (p *parser) callonUnionNames10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionNames10(stack["fn"])
}

func (c *current) onImplements1(n interface{}) (interface{}, error) {
	return TR.MaybeGet(n.(string))

}

func (p *parser) callonImplements1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImplements1(stack["n"])
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

func (c *current) onAnotherName1(n interface{}) (interface{}, error) {
	return n.(string), nil
}

func (p *parser) callonAnotherName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnotherName1(stack["n"])
}

func (c *current) onFieldSet1(fields interface{}) (interface{}, error) {
	fa := make([]Field, 0)
	for _, f := range fields.([]interface{}) {
		fa = append(fa, f.(Field))
	}
	return fa, nil

}

func (p *parser) callonFieldSet1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldSet1(stack["fields"])
}

func (c *current) onField1(n, args, t interface{}) (interface{}, error) {
	if args == nil {
		return Field{N: n.(string), T: t.(Type)}, nil
	}
	return Field{N: n.(string), Args: args.([]Field), T: t.(Type)}, nil

}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["n"], stack["args"], stack["t"])
}

func (c *current) onArguments1(args interface{}) (interface{}, error) {
	fields := make([]Field, 0)
	for _, a := range args.([]interface{}) {
		fields = append(fields, a.(Field))
	}
	return fields, nil

}

func (p *parser) callonArguments1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArguments1(stack["args"])
}

func (c *current) onArgument1(n, t interface{}) (interface{}, error) {
	name := n.(string)
	typ := t.(Type)
	return Field{N: name, T: typ}, nil

}

func (p *parser) callonArgument1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArgument1(stack["n"], stack["t"])
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
