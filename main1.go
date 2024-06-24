package main

import (
	"fmt"
	"strings"
)

func main() {
	str1 := "et-0/0/0"
	str2 := "et-0/0/0:1"
	suffix := "0/0/0"

	fmt.Println(strings.HasSuffix(RemoveChannalisation(str1), suffix)) // true
	fmt.Println(strings.HasSuffix(RemoveChannalisation(str2), suffix)) // true
}

// RemoveChannalisationOld removes the string after ":"
// It will take the input as "et-0/0/0:1" and returns output as et-0/0/0"
func RemoveChannalisationOld(str2 string) string {
	if strings.Contains(str2, ":") {
        str2 = strings.Split(str2, ":")[0]
	}
	return str2
}

// RemoveChannalisation removes the string after ":"
// It will take the input as "et-0/0/0:1" and returns output as -0/0/0"
func RemoveChannalisation(str2 string) string {
	if strings.Contains(str2, ":") {
        str2 = strings.Split(str2, ":")[0]
	}
	if strings.Contains(str2, "-") {
        str2 = strings.Split(str2, "-")[1]
	}
	fmt.Printf("%v\n", str2)
	return str2
}