package main

import (
	"fmt"
	DB "github.com/connorvoisey/shgrid_api/pkg/db"
	"github.com/connorvoisey/shgrid_api/pkg/load"
	_ "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"sync"
)

const (
	batchLoops      = 1
	batches         = 5
	recordsPerBatch = 10_000
)

func main() {
	_, db, err := load.Init()
	panicOnError(err, "Failed to init")
	defer db.Close()

	for i := 0; i < batchLoops; i++ {
		var wg sync.WaitGroup
		for i := 0; i < batches; i++ {
			wg.Add(1)
			go DB.Seed(db, recordsPerBatch, recordsPerBatch, recordsPerBatch, &wg)
		}
		wg.Wait()
	}
}

func panicOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}
