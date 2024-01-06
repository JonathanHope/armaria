package order

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
)

func TestMakeSymbolTable(t *testing.T) {
	type test struct {
		symbols string
		table   SymbolTable
	}

	tests := []test{
		{
			symbols: "abc",
			table: SymbolTable{
				numberToSymbol: []string{"a", "b", "c"},
				symbolToNumber: map[string]int{
					"a": 0,
					"b": 1,
					"c": 2,
				},
				base: 3,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.symbols, func(t *testing.T) {
			table, err := MakeSymbolTable(tc.symbols)
			if err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			unexported := cmp.AllowUnexported(SymbolTable{})
			diff := cmp.Diff(table, tc.table, unexported)
			if diff != "" {
				t.Errorf("Expected and actual symbol tables different:\n%s", diff)
			}
		})
	}
}

func TestNumberToDigits(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	hex, err := MakeSymbolTable(Hex)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table  SymbolTable
		number int
		digits []int
	}

	tests := []test{
		{table: decimal, number: 123, digits: []int{1, 2, 3}},
		{table: hex, number: 123, digits: []int{7, 11}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.number), func(t *testing.T) {
			digits := tc.table.numberToDigits(tc.number)

			diff := cmp.Diff(digits, tc.digits)
			if diff != "" {
				t.Errorf("Expected and actual digits different:\n%s", diff)
			}
		})
	}
}

func TestDigitsToString(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	hex, err := MakeSymbolTable(Hex)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table  SymbolTable
		digits []int
		str    string
	}

	tests := []test{
		{table: decimal, digits: []int{1, 2, 3}, str: "123"},
		{table: hex, digits: []int{7, 11}, str: "7b"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.digits), func(t *testing.T) {
			str, err := tc.table.digitsToString(tc.digits)
			if err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			diff := cmp.Diff(str, tc.str)
			if diff != "" {
				t.Errorf("Expected and actual strings different:\n%s", diff)
			}
		})
	}
}

func TestStringToDigits(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table  SymbolTable
		str    string
		digits []int
	}

	tests := []test{
		{table: decimal, str: "123", digits: []int{1, 2, 3}},
	}

	for _, tc := range tests {
		t.Run(tc.str, func(t *testing.T) {
			digits, err := tc.table.stringToDigits(tc.str)
			if err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			diff := cmp.Diff(digits, tc.digits)
			if diff != "" {
				t.Errorf("Expected and actual digits different:\n%s", diff)
			}
		})
	}
}

func TestDigitsToNumber(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	hex, err := MakeSymbolTable(Hex)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table  SymbolTable
		digits []int
		number int
	}

	tests := []test{
		{table: decimal, digits: []int{1, 2, 3}, number: 123},
		{table: hex, digits: []int{7, 11}, number: 123},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.digits), func(t *testing.T) {
			num := tc.table.digitsToNumber(tc.digits)

			diff := cmp.Diff(num, tc.number)
			if diff != "" {
				t.Errorf("Expected and actual numbers different:\n%s", diff)
			}
		})
	}
}

func TestNumberToString(t *testing.T) {
	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table  SymbolTable
		number int
		str    string
	}
	tests := []test{
		{table: base62, number: 0, str: "0"},
		{table: base62, number: 1, str: "1"},
		{table: base62, number: 9, str: "9"},
		{table: base62, number: 10, str: "A"},
		{table: base62, number: 35, str: "Z"},
		{table: base62, number: 36, str: "a"},
		{table: base62, number: 37, str: "b"},
		{table: base62, number: 61, str: "z"},
		{table: base62, number: 62, str: "10"},
		{table: base62, number: 63, str: "11"},
		{table: base62, number: 1945, str: "VN"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.number), func(t *testing.T) {
			str, err := tc.table.NumberToString(tc.number)
			if err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			diff := cmp.Diff(str, tc.str)
			if diff != "" {
				t.Errorf("Expected and actual strings different:\n%s", diff)
			}
		})
	}
}

func TestStringToNumber(t *testing.T) {
	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table  SymbolTable
		str    string
		number int
	}

	tests := []test{
		{table: base62, number: 0, str: "0"},
		{table: base62, number: 1, str: "1"},
		{table: base62, number: 9, str: "9"},
		{table: base62, number: 10, str: "A"},
		{table: base62, number: 35, str: "Z"},
		{table: base62, number: 36, str: "a"},
		{table: base62, number: 37, str: "b"},
		{table: base62, number: 61, str: "z"},
		{table: base62, number: 62, str: "10"},
		{table: base62, number: 63, str: "11"},
		{table: base62, number: 1945, str: "VN"},
	}

	for _, tc := range tests {
		t.Run(tc.str, func(t *testing.T) {
			num, err := tc.table.StringToNumber(tc.str)
			if err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			diff := cmp.Diff(num, tc.number)
			if diff != "" {
				t.Errorf("Expected and actual numbers different:\n%s", diff)
			}
		})
	}
}

func TestLongAddSameLength(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table       SymbolTable
		addend2     []int
		addend1     []int
		remainder   int
		denominator int
		result      longResult
	}

	tests := []test{
		{
			table:       decimal,
			addend1:     []int{1, 2},
			addend2:     []int{4, 5},
			remainder:   0,
			denominator: 1,
			result: longResult{
				result:    []int{5, 7},
				carry:     false,
				remainder: 0,
			},
		},
		{
			table:       decimal,
			addend1:     []int{1, 2},
			addend2:     []int{9, 9},
			remainder:   0,
			denominator: 1,
			result: longResult{
				result:    []int{1, 1},
				carry:     true,
				remainder: 0,
			},
		},
		{
			table:       decimal,
			addend1:     []int{1},
			addend2:     []int{1},
			remainder:   1,
			denominator: 1,
			result: longResult{
				result:    []int{3},
				carry:     false,
				remainder: 0,
			},
		},
		{
			table:       decimal,
			addend1:     []int{1},
			addend2:     []int{9},
			remainder:   1,
			denominator: 1,
			result: longResult{
				result:    []int{1},
				carry:     true,
				remainder: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(
			"addend1: %+v, addend2: %+v, remainder: %d, denominator: %d",
			tc.addend2, tc.addend1, tc.remainder, tc.denominator),
			func(t *testing.T) {
				result, err := tc.table.longAdd(tc.addend1, tc.addend2, tc.remainder, tc.denominator)
				if err != nil {
					t.Fatalf("Unexpected error occurred: %s", err)
				}

				unexported := cmp.AllowUnexported(longResult{})
				diff := cmp.Diff(result, tc.result, unexported)
				if diff != "" {
					t.Errorf("Expected and actual sums different:\n%s", diff)
				}
			})
	}
}

func TestLongDiv(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table    SymbolTable
		dividend []int
		divisor  int
		result   longResult
	}

	tests := []test{
		{
			table:    decimal,
			dividend: []int{1, 0},
			divisor:  2,
			result: longResult{
				result:    []int{0, 5},
				remainder: 0,
			},
		},
		{
			table:    decimal,
			dividend: []int{5, 0, 0},
			divisor:  4,
			result: longResult{
				result:    []int{1, 2, 5},
				remainder: 0,
			},
		},
		{
			table:    decimal,
			dividend: []int{7, 5},
			divisor:  4,
			result: longResult{
				result:    []int{1, 8},
				remainder: 3,
			},
		},
		{
			table:    decimal,
			dividend: []int{4, 3, 5},
			divisor:  25,
			result: longResult{
				result:    []int{0, 1, 7},
				remainder: 10,
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(
			"dividend: %+v, divisior: %+v",
			tc.dividend, tc.divisor),
			func(t *testing.T) {
				result := tc.table.longDiv(tc.dividend, tc.divisor)

				unexported := cmp.AllowUnexported(longResult{})
				diff := cmp.Diff(result, tc.result, unexported)
				if diff != "" {
					t.Errorf("Expected and actual quotients different:\n%s", diff)
				}
			})
	}
}

