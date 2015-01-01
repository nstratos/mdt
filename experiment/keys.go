package main

import (
	"fmt"
	"strconv"
)

func main() {
	// keys := []rune{'q', 'w', 'e', 'a', 's', 'd'}
	/*for k := range keys {
		fmt.Printf("%s = %v\n", strconv.QuoteRune(k), k)
	}*/
	fmt.Printf("%s = %v\n", strconv.QuoteRune('a'), 'a')
	fmt.Printf("%s = %v\n", strconv.QuoteRune('d'), 'd')
	fmt.Printf("%s = %v\n", strconv.QuoteRune('e'), 'e')
	fmt.Printf("%s = %v\n", strconv.QuoteRune('q'), 'q')
	fmt.Printf("%s = %v\n", strconv.QuoteRune('s'), 's')
	fmt.Printf("%s = %v\n", strconv.QuoteRune('w'), 'w')
}
