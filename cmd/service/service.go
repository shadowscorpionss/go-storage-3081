package main

import (
	"fmt"
	"gostorage3081/pkg/storage/interface"
	"gostorage3081/pkg/storage/memorydb"
	"gostorage3081/pkg/storage/postgres"
	"log"
	"os"
)

var db storage.DbInterface

func main() {
	var err error
	authInfo := os.Getenv("crend")
	if authInfo == "" {
		os.Exit(1)
	}

	postgreConnStr := fmt.Sprintf("postgresql://%s@localhost/db30_8_1", authInfo)

	db, err = postgresdb.New(postgreConnStr)
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := db.Tasks(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)

	db = memorydb.DB{}

	tasks, err = db.Tasks(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)

}
