package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

type Task struct{ 
	ID int64
	Task string
	Days int16
}

func main() {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=todolist_db sslmode=disable ")

	if err!= nil{
		log.Fatal(err)
	}

	defer db.Close()

	if err := db.Ping(); err!=nil{
		log.Fatal(err)
	} 
	
	
	// err = insertTask(db, Task{
	// 	Task: "poop time", 
	// 	Days: 1,
	// })
	
	
	users,err := getUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	
	for _,c := range users{ 
		fmt.Println(c)
	}

}




func getUsers(db *sql.DB) ([]Task, error) { 
	rows, err := db.Query("select * from taskinf")
	if err!=nil{ 
		return nil, err
	}

	defer rows.Close()

	users := make([]Task, 0)
	for rows.Next()  {
		u:= Task{}
		err := rows.Scan(&u.ID, &u.Task, &u.Days)
		if err!=nil{ 
			log.Fatal(err)
		}
		users = append(users, u)
	}

	err =rows.Err()
	if err!=nil{ 
		return nil, err
	}

	return users, nil

}


func getUserByDays(db *sql.DB, days int) (Task, error) { 
	var u Task
	err := db.QueryRow("select * from taskinf where days = $1", days).Scan(&u.ID, &u.Task, &u.Days)

	return u, err
}


func insertTask(db *sql.DB, u Task) error {
	tx,err := db.Begin()
	if err!=nil{ 
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("insert into taskinf(task, days) values ($1, $2)", u.Task, u.Days)
	if err!=nil{ 
		return  err
	}
	
	_, err = tx.Exec("insert into logs(entity, action) values ($1, $2)", "task", "created")
	if err!=nil{ 
		return  err
	}


	return tx.Commit()
}


func deleteTask(db *sql.DB, id int) error {
	tx,err := db.Begin()
	if err!=nil{ 
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("delete from taskinf where id = $1", id)
	if err!=nil{ 
		return  err
	}
	
	_, err = tx.Exec("insert into logs(entity, action) values ($1, $2)", "task", "deleted")
	if err!=nil{ 
		return  err
	}
	return tx.Commit()
}


func updateTask(db *sql.DB, id int, newTask Task) error {
	tx,err := db.Begin()
	if err!=nil{ 
		return err
	}

	defer tx.Rollback()
	_, err = db.Exec("update taskinf set task=$1, days=$2 where id = $3",newTask.Task, newTask.Days, id)
	if err!=nil{ 
		return  err
	}
	
	_, err = tx.Exec("insert into logs(entity, action) values ($1, $2)", "task", "updated")
	if err!=nil{ 
		return  err
	}
	return tx.Commit()
}