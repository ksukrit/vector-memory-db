package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

// this file has mainly the implementation of the vector memory database and db specific functions

type InMemVectorDB struct {
	db_name                string
	data                   map[string]Vector
	embedding_size         int
	totalCommandsProcessed int
	mu                     sync.RWMutex
}

func NewVecMemDB(db_name string, embedding_size int) *InMemVectorDB {
	return &InMemVectorDB{
		db_name:                db_name,
		data:                   make(map[string]Vector),
		embedding_size:         embedding_size,
		totalCommandsProcessed: 0,
	}
}

func (db *InMemVectorDB) GetEmbeddingSize() int {
	return db.embedding_size
}

func (db *InMemVectorDB) Get(key string) Vector {
	db.mu.RLock()
	defer db.mu.RUnlock()
	db.totalCommandsProcessed++
	return db.data[key]
}

func (db *InMemVectorDB) Insert(value Vector) string {
	// generate a unique UUID for the vector
	uuid := uuid.New()

	uuidString := uuid.String()

	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[uuidString] = value
	db.totalCommandsProcessed++
	return uuidString
}

func (db *InMemVectorDB) InsertKey(key string, value Vector) string {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
	db.totalCommandsProcessed++
	return key
}

func (db *InMemVectorDB) Exists(key string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, ok := db.data[key]
	db.totalCommandsProcessed++
	return ok
}

func (db *InMemVectorDB) Update(key string, value Vector) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.totalCommandsProcessed++
	if db.Exists(key) {
		db.data[key] = value
		return true
	} else {
		return false
	}
}

func (db *InMemVectorDB) Delete(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.totalCommandsProcessed++
	if db.Exists(key) {
		delete(db.data, key)
		return true
	} else {
		return false
	}
}

func (db *InMemVectorDB) Size() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	db.totalCommandsProcessed++
	return len(db.data)
}

func (db *InMemVectorDB) DumpFile(filepath string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	db.totalCommandsProcessed++
	fmt.Println("Dumping data to file: ", filepath)
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	defer f.Close()

	jsonData, err := json.Marshal(db.GetDBData())
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return false
	}

	f.Write(jsonData)

	return true
}

func (db *InMemVectorDB) GetProcessedCommand() int {
	return db.totalCommandsProcessed
}

func (db *InMemVectorDB) GetDBName() string {
	return db.db_name
}

func (db *InMemVectorDB) IncreaseEmbeddingSize(newEmbeddingSize int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.embedding_size = newEmbeddingSize
	db.totalCommandsProcessed++
	for key, value := range db.data {
		db.data[key] = PadVector(value, newEmbeddingSize)
	}
}

func (db *InMemVectorDB) GetDBData() map[string]Vector {
	db.mu.RLock()
	defer db.mu.RUnlock()
	db.totalCommandsProcessed++
	return db.data
}

func (db *InMemVectorDB) SetDBData(data map[string]Vector) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.totalCommandsProcessed++
	db.data = data
}
