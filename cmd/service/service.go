package main

import (
	"fmt"
	"gostorage3081/pkg/storage/interface"
	"gostorage3081/pkg/storage/memorydb"
	"gostorage3081/pkg/storage/postgres"
	"log"
	"os"
)


func main() {
	var err error
	var db storage.DbInterface
	authInfo := os.Getenv("crend")
	if authInfo == "" {
		os.Exit(1)
	}

	postgreConnStr := fmt.Sprintf("postgresql://%s@localhost/db30_8_1", authInfo)

	db, err = postgresdb.New(postgreConnStr)
	if err != nil {
		log.Fatal(err)
	}
	testTasks(db)
	db = memorydb.DB{}
	testTasks(db)
}


func testTasks (db storage.DbInterface){
	//all tasks
	fmt.Println("ALL")
	tasks, err := db.Tasks(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)

	//new
	fmt.Println("NEW")
	task:=storage.Task{ Content: "ntsk", Title: "ntitle"}
	task.ID,err = db.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(task)

	//update
	fmt.Println("UDPATE")
	var affectedRows int
	task.Content="new ntsk"
	affectedRows, err = db.EditTask(task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(affectedRows)
	
	//delete
	fmt.Println("DELETE")
	affectedRows, err=db.DeleteTask(task.ID)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(affectedRows)
}