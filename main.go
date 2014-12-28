package main

import "fmt"

func main() {
	var key string
	fmt.Printf("Enter something:\n")
	fmt.Scanf("%s", &key)
	fmt.Printf("You have entered %s.\n", key)
}
