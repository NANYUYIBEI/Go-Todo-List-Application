package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// var Todos []Todo

func Init() {
	// init db
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&Todo{})
}

func CreateTodo(todo Todo) error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	result := db.Create(&todo)
	return result.Error
}

func GetAllTodos() ([]Todo, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var todos []Todo
	result := db.Find(&todos)
	return todos, result.Error
}

func UpdateTodo(todo Todo) error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	result := db.Save(&todo)
	return result.Error
}

func DeleteTodo(id string) error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	// result := db.Delete(&Todo{}, id)
	result := db.Delete(&Todo{}, "id = ?", id)
	return result.Error
}
