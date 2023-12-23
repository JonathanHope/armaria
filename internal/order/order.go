// order includes helpers to allow users to manually order lists of things efficiently.
// Allowing users to sort a list of things manually is harder than it seems.
// You could increment the order of everything after the item that was moved.
// But then how many items will you have to modify when something moves?
// A common tactic is to use floating point numbers.
// For instance to add an item x after an item y you could do y.pos =  x.pos + x.pos / 2.
// The issue is that you will run out of precision at some point.
// Another less commonly used tactics is strings.
// They can be lexicographical sorted, and effectively have infinite precision.
// Someone on Stackoverflow designed an elegant algorithm that can be used for this:
// https://stackoverflow.com/questions/38923376
// Someone else wrote a library that fleshed out the idea more:
// https://github.com/fasiha/mudderjs
// What follows is not my invention, merely a translation to Go.
package order

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/samber/lo"
)

// Decimal is the symbols for a decimal number system.
const Decimal = "0123456789"

// Hex is the symbols for a hexadecimal number system.
const Hex = "0123456789abcdef"

// Base62 is the symbols for a Base62 number system.
const Base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// ErrSymbolsMustBeUnique is returned when duplicate symbols are provided for a numeric system.
var ErrSymbolsMustBeUnique = errors.New("symbols must be unique")

// ErrBetweenOrderWrong is returned when b < a for a call to Between.
var ErrBetweenOrderWrong = errors.New("b cannot be less than a")

// SymbolTable is the primary underlying data structure for this algorithm.
// Effectively it lets you define a numeric system.
// You can then do some conversions with it: positive integers <-> digits <-> strings.
// The README for mudder has some fun examples of numeric systems.
// There are some simplification made here:
// - It only supports building a symbol table out of a string.
// - The base is not configurable (it is the number of symbols).
// - The symbols must be unique (this means the symbol table always has the prefix property).
type SymbolTable struct {
	numberToSymbol []string       // lookup table for number -> symbol
	symbolToNumber map[string]int // lookup table for symbol -> number
	base           int            // the base of the numeric system
}

// MakeSymbolTable builds a symbol table from  a string of symbols.
// Each character in the string will become a symbol.
func MakeSymbolTable(symbols string) (SymbolTable, error) {
	table := SymbolTable{}
	table.numberToSymbol = strings.Split(symbols, "")

	if len(table.numberToSymbol) != len(lo.Uniq(table.numberToSymbol)) {
		return SymbolTable{}, ErrSymbolsMustBeUnique
	}

	table.symbolToNumber = map[string]int{}
	for i, str := range table.numberToSymbol {
		table.symbolToNumber[str] = i
	}

	table.base = len(table.numberToSymbol)

	return table, nil
}

// numberToDigits converts a number to a digit.
func (table SymbolTable) numberToDigits(number int) []int {
	digits := make([]int, 0)
	for number >= 1 {
		digits = append(digits, number%table.base)
		number = number / table.base
	}

	if len(digits) > 0 {
		return lo.Reverse(digits)
	}

	return []int{0}
}

// digitsToString converts digits to a string.
func (table SymbolTable) digitsToString(digits []int) (string, error) {
	for _, digit := range digits {
		if digit > len(table.numberToSymbol)-1 {
			return "", errors.New("cannot map between symbol and number")
		}
	}

	strs := lo.Map(digits, func(x int, _ int) string {
		return table.numberToSymbol[x]
	})

	return strings.Join(strs, ""), nil
}

// stringToDigits converts a string to digits.
func (table SymbolTable) stringToDigits(str string) ([]int, error) {
	tokens := strings.Split(str, "")

	for _, token := range tokens {
		_, ok := table.symbolToNumber[token]
		if !ok {
			return nil, errors.New("cannot map between symbol and number")
		}
	}

	digits := lo.Map(tokens, func(x string, _ int) int {
		return table.symbolToNumber[x]
	})

	return digits, nil
}

// digitsToNumber converts digits to a number.
func (table SymbolTable) digitsToNumber(digits []int) int {
	currBase := 1
	result := lo.ReduceRight(digits, func(accum int, curr int, _ int) int {
		ret := accum + curr*currBase
		currBase *= table.base
		return ret
	}, 0)

	return result
}

// NumberToString converts a number to a string.
func (table SymbolTable) NumberToString(number int) (string, error) {
	digits := table.numberToDigits(number)

	str, err := table.digitsToString(digits)
	if err != nil {
		return "", fmt.Errorf("error converting number to string: %w", err)
	}

	return str, nil
}

