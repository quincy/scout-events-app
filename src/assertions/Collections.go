package assertions

import (
	"log"
	"reflect"
)

type Collection[T any] []T

// ContainsExactly returns true if the `actual` Collection contains exactly the same elements
// in the same order as the `expected`.
func (actual Collection[T]) ContainsExactly(expected []T) bool {
	if len(actual) != len(expected) {
		log.Printf("Expected length of %v but got %v", len(expected), len(actual))
		return false
	}

	for i, expected := range expected {
		if !reflect.DeepEqual(expected, actual[i]) {
			log.Printf(
				"Expected value of\n%v\nbut got\n%v", expected, actual[i])
			return false
		}
	}

	return true
}
