package Maria

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetStudentInfo(campemail,term string,year int,db *sql.DB) (string,string) {
	var ID string
	var Program_code string
	// Query the database
	statementinfo := `SELECT ID FROM Student WHERE CampEmail=?;`
	row := db.QueryRow(statementinfo, campemail)
	switch err := row.Scan(&ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		
		statementPro:=`SELECT Program_code FROM Enrollment WHERE StudentID=? AND Term=? AND TermYear=? LIMIT 1;`
		row2 := db.QueryRow(statementPro, ID,term,year)
		row2.Scan(&Program_code);
		fmt.Println(ID,Program_code)
	default:
		panic(err)

	}
	return ID,Program_code
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