// StringToNumber converts a string to a number.
func (table SymbolTable) StringToNumber(str string) (int, error) {
	digits, err := table.stringToDigits(str)
	if err != nil {
		return 0, fmt.Errorf("error converting string to a number: %w", err)
	}

	num := table.digitsToNumber(digits)

	return num, nil
}

// rightpad right pads an input array with zeros to the desired length.
func rightpad(input []int, length int) []int {
	if len(input) < length {
		for len(input) < length {
			input = append(input, 0)
		}
	}

	return input
}

// leftpad left pads an input array with zeroes to the desired length.
func leftpad(input []int, length int) []int {
	reversed := lo.Reverse(input)
	rightpadded := rightpad(reversed, length)

	return lo.Reverse(rightpadded)
}

// roundFraction rounds a fraction to a decimal.
func (table SymbolTable) roundFraction(numerator int, denominator int) []int {
	places := math.Ceil(math.Log(float64(denominator)) / math.Log(float64(table.base)))
	scale := math.Pow(float64(table.base), places)
	scaled := math.Round(float64(numerator) / float64(denominator) * scale)
	digits := table.numberToDigits(int(scaled))

	return leftpad(digits, int(places))
}

// chopDigits returns the digits in a up to (and including) the first digit that doesn't match b.
func chopDigits(current []int, previous []int, placesToKeep int) []int {
	get := func(digits []int, idx int) int {
		if idx < len(digits) {
			return digits[idx]
		}
		return -1
	}

	for idx := placesToKeep; idx < len(current); idx++ {
		if get(current, idx) != -1 && get(current, idx) != 0 && get(previous, idx) != get(current, idx) {
			return current[:idx+1]
		}
	}

	return current
}

// chopSuccessiveDigits removes as many digits as it can while maintaining the lexicographic order.
func chopSuccessiveDigits(startDigits []int, digits [][]int, endDigits []int, placesToKeep int) [][]int {
	digits = append([][]int{startDigits}, digits...)
	digits = append(digits, endDigits)

	chopped := lo.Reduce(digits[1:], func(agg [][]int, item []int, _ int) [][]int {
		agg = append(agg, chopDigits(item, agg[len(agg)-1], placesToKeep))
		return agg
	}, [][]int{digits[0]})

	return chopped[1 : len(digits)-1]
}

// longResult is a result struct for the long* functions.
type longResult struct {
	result    []int // result of the operation
	remainder int   // remainder from the operation
	carry     bool  // if true there is a carry digit that needs to be accounted for in the next operation
}

// longDiv divides a set of digits by a divisor.
// This is standard long division without any caveats.
func (table SymbolTable) longDiv(dividend []int, divisor int) longResult {
	var remainder int
	var quotient []int
	for _, digit := range dividend {
		r := digit + remainder*table.base
		quotient = append(quotient, r/divisor)
		remainder = r % divisor
	}

	return longResult{
		result:    quotient,
		remainder: remainder,
	}
}

// longAdd adds two sets of digits together.
// It is assumed that the digits have radix point before the first digit.
// For example [1, 2] + [4, 5, 6] => 0.12 + 0.456 = 0.576
// If carry is true it can be assumed there should be a one to the left of the radix point.
func (table SymbolTable) longAdd(addend1 []int, addend2 []int, remainder int, denominator int) (longResult, error) {
	if len(addend1) != len(addend2) {
		return longResult{}, errors.New("addends must be same length")
	}

	carry := remainder >= denominator
	if carry {
		remainder -= denominator
	}

	sum := make([]int, len(addend2))

	for i := len(addend1) - 1; i >= 0; i-- {
		s := addend1[i] + addend2[i]
		if carry {
			s += 1
		}

		carry = s >= table.base
		if carry {
			sum[i] = s - table.base
		} else {
			sum[i] = s
		}
	}

	return longResult{
		result:    sum,
		remainder: remainder,
		carry:     carry,
	}, nil
}