func TestLongSub(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table               SymbolTable
		minuend             []int
		subtrahend          []int
		minuendRemainder    int
		subtrahendRemainder int
		denominator         int
		result              longResult
	}

	tests := []test{
		{
			table:               decimal,
			minuend:             []int{6, 8},
			subtrahend:          []int{2, 5},
			minuendRemainder:    0,
			subtrahendRemainder: 0,
			denominator:         0,
			result: longResult{
				result:    []int{4, 3},
				remainder: 0,
			},
		},
		{
			table:               decimal,
			minuend:             []int{7, 2},
			subtrahend:          []int{4, 9},
			minuendRemainder:    0,
			subtrahendRemainder: 0,
			denominator:         0,
			result: longResult{
				result:    []int{2, 3},
				remainder: 0,
			},
		},
		{
			table:               decimal,
			minuend:             []int{7, 2},
			subtrahend:          []int{4, 9},
			minuendRemainder:    4,
			subtrahendRemainder: 2,
			denominator:         0,
			result: longResult{
				result:    []int{2, 3},
				remainder: 2,
			},
		},
		{
			table:               decimal,
			minuend:             []int{7, 2},
			subtrahend:          []int{4, 9},
			minuendRemainder:    4,
			subtrahendRemainder: 2,
			denominator:         3,
			result: longResult{
				result:    []int{2, 3},
				remainder: 2,
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(
			"minuend: %+v, subtrahend: %+v, minuendRemainder: %d, subtrahendRemainder: %d, denominator: %d",
			tc.minuend, tc.subtrahend, tc.minuendRemainder, tc.subtrahendRemainder, tc.denominator),
			func(t *testing.T) {
				result, err := tc.table.longSubtract(
					tc.minuend,
					tc.subtrahend,
					tc.minuendRemainder,
					tc.subtrahendRemainder,
					tc.denominator)
				if err != nil {
					t.Fatalf("Unexpected error occurred: %s", err)
				}

				unexported := cmp.AllowUnexported(longResult{})
				diff := cmp.Diff(result, tc.result, unexported)
				if diff != "" {
					t.Errorf("Expected and actual differences different:\n%s", diff)
				}
			})
	}
}

func TestRightpad(t *testing.T) {
	type test struct {
		input  []int
		length int
		result []int
	}

	tests := []test{
		{input: []int{1, 2}, length: 4, result: []int{1, 2, 0, 0}},
		{input: []int{1, 2}, length: 3, result: []int{1, 2, 0}},
		{input: []int{1, 2}, length: 2, result: []int{1, 2}},
		{input: []int{1, 2}, length: 1, result: []int{1, 2}},
		{input: []int{1, 2}, length: 0, result: []int{1, 2}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("input: %+v, length: %d", tc.input, tc.length), func(t *testing.T) {
			result := rightpad(tc.input, tc.length)

			diff := cmp.Diff(result, tc.result)
			if diff != "" {
				t.Errorf("Expected and actual slices different:\n%s", diff)
			}
		})
	}
}

func TestLeftPad(t *testing.T) {
	type test struct {
		input  []int
		length int
		result []int
	}

	tests := []test{
		{input: []int{1, 2}, length: 4, result: []int{0, 0, 1, 2}},
		{input: []int{1, 2}, length: 3, result: []int{0, 1, 2}},
		{input: []int{1, 2}, length: 2, result: []int{1, 2}},
		{input: []int{1, 2}, length: 1, result: []int{1, 2}},
		{input: []int{1, 2}, length: 0, result: []int{1, 2}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("input: %+v, length: %d", tc.input, tc.length), func(t *testing.T) {
			result := leftpad(tc.input, tc.length)

			diff := cmp.Diff(result, tc.result)
			if diff != "" {
				t.Errorf("Expected and actual slices different:\n%s", diff)
			}
		})
	}
}

