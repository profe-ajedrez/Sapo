package file

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	json, err := ReadJSON("./testdata/good.json")

	if err == nil {
		fmt.Println(json)

		erro := WriteJSON(json, "./testdata/test2.json", PrettyPrint)
		if erro != nil {
			fmt.Println(erro)
			panic("shuata falló el WriteJson")
		}
	} else {
		fmt.Print(err)
		panic("shuata falló el ReadJson")
	}
}

func BenchmarkRead(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := ReadJSON("./testdata/good.json")

		if err != nil {
			fmt.Print(err)
		}
	}
}
