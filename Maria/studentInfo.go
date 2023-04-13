package Maria

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetStudentInfo(campemail string, db *sql.DB) (string, string) {
	var ID string
	var Birthday string
	// Query the database
	statement := `SELECT ID,Birthday FROM Student WHERE CampEmail=?;`
	row := db.QueryRow(statement, campemail)
	switch err := row.Scan(&ID, &Birthday); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(ID, Birthday)
	default:
		panic(err)

	}
	return ID, Birthday
}

func InitMaria() *sql.DB {
	// Open a new database connection
	db, err := sql.Open("mysql", "it:Acumen321@tcp(192.168.5.238:3306)/Ace?charset=utf8")
	if err != nil {
		panic(err)
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db

}
