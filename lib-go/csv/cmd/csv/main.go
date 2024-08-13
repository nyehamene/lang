package main

import (
	"fmt"
	"os"

	"github.com/nyehamene/lang/csv/parser"
	"github.com/nyehamene/lang/csv/tokenizer"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	filename := os.Args[1]
	source, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	t := tokenizer.New(string(source))
	p := parser.New(*t)

	records, err := p.ParseAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	summarize(records)
}

func usage() {
	fmt.Println("Usage: csvgo file")
	os.Exit(1)
}

func summarize(records []parser.Record) {
	countRecords := func() {
		fmt.Println("- Records:", len(records))
	}

	countFields := func() {
		group := map[int]int{}

		for _, v := range records {
			cols := len(v)
			group[cols]++
		}

		if len(group) == 1 {
			for k := range group {
				fmt.Println("- Fields: ", k)
			}
			return
		}

		fmt.Printf("- Fields:")

		for k, v := range group {
			fmt.Printf("- %d records have %d fields\n", v, k)
		}
	}

	fmt.Printf("# Summary\n\n")
	countRecords()
	countFields()
}
