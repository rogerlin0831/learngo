package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	server "test/learngo/tryserver/services"

	"github.com/gorilla/mux"
)

type User struct {
	Name     string
	Password string
}

type ApiResponse struct {
	ResultCode    string
	ResultMeggage interface{}
}

var UserData map[string]string

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

func CheckLogin(use User) string {

	if CheckUser(use.Name) {
		if CheckPassword(use.Name, use.Password) {
			return "Login Success"
		}
		return "Password Error"
	}
	return "UserName Error"

}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		// read error.
		fmt.Println(err)
		return
	}

	var loginData User
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
	var newData User
	_ = json.Unmarshal(body, &newData)

	if UserData == nil {
		UserData = map[string]string{}
	}
	if _, find := UserData[newData.Name]; find {
		// error name same
		// do something...
	}
	UserData[newData.Name] = newData.Password

	defer r.Body.Close()
	response := ApiResponse{"200", UserData}

	server.ResopnseWithJson(w, http.StatusOK, response)
}

func GetPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	queryUser := vars["user"]
	fmt.Println(queryUser)
	var target User

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
