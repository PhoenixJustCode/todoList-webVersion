package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
)



// func filterRes(w http.ResponseWriter, result map[string]string){ 
// 	for m,c := range result {
// 		fmt.Fprintf(w, "Task:%s <-> Time: %s\n",m,c)
// 	}
// }


func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Parseform() err :%v", err), http.StatusBadRequest)
		return
	}

	task := r.FormValue("task")
	day := r.FormValue("day")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"task": task,
		"day":  day,
	})
}

// func helloUser(writer http.ResponseWriter, request *http.Request) {
// 	var greeting = "Hello user. Welcome to our Todolist app"
// 	fmt.Fprint(writer, greeting)
// }


func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	fmt.Printf("Starting server at port 8080 \n")
	fmt.Printf("Link : http://localhost:8080/ \n")
	
	http.HandleFunc("/form", formHandler)

	if err:= http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}


