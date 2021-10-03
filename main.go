package main

import (
	"flag"
	"fmt"
	"os"

	mainUtils "github.com/profe-ajedrez/sapo/pkg/utils"
)

const VERSION = "0.0.1"

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s -file=<json config file> -table=<table name> [-extended=<bool> -column=<column name>] [-structure=<bool>] [-maxopenconns=<int>] [-maxidleconns=<int>]", os.Args[0])
	}

	file := flag.String("file", "", "Database config file")
	daTable := flag.String("table", "", "Table to check relations")
	maxConns := flag.Int("maxopenconns", 16, "Maximum open connections to database\nyour server can handle at any one time")
	maxIdles := flag.Int("maxidleconns", 16, "Maximum idle connections to database\nthat can be allocated from this pool at the same time")
	extended := flag.Bool("extended", false, "true to force search the columns names presence in other tables,\nfalse to use `information_schema.key_column_usage`")
	columnToSearch := flag.String("column", "", "Column name to extensive search. Valid only with -extended=true")
	structure := flag.Bool("structure", false, "Whether get the table structure")
	version := flag.Bool("version", false, "shows version")

	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	err := mainUtils.CheckParams(*file, *daTable, *maxConns, *maxIdles, *extended, *columnToSearch)

	if err == nil {
		data, err := mainUtils.GetData(*file, *daTable, *maxConns, *maxIdles, *extended, *columnToSearch, *structure)

		if err == nil {
			fmt.Println(data)
			return
		}
	}
	mainUtils.GracefullyStop(err)
}