// longSubtract does subtraction on two sets of digits.
// The code here is a little busy and should be refactored at some point.
func (table SymbolTable) longSubtract(minuend []int, subtrahend []int, minuendRemainder int, subtrahendRemainder int, denominator int) (longResult, error) {
	if len(minuend) != len(subtrahend) {
		return longResult{}, errors.New("minuend/subtrahend must be same length")
	}

	// Clone the minuend and subtrahend.
	mc := make([]int, len(minuend))
	copy(mc, minuend)
	sc := make([]int, len(subtrahend))
	copy(sc, subtrahend)

	// If there are remainders append them to the minuend/subtrahend.
	hasRemainder := minuendRemainder != 0 || subtrahendRemainder != 0
	if hasRemainder {
		mc = append(mc, minuendRemainder)    //nolint:makezero
		sc = append(sc, subtrahendRemainder) //nolint:makezero
	}

	difference := make([]int, len(mc))

OUTER:
	for i := len(mc) - 1; i >= 0; i-- {

		if mc[i] >= sc[i] {
			difference[i] = mc[i] - sc[i]
			continue
		}

		for j := i - 1; j >= 0; j-- {
			if mc[j] > 0 {
				// found a non-zero digit. Decrement it
				mc[j]--
				// increment digits to its right by `base-1`
				for k := j + 1; k < i; k++ {
					mc[k] += table.base - 1
				}
				// until you reach the digit you couldn't subtract
				if hasRemainder && i == len(mc)-1 {
					difference[i] = mc[i] + denominator - sc[i]
				} else {
					difference[i] = mc[i] + table.base - sc[i]
				}
				continue OUTER
			}
		}

		return longResult{}, errors.New("failed to find digit to borrow from")
	}

	if hasRemainder {
		return longResult{
			result:    difference[:len(difference)-1],
			remainder: difference[len(difference)-1],
		}, nil
	}

	return longResult{
		result:    difference,
		remainder: 0,
	}, nil
}

// longLinspace returns number linearly spaced numbers between start and end.
// More specifically it is the equation (a + (b-a)/M*n) for n=[1, 2, ..., N], where N<M
func (table SymbolTable) longLinspace(start []int, end []int, number int, divisions int) ([]longResult, error) {
	if len(start) < len(end) {
		start = rightpad(start, len(end))
	} else if len(end) < len(start) {
		end = rightpad(end, len(start))
	}

	if slices.Equal(start, end) {
		return nil, errors.New("start and end strings lexicographically inseparable")
	}

	aDiv := table.longDiv(start, divisions)
	bDiv := table.longDiv(end, divisions)

	aPrev, err := table.longSubtract(start, aDiv.result, 0, aDiv.remainder, divisions)
	bPrev := bDiv
	if err != nil {
		return nil, fmt.Errorf("error finding linearly spaced numbers: %w", err)
	}

	var results []longResult
	for n := 1; n <= number; n++ {
		result, err := table.longAdd(aPrev.result, bPrev.result, aPrev.remainder+bPrev.remainder, divisions)
		if err != nil {
			return nil, fmt.Errorf("error finding linearly spaced r: %w", err)
		}
		results = append(results, result)

		aPrev, err = table.longSubtract(aPrev.result, aDiv.result, aPrev.remainder, aDiv.remainder, divisions)
		if err != nil {
			return nil, fmt.Errorf("error finding linearly spaced numbers: %w", err)
		}

		bPrev, err = table.longAdd(bPrev.result, bDiv.result, bPrev.remainder+bDiv.remainder, divisions)
		if err != nil {
			return nil, fmt.Errorf("error finding linearly spaced numbers: %w", err)
		}
	}

	return results, nil
}

// Between returns number strings between start and end.
// These strings should be evenly spaced in the address space between start and end.
// For example if you request 3 numbers between 0 and 10 in base10 you would get 3, 4, and 7.
// All parameters to this are optional:
// - start defaults to the first symbol in the symbol table.
// - end defaults to the last symbol in the symbol table repeated a number of times.
// - number defaults to 1.
// - divisions defaults to number + 1.
// This only works ascending so start must be > end.
func (table SymbolTable) Between(start string, end string, number int, divisions int, placesToKeep int) ([]string, error) {
	if start == "" {
		start = table.numberToSymbol[0]
	}

	if end == "" {
		end = strings.Repeat(table.numberToSymbol[len(table.numberToSymbol)-1], len(start)+6)
	}

	if number == 0 {
		number = 1
	}

	if divisions == 0 {
		divisions = number + 1
	}

	if end < start {
		return nil, ErrBetweenOrderWrong
	}

	startDigits, err := table.stringToDigits(start)
	if err != nil {
		return nil, fmt.Errorf("error finding string between strings: %w", err)
	}

	endDigits, err := table.stringToDigits(end)
	if err != nil {
		return nil, fmt.Errorf("error finding string between strings: %w", err)
	}

	intermediateDigits, err := table.longLinspace(startDigits, endDigits, number, divisions)
	if err != nil {
		return nil, fmt.Errorf("error finding string between strings: %w", err)
	}

	finalDigits := lo.Map(intermediateDigits, func(lr longResult, index int) []int {
		return append(lr.result, table.roundFraction(lr.remainder, divisions)...)
	})

	var results []string
	for _, digits := range chopSuccessiveDigits(startDigits, finalDigits, endDigits, placesToKeep) {
		result, err := table.digitsToString(digits)
		if err != nil {
			return nil, fmt.Errorf("error finding string between strings: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}
