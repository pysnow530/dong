package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type ListResponse struct {
	Data []Paper `json:"data"`
	Code int     `json:"code"`
}

type DetailResponse struct {
	Data Paper `json:"data"`
	Code int   `json:"code"`
}

type MessageResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Secreter struct {
	Token  string
	Secret []byte
}

func NewSecreter(r *http.Request) *Secreter {
	var secreter Secreter

	secreter.Token = r.Header.Get("ABRACADABRA-TOKEN")
	secreter.Secret = CalcSecret(secreter.Token, *FlagSalt)

	return &secreter
}

func (secreter *Secreter) CalcId(idx int) string {
	return CalcId(secreter.Token, idx)
}

// MessageListHandler
func MessageListHandler(w http.ResponseWriter, r *http.Request) {
	secreter := NewSecreter(r)

	if r.Method == "GET" {
		var papers []Paper

		for idx := 0; idx < 100; idx += 1 {
			var paper Paper
			id := secreter.CalcId(idx)
			DB.Unscoped().Find(&paper, "id = ?", id)
			if paper.ID == "" { // reach the tail
				break
			}

			if paper.DeletedAt == nil {
				err := paper.Decrypt(secreter.Secret)
				if err != nil {
					fmt.Fprintf(w, "Decrypt error: %s\n", err.Error())
					return
				}
				papers = append([]Paper{paper}, papers...)
			}
		}

		buf, err := json.Marshal(ListResponse{Data: papers, Code: 0})
		if err != nil {
			fmt.Fprintf(w, "Serialize error: %s\n", err.Error())
			return
		}
		w.Write(buf)

		return

	} else if r.Method == "POST" {
		var params map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			fmt.Fprintf(w, "Decode param error: %s\n", err.Error())
			return
		}

		argData, exist := params["data"]
		if !exist {
			fmt.Fprintf(w, "Param required: data")
			return
		}

		for idx := 0; idx < 100; idx += 1 {
			var paper Paper
			id := secreter.CalcId(idx)
			DB.Unscoped().Find(&paper, "id = ?", id)
			if paper.ID != "" { // keep going to the tail
				continue
			}

			paper.ID = id
			paper.Data = []byte(argData.(string))
			paper.Encrypt(secreter.Secret)
			DB.Create(&paper)

			paper.Decrypt(secreter.Secret)
			buf, err := json.Marshal(paper)
			if err != nil {
				fmt.Fprintln(w, "Serialize error: %s\n", err.Error())
				return
			}

			w.Write(buf)

			return
		}

	}
}

func MessageDetailHandler(w http.ResponseWriter, r *http.Request) {
	secreter := NewSecreter(r)

	vars := mux.Vars(r)
	id := vars["id"]

	var paper Paper
	DB.Find(&paper, "id = ?", id)
	paper.Decrypt(secreter.Secret)

	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode(paper)
		if err != nil {
			fmt.Fprintf(w, "Serialize error: %s\n", err.Error())
		}
		return

	} else if r.Method == "PUT" {
		var params map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			fmt.Fprintf(w, "Decode param error: %s\n", err.Error())
			return
		}

		argData, exist := params["data"]
		if !exist {
			fmt.Fprintf(w, "Param required: data")
			return
		}

		paper.Data = []byte(argData.(string))
		paper.Encrypt(secreter.Secret)
		DB.Save(&paper)

		paper.Decrypt(secreter.Secret)
		err = json.NewEncoder(w).Encode(paper)
		if err != nil {
			fmt.Fprintf(w, "Serialize error: %s\n", err.Error())
		}

		return

	} else if r.Method == "DELETE" {
		DB.Delete(paper)
		fmt.Fprintln(w, "deleted")
		return
	}

}

func GetGlobalHandler() *mux.Router {
	var router = mux.NewRouter()
	router.HandleFunc("/api/v1/papers", MessageListHandler)
	router.HandleFunc("/api/v1/papers/{id}", MessageDetailHandler)
	router.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir("app"))))

	return router
}