func TestLongLinspace(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table     SymbolTable
		start     []int
		end       []int
		number    int
		divisions int
		result    []longResult
	}

	tests := []test{
		{
			table:     decimal,
			start:     []int{0, 0},
			end:       []int{1, 0},
			number:    1,
			divisions: 2,
			result: []longResult{
				{result: []int{0, 5}},
			},
		},
		{
			table:     decimal,
			start:     []int{0, 0},
			end:       []int{2, 0},
			number:    3,
			divisions: 4,
			result: []longResult{
				{result: []int{0, 5}},
				{result: []int{1, 0}},
				{result: []int{1, 5}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(
			"start: %+v, end: %+v, number: %d, divisions:: %d",
			tc.start, tc.end, tc.number, tc.divisions),
			func(t *testing.T) {
				result, err := tc.table.longLinspace(tc.start, tc.end, tc.number, tc.divisions)
				if err != nil {
					t.Fatalf("Unexpected error occurred: %s", err)
				}

				unexported := cmp.AllowUnexported(longResult{})
				diff := cmp.Diff(result, tc.result, unexported)
				if diff != "" {
					t.Errorf("Expected and actual linearly spaced numbers different:\n%s", diff)
				}
			})
	}
}

func TestRoundFraction(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table       SymbolTable
		numerator   int
		denominator int
		rounded     []int
	}

	tests := []test{
		{
			table:       decimal,
			numerator:   3,
			denominator: 4,
			rounded:     []int{8},
		},
		{
			table:       decimal,
			numerator:   4,
			denominator: 3,
			rounded:     []int{1, 3},
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf(
			"numerator: %d, denominator: %d", tc.numerator, tc.denominator), func(t *testing.T) {
			rounded := tc.table.roundFraction(tc.numerator, tc.denominator)

			diff := cmp.Diff(rounded, tc.rounded)
			if diff != "" {
				t.Errorf("Expected and actual rounded numbers different:\n%s", diff)
			}
		})
	}
}

func TestChopDigits(t *testing.T) {
	type test struct {
		current  []int
		previous []int
		chopped  []int
	}

	tests := []test{
		{current: []int{2, 9}, previous: []int{1, 9}, chopped: []int{2}},
		{current: []int{1, 2}, previous: []int{1, 3}, chopped: []int{1, 2}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", t), func(t *testing.T) {
			chopped := chopDigits(tc.current, tc.previous, 0)

			diff := cmp.Diff(chopped, tc.chopped)
			if diff != "" {
				t.Errorf("Expected and actual chopped digits different:\n%s", diff)
			}
		})
	}
}

func TestChopSuccessiveDigits(t *testing.T) {
	type test struct {
		startDigits []int
		digits      [][]int
		endDigits   []int
		chopped     [][]int
	}

	tests := []test{
		{
			startDigits: []int{0},
			digits: [][]int{
				{2, 4, 9},
				{4, 9, 9},
				{7, 4, 9},
			},
			endDigits: []int{1, 0},
			chopped: [][]int{
				{2},
				{4},
				{7},
			},
		},
		{
			startDigits: []int{0},
			digits: [][]int{
				{2},
				{4},
				{7},
			},
			endDigits: []int{1, 0},
			chopped: [][]int{
				{2},
				{4},
				{7},
			},
		},
	}

	for _, tc := range tests {
		t.Run(
			fmt.Sprintf("start digits: %+v, digits: %+v, end digits: %+v",
				tc.startDigits, tc.digits, tc.endDigits),
			func(t *testing.T) {
				chopped := chopSuccessiveDigits(tc.startDigits, tc.digits, tc.endDigits, 0)

				diff := cmp.Diff(chopped, tc.chopped)
				if diff != "" {
					t.Errorf("Expected and actual chopped digits different:\n%s", diff)
				}
			})
	}
}

