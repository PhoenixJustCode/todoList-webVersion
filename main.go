package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)


var db *DB

func main() {
	var err error
	db, err = NewDB("host=127.0.0.1 port=5432 user=postgres dbname=todolist_db sslmode=disable")
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/tasks", tasksHandler)       // GET - список задач
	http.HandleFunc("/tasks/add", addTaskHandler) // POST - добавить
	http.HandleFunc("/tasks/delete", deleteTaskHandler)   // POST - удалить
	http.HandleFunc("/tasks/update", updateTaskHandler)   // POST - обновить

	fmt.Println("Starting server at :8080")
	fmt.Printf("Link : http://localhost:8080/ \n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	tasks, err := db.GetAllTasks()
	if err != nil {
		http.Error(w, "Failed to get tasks", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error", 400)
		return
	}
	task := r.FormValue("task")
	dayStr := r.FormValue("day")
	dayInt, err := dayStringToInt(dayStr)
	if err != nil {
		http.Error(w, "Invalid day", 400)
		return
	}

	newTask := Task{Task: task, Days: int16(dayInt)}
	err = db.InsertTask(newTask)
	if err != nil {
		http.Error(w, "Failed to add task", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error", 400)
		return
	}
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		http.Error(w, "Failed to delete task", 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error", 400)
		return
	}
	idStr := r.FormValue("id")
	task := r.FormValue("task")
	dayStr := r.FormValue("day")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}
	dayInt, err := dayStringToInt(dayStr)
	if err != nil {
		http.Error(w, "Invalid day", 400)
		return
	}

	newTask := Task{Task: task, Days: int16(dayInt)}
	err = db.UpdateTask(id, newTask)
	if err != nil {
		http.Error(w, "Failed to update task", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// Преобразуем день из строки в число для хранения
// Monday=1, Tuesday=2 и т.д.
func dayStringToInt(day string) (int, error) {
	switch day {
	case "monday":
		return 1, nil
	case "tuesday":
		return 2, nil
	case "wednesday":
		return 3, nil
	case "thursday":
		return 4, nil
	case "friday":
		return 5, nil
	case "saturday":
		return 6, nil
	case "sunday":
		return 7, nil
	default:
		return 0, fmt.Errorf("invalid day")
	}
}

// и обратно если надо (для отображения)
func dayIntToString(day int16) string {
	switch day {
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	case 7:
		return "Sunday"
	default:
		return "Unknown"
	}
}
