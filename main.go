package main

import (
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "strconv"
   "github.com/gorilla/mux"
)

// Pessoa representa uma pessoa com ID e Nome
type Pessoa struct {
   ID   int    `json:"id"`
   Nome string `json:"nome"`
}

// Slice para armazenar as pessoas
var pessoas []Pessoa
var nextID int = 1

// Função para retornar todas as pessoas
func getPessoas(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(pessoas)
}

// Função para retornar uma pessoa por ID
func getPessoa(w http.ResponseWriter, r *http.Request) {
   params := mux.Vars(r)
   id, err := strconv.Atoi(params["id"])
   if err != nil {
      http.Error(w, "ID inválido", http.StatusBadRequest)
      return
   }

   for _, pessoa := range pessoas {
      if pessoa.ID == id {
         w.Header().Set("Content-Type", "application/json")
         json.NewEncoder(w).Encode(pessoa)
         return
      }
   }

   http.Error(w, "Pessoa não encontrada", http.StatusNotFound)
}

// Função para adicionar uma nova pessoa
func createPessoa(w http.ResponseWriter, r *http.Request) {
   var pessoa Pessoa
   if err := json.NewDecoder(r.Body).Decode(&pessoa); err != nil {
      http.Error(w, "Dados inválidos", http.StatusBadRequest)
      return
   }
   pessoa.ID = nextID
   nextID++

   pessoas = append(pessoas, pessoa)

   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(pessoa)
}

// Função para remover uma pessoa por ID
func deletePessoa(w http.ResponseWriter, r *http.Request) {
   params := mux.Vars(r)
   id, err := strconv.Atoi(params["id"])
   if err != nil {
      http.Error(w, "ID inválido", http.StatusBadRequest)
      return
   }

   for i, pessoa := range pessoas {
      if pessoa.ID == id {
         pessoas = append(pessoas[:i], pessoas[i+1:]...)
         w.WriteHeader(http.StatusNoContent)
         return
      }
   }

   http.Error(w, "Pessoa não encontrada", http.StatusNotFound)
}

func main() {
   r := mux.NewRouter()

   // Rotas
   r.HandleFunc("/pessoa", getPessoas).Methods("GET")
   r.HandleFunc("/pessoa/{id}", getPessoa).Methods("GET")
   r.HandleFunc("/pessoa", createPessoa).Methods("POST")
   r.HandleFunc("/pessoa/{id}", deletePessoa).Methods("DELETE")

   fmt.Println("Servidor iniciado na porta 8080")
   log.Fatal(http.ListenAndServe(":8080", r))
}