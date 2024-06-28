package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"firebase.google.com/go/db"
)

type Pessoa struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

var client *db.Client

func getPessoas(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var pessoas map[string]Pessoa
	if err := client.NewRef("pessoa").Get(ctx, &pessoas); err != nil {
		http.Error(w, "Erro ao buscar pessoas", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)
}

func getPessoa(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	params := mux.Vars(r)
	id := params["id"]

	var pessoa Pessoa
	if err := client.NewRef("pessoa/"+id).Get(ctx, &pessoa); err != nil {
		http.Error(w, "Pessoa não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoa)
}

func createPessoa(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var p Pessoa
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	ref, err := client.NewRef("pessoa").Push(ctx, nil)
	if err != nil {
		http.Error(w, "Erro ao adicionar pessoa", http.StatusInternalServerError)
		return
	}

	p.ID = ref.Key
	if err := ref.Set(ctx, p); err != nil {
		http.Error(w, "Erro ao adicionar pessoa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deletePessoa(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	params := mux.Vars(r)
	id := params["id"]

	if err := client.NewRef("pessoa/" + id).Delete(ctx); err != nil {
		http.Error(w, "Erro ao deletar pessoa", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://api-restiful-default-rtdb.firebaseio.com/",
	}
	opt := option.WithCredentialsFile("api-restiful-firebase-adminsdk-9g7ut-c6f9cd0b86.json")
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("Erro ao inicializar app Firebase: %v", err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalf("Erro ao inicializar Realtime Database: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/pessoa", getPessoas).Methods("GET")
	r.HandleFunc("/pessoa/{id}", getPessoa).Methods("GET")
	r.HandleFunc("/pessoa", createPessoa).Methods("POST")
	r.HandleFunc("/pessoa/{id}", deletePessoa).Methods("DELETE")

	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
