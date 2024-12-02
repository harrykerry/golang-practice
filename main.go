package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)



type User struct {
	ID         int64
	FirstName  string
	SecondName string
	Email      string
}

func main() {

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	// Set log output to file
	log.SetOutput(logFile)

	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cnfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               os.Getenv("DBDATABASE"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cnfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping()

	if pingErr != nil {

		log.Fatal(pingErr)
	}

	fmt.Println("Database Connction established")

	users,err := getUserDetails(db)

	if (err !=nil){
		log.Fatal(err)
	}
	for i, user := range users {
		fmt.Printf("Index: %d, User: %s %s, Email: %s\n", i, user.FirstName, user.SecondName, user.Email)
	}

	user,err := getSingleUser(db,1)

	if (err !=nil){
		log.Fatal(err)
	}

	fmt.Println(user)

	addedId,err := insertSIngleUser(db,User{
		FirstName: "Omosh",
		SecondName: "Olosh",
		Email: "omosh.lohso@gmail.com",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(addedId)


}


func getUserDetails(db *sql.DB) ([]User,error){

	var users []User

	rows,err := db.Query("SELECT fname,sname,email FROM user")

	if(err !=nil){

		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var user User

		if err := rows.Scan(&user.FirstName, &user.SecondName, &user.Email);err!=nil{
			return nil, err 
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func getSingleUser(db *sql.DB,id int) (User,error){

	var user User

	row := db.QueryRow("SELECT fname,sname,email FROM user where id = ?",id)

    err := row.Scan(&user.FirstName,&user.SecondName, &user.Email)

	if err != nil {
        return user, err
    }

	return user,err

}

func insertSIngleUser(db *sql.DB,singleuser User) (int64,error){

	result,err := db.Exec("INSERT INTO user (fname,sname,email) VALUES (?,?,?)", singleuser.FirstName,singleuser.SecondName,singleuser.Email)

	if (err !=nil){

		return 0, err
	}

	id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

	fmt.Println(result)
    return id, nil
	
}







