package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	// Usage go run . -f <filepath> to load existing data from a file
	var data map[string]Vector
	var embedding_size int

	if len(os.Args) > 2 {
		if os.Args[1] == "-f" {
			file, err := os.Open(os.Args[2])
			if err != nil {
				fmt.Println("Error: Reading previous db exiting ", err)
				return
			}
			defer file.Close()

			decoder := json.NewDecoder(file)

			if err := decoder.Decode(&data); err != nil {
				fmt.Println("Error: Parsing the json doc", err)
				return
			}

			for _, v := range data {
				embedding_size = len(v)
				break
			}

		}
	} else {
		RunBenchmark()
		return
	}

	if embedding_size == 0 {
		embedding_size = 128
	}

	db := NewVecMemDB("sample_db", embedding_size)

	if len(data) > 0 {
		db.SetDBData(data)
	}

	k1 := db.Insert(Vector{3.0, 2.0, 1.0})
	fmt.Println(k1)
	k2 := db.Insert(Vector{1.0, 1.0, 2.0})
	fmt.Println(k2)

	StartServer(db)

}
