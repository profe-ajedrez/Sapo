package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/profe-ajedrez/sapo/pkg/bd/drivers"
	"github.com/profe-ajedrez/sapo/pkg/morestring"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s -file=<json config file> -maxopenconns=<int> -maxidleconns=<int>", os.Args[0])
	}

	file := flag.String("file", "", "Database config file")
	maxConns := flag.Int("maxopenconns", 16, "Maximum open connections to database\nyour server can handle at any one time")
	maxIdles := flag.Int("maxidleconns", 16, "Maximum idle connections to database\nthat can be allocated from this pool at the same time")
	extended := flag.Bool("extended", false, "true to force search the columns names presence in other tables,\nfalse to use `information_schema.key_column_usage`")
	columnToSearch := flag.String("column", "", "Column name to extensive search. Valid only with -extended=true")

	flag.Parse()

	if fileExists(*file) {
		data, err := GetData(*file, *maxConns, *maxIdles, *extended, *columnToSearch)

		if err == nil {
			fmt.Println(data)
			return
		}
		panic(err)
	}

	fmt.Println("file " + *file + " doesnt exists")
}

func GetData(file string, maxConns int, maxIdles int, extended bool, columnToSearch string) (string, error) {
	drivers.ExtendedQuery = extended
	drivers.ColumnToSearch = columnToSearch

	db := drivers.Maria{}
	db.Connect(file, maxConns, maxIdles, time.Minute*4)
	defer db.Close()

	rels, err := db.Relations("cart_item")
	if err == nil {
		json, err := rels.ToPrettyJson()
		if err == nil {
			return json, nil
		}
	}
	return morestring.Null, err
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
