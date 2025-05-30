package main

import (
	"encoding/json"
	"log"
	"todo/db"

	"github.com/google/uuid"

	// "todo/db"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func main() {

	db.Init()

	// create a new todo
	http.HandleFunc("/create", enableCORS(handleCreateTodo))

	// get all todos
	http.HandleFunc("/get-all-todos", enableCORS(handleGetAllTodos))

	// update a todo
	http.HandleFunc("/update", enableCORS(handleUpdate))

	// delete a todo
	http.HandleFunc("/delete", enableCORS(handleDelete))

	log.Println("Starting server on :8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	// 1. 读取前端传来的参数
	params := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 3. 处理参数
	name := params["name"]
	description := params["description"]
	id := uuid.New().String()
	// 4. 构造数据
	var newTodo db.Todo = db.Todo{
		ID:          id,
		Name:        name,
		Description: description,
		Completed:   false,
	}

	// 5. 将数据存入数据库
	err = db.CreateTodo(newTodo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 6. 返回结果
	w.WriteHeader(http.StatusOK)
}

func handleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// log.Println("handleGetAllTodos:", db.Todos)
	// 2. 读取数据库中的数据
	todos, err := db.GetAllTodos()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO
	params := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := params["id"]
	name := params["name"]
	description := params["description"]
	completed := params["completed"]

	err = db.UpdateTodo(db.Todo{
		ID:          id,
		Name:        name,
		Description: description,
		Completed:   completed == "true",
	})
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := params["id"]

	if err = db.DeleteTodo(id); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 返回结果
	w.WriteHeader(http.StatusOK)
}
