package memorydb

import "gostorage3081/pkg/storage/interface"

type DB []storage.Task

func (db DB) Tasks (int, int) ([]storage.Task, error ){
	return db, nil
}

func (db DB) NewTask( storage.Task) (int, error){
	return 1, nil
}

func (db DB) EditTask (storage.Task) (int, error){
	return 1,nil
}

func (db DB) DeleteTask(int) (int, error){
	return 1, nil
}