package main

import (
	"HttpCalculator/WorkWithSQL"
	"log"
	"net/http"
)

func registrate(w http.ResponseWriter, r *http.Request) {
	Registrate(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	Login(w, r)
}

func expressions(w http.ResponseWriter, r *http.Request) {
	Expressions(w, r)
}

func orchestrator(w http.ResponseWriter, r *http.Request) {
	Orchestrator(w, r)
}

func orchestratorReturn(w http.ResponseWriter, r *http.Request) {
	OrchestratorReturn(w, r)
}

func main() {
	db, err := WorkWithSQL.CreateBD()
	if err != nil {
		log.Println("Ошибка создания db:", err)
		return
	}
	if db == nil {
		log.Println("db == nil")
		return
	}

	UserRepo = WorkWithSQL.NewSQLiteUserRepository(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/register", registrate)
	mux.HandleFunc("/api/v1/login", login)
	mux.HandleFunc("/api/v1/expressions/", expressions)
	mux.HandleFunc("/api/v1/expressions", expressions)
	mux.HandleFunc("/api/v1/calculate", orchestrator)
	mux.HandleFunc("/internal/task", orchestratorReturn)

	log.Println("Сервер запущен на порту 9090...")
	log.Fatal(http.ListenAndServe(":9090", mux))
}
