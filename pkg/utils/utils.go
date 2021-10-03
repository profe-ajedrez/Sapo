package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/profe-ajedrez/sapo/pkg/bd/drivers"
	"github.com/profe-ajedrez/sapo/pkg/errors"
	"github.com/profe-ajedrez/sapo/pkg/morestring"
)

func GracefullyStop(err error) {
	json := errors.ToJson(err)
	fmt.Println(json)
}

func CheckParams(file string, daTable string, maxConns int, maxIdles int, extended bool, columnToSearch string) error {

	if len(file) == 0 {
		return errors.ErrorEmptyFileName()
	}
	if !FileExists(file) {
		return errors.ErrorNoFile()
	}
	if len(daTable) == 0 {
		return errors.ErrorNoTable()
	}

	return nil
}

func GetData(file string, daTable string, maxConns int, maxIdles int, extended bool, columnToSearch string, structure bool) (string, error) {
	drivers.ExtendedQuery = extended
	drivers.ColumnToSearch = columnToSearch

	db := drivers.Maria{}
	db.Connect(file, maxConns, maxIdles, time.Minute*4)
	defer db.Close()

	rels, err := db.Relations(daTable)
	if err == nil {
		json, err := rels.ToPrettyJson()
		if err == nil {

			if structure {
				daStruct, _ := db.GetStructure(daTable)
				json = morestring.ReplaceLast(json, "}", ",\n    \"structure\" : "+daStruct+"\n}")
			}
			return json, nil
		}
	}
	return morestring.Null, err
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
