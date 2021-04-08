package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name     string
	Password string
}

const (
	userName = "root"
	password = "Ab123456"
	host     = "127.0.0.1"
	port     = "3306"
	dbName   = "golang"

	CreateUserTable = `
	create table userData(
		u_id INT NOT NULL AUTO_INCREMENT,
		u_Name VARCHAR(64) NOT NULL,
		u_Password VARCHAR(64) NOT NULL,
		datatime datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		createtime datetime DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY ( u_id )
	 );
	`
)

var Gomysql *sql.DB

func init() {

	path := strings.Join([]string{userName, ":", password, "@tcp(", host, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println(path)

	gosql, err := sql.Open("mysql", path)
	if err != nil {
		fmt.Println(err)
	}
	Gomysql = gosql
	// Gomysql.Query( "SHOW TABLES LIKE 'UserData'" )

	rows, err := gosql.Query("SHOW TABLES LIKE 'UserData'")
	if err != nil {
		fmt.Println(err)
	}

	if rows.Next() {
		// find table.
		gosql.Query("SELECT * FROM UserData")
	} else {
		gosql.Query(CreateUserTable)
	}
}

func GetDBData() []User {
	data := []User{}
	row, err := Gomysql.Query("SELECT * FROM UserData")
	if err != nil {
		fmt.Print(err)
		return nil
	}
	for row.Next() {
		var id int
		var name string
		var password string
		var updatetime string
		var createtime string
		err = row.Scan(&id, &name, &password, &updatetime, &createtime)
		if err != nil {
			fmt.Print(err)
			return nil
		}
		data = append(data, User{name, password})
	}
	return data
}

func DBAddData(data User) {

	//add data
	stmt, err := Gomysql.Prepare("INSERT UserData SET u_Name=?,u_Password=?")
	checkErr(err)

	res, err := stmt.Exec(data.Name, data.Password)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println("DBAddData new :", id)
}

func DBDeleteData(data User) {

	stmt, err := Gomysql.Prepare("DELETE FROM UserData WHERE u_Name=? AND u_Password=? ")
	checkErr(err)

	res, err := stmt.Exec(data.Name, data.Password)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)
}

func DBUpdateData(data User) {
	//update
	stmt, err := Gomysql.Prepare("UPDATE UserData SET u_Password=? WHERE u_Name=?")
	checkErr(err)

	res, err := stmt.Exec(data.Password, data.Name)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)
}

func Close() {
	Gomysql.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