func TestBetween(t *testing.T) {
	decimal, err := MakeSymbolTable(Decimal)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	type test struct {
		table        SymbolTable
		start        string
		end          string
		number       int
		divisions    int
		placesToKeep int
		results      []string
	}

	tests := []test{
		{
			table:        decimal,
			start:        "",
			end:          "",
			number:       3,
			divisions:    0,
			placesToKeep: 0,
			results:      []string{"2", "4", "7"},
		},
		{
			table:        decimal,
			start:        "2",
			end:          "3",
			number:       3,
			divisions:    0,
			placesToKeep: 0,
			results:      []string{"23", "25", "28"},
		},
		{
			table:        base62,
			start:        "",
			end:          "",
			number:       3,
			divisions:    0,
			placesToKeep: 0,
			results:      []string{"F", "U", "k"},
		},
		{
			table:        base62,
			start:        "aV",
			end:          "b",
			number:       3,
			divisions:    0,
			placesToKeep: 0,
			results:      []string{"ac", "ak", "as"},
		},
		{
			table:        base62,
			start:        "",
			end:          "",
			number:       10,
			divisions:    0,
			placesToKeep: 0,
			results:      []string{"5", "B", "G", "M", "S", "X", "d", "j", "o", "u"},
		},
		{
			table:        base62,
			start:        "",
			end:          "",
			number:       1,
			divisions:    10000,
			placesToKeep: 4,
			results:      []string{"00Npd"},
		},
		{
			table:        base62,
			start:        "",
			end:          "",
			number:       1,
			divisions:    10000,
			placesToKeep: 0,
			results:      []string{"00N"},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(
			"start: %+v, end: %+v, number: %d, divisions: %d, placesToKeep: %d",
			tc.start, tc.end, tc.number, tc.divisions, tc.placesToKeep),
			func(t *testing.T) {
				got, err := tc.table.Between(tc.start, tc.end, tc.number, tc.divisions, tc.placesToKeep)
				if err != nil {
					t.Fatalf("Unexpected error occurred: %s", err)
				}

				diff := cmp.Diff(got, tc.results)
				if diff != "" {
					t.Errorf("Expected and actual results numbers different:\n%s", diff)
				}
			})
	}
}

// This is testing adding a bunch of things one at at time to the end.
func TestAddLotsOfItems(t *testing.T) {
	// A large number of divisions is critical in order to keep the strings short.
	const divisions = 10000
	const placesToKeep = 4

	type item struct {
		index int
		order string
	}

	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	// Place the first item dead in the middle of the address pace.
	first, err := base62.Between("", "", 1, divisions, placesToKeep)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	var items []item
	items = append(items, item{index: 0, order: first[0]})

	for i := 1; i < 10000; i++ {
		curr, err := base62.Between(items[len(items)-1].order, "", 1, divisions, placesToKeep)
		if err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		items = append(items, item{index: i, order: curr[0]})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].order < items[j].order
	})

	got := lo.Map(items, func(x item, _ int) int {
		return x.index
	})

	var want []int
	for i := 0; i < 10000; i++ {
		want = append(want, i)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual sorted items different:\n%s", diff)
	}
}

// This is testing moving two items around a lot of times.
func TestMoveItemLotsOfTimes(t *testing.T) {
	const divisions = 10000
	const placesToKeep = 4

	type item struct {
		index int
		order string
	}

	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	initial, err := base62.Between("", "", 4, divisions, placesToKeep)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	items := []item{
		{index: 0, order: initial[0]},
		{index: 1, order: initial[1]},
		{index: 2, order: initial[2]},
		{index: 3, order: initial[3]},
	}

	for i := 0; i < 1000; i++ {
		curr, err := base62.Between(items[0].order, items[1].order, 1, 0, 0)
		if err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		items[2].order = curr[0]

		sort.Slice(items, func(i, j int) bool {
			return items[i].order < items[j].order
		})
	}

	got := lo.Map(items, func(x item, _ int) int {
		return x.index
	})

	want := []int{0, 1, 2, 3}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual sorted items different:\n%s", diff)
	}
}
