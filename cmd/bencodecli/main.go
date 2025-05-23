package main

import (
	"fmt"
	"BencodeParser/pkg/bencode"
)

func main() {
	// ---- Decoding Bencode ----
	fmt.Println("The below are decoding bencode")

	// Parsing integer
	intData := "i3e"
	parsedInt, err := bencode.ParseString(intData)
	if err != nil {
		fmt.Printf("Error parsing integer '%s'': %v\n", intData, err)
	} else {
		fmt.Printf("The parsed integer from '%s' : %v (Type: %T)\n", intData, parsedInt, parsedInt)
	}
	intDataNeg := "i-7e"
	parsedIntNeg, err := bencode.ParseString(intDataNeg)
	if err != nil {
		fmt.Printf("Error parsing integer '%s': %v\n", intDataNeg, err)
	} else {
		fmt.Printf("Parsed integer from '%s': %v (Type: %T)\n", intDataNeg, parsedIntNeg, parsedIntNeg)
	}

	// Leading integer examples
	invalidIntData := "i07e"
	_, err = bencode.ParseString(invalidIntData)
	fmt.Printf("Parsing invalid integer '%s', expected error: %v\n", invalidIntData, err)

	// Parsing String
	stringData := "5:hello"
	parsedString, err := bencode.ParseString(stringData)
	if err != nil {
		fmt.Printf("Error parsing string '%s': %v\n", stringData, err)
	} else {
		fmt.Printf("Parsed string from '%s': \"%s\" (Type: %T)\n", stringData, parsedString, parsedString)
	}

	emptyStringData := "0:"
	parsedEmptyString, err := bencode.ParseString(emptyStringData)
	if err != nil {
		fmt.Printf("Error parsing empty string '%s': %v\n", emptyStringData, err)
	} else {
		fmt.Printf("Parsed empty string from '%s': \"%s\" (Type: %T)\n", emptyStringData, parsedEmptyString, parsedEmptyString)
	}

	// Parsing List
	listData := "li10e4:spam3:eggse" // [10, "spam", "eggs"]
	parsedList, err := bencode.ParseString(listData)
	if err != nil {
		fmt.Printf("Error parsing list '%s': %v\n", listData, err)
	} else {
		fmt.Printf("Parsed list from '%s': %v (Type: %T)\n", listData, parsedList, parsedList)
	}

	// Parsing Dictionary
	dictData := "d3:cow3:moo4:spam4:eggse" // {"cow": "moo", "spam": "eggs"}
	parsedDict, err := bencode.ParseString(dictData)
	if err != nil {
		fmt.Printf("Error parsing dictionary '%s': %v\n", dictData, err)
	} else {
		fmt.Printf("Parsed dictionary from '%s': %v (Type: %T)\n", dictData, parsedDict, parsedDict)
	}

}
