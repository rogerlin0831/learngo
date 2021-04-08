package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"test/learngo/tryserver/db"
	server "test/learngo/tryserver/services"

	"github.com/gorilla/mux"
)

const (
	UserCheckOK   = "Login Success"
	PasswordError = "Password Error"
	UserNameError = "UserName Error"
)

type ApiResponse struct {
	ResultCode    string
	ResultMeggage interface{}
}

var UserData map[string]string

func init() {
	UserData = map[string]string{}
}

func InitData() {
	dbData := db.GetDBData()
	for _, row := range dbData {
		UserData[row.Name] = row.Password
	}
	fmt.Println(UserData)
}

func CheckUser(usename string) bool {
	if UserData != nil {
		_, isFind := UserData[usename]
		return isFind
	}
	return false
}

func CheckPassword(usename string, password string) bool {
	if UserData != nil {
		return UserData[usename] == password
	}
	return false
}

func CheckLogin(use db.User) string {

	if CheckUser(use.Name) {
		if CheckPassword(use.Name, use.Password) {
			return UserCheckOK
		}
		return PasswordError
	}
	return UserNameError

}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		// read error.
		fmt.Println(err)
		return
	}

	var loginData db.User
	err = json.Unmarshal(body, &loginData)
	if err != nil {
		// data error.
		return
	}
	result := CheckLogin(loginData)

	defer r.Body.Close()
	response := ApiResponse{"200", result}

	server.ResopnseWithJson(w, http.StatusOK, response)
}

func AddData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		fmt.Println(err)
		return
	}
	var newData db.User
	_ = json.Unmarshal(body, &newData)

	if _, find := UserData[newData.Name]; find {
		// error name same
		// do something...
	}
	UserData[newData.Name] = newData.Password

	db.DBAddData(newData)

	defer r.Body.Close()
	response := ApiResponse{"200", UserData}

	server.ResopnseWithJson(w, http.StatusOK, response)
}

func DeleteData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		fmt.Println(err)
		return
	}
	var deleteData db.User
	_ = json.Unmarshal(body, &deleteData)

	result := CheckLogin(deleteData)
	response := ApiResponse{"200", UserData}

	if result == UserCheckOK {
		delete(UserData, deleteData.Name)
		db.DBDeleteData(deleteData)
	} else {
		response.ResultMeggage = "login fail, please check name and password"
	}

	defer r.Body.Close()

	server.ResopnseWithJson(w, http.StatusOK, response)
}

func UpdateData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		fmt.Println(err)
		return
	}
	var updateData db.User
	_ = json.Unmarshal(body, &updateData)

	_, find := UserData[updateData.Name]

	response := ApiResponse{"200", UserData}

	if !find {
		response.ResultMeggage = "UserName Not Found"
	} else {
		db.DBUpdateData(updateData)
		UserData[updateData.Name] = updateData.Password
	}

	defer r.Body.Close()

	server.ResopnseWithJson(w, http.StatusOK, response)
}

func GetPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	queryUser := vars["user"]
	fmt.Println(queryUser)
	var target db.User

	response := ApiResponse{"200", &target}
	if CheckUser(queryUser) {
		target.Name = queryUser
		target.Password = UserData[queryUser]
	} else {
		response.ResultMeggage = "Error UserName"
	}
	fmt.Println(target)
	server.ResopnseWithJson(w, http.StatusOK, response)
}
