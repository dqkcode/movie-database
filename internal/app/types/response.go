package types

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Response struct {
		Code    int
		Data    interface{}
		Message string
	}
	ResponseInfo struct {
		Status  int    `yaml:"status"`
		Code    int    `yaml:"code"`
		Message string `yaml:"message"`
	}
	UserResponse struct {
		Created        ResponseInfo `yaml:"created"`
		DuplicateEmail ResponseInfo `yaml:"duplicate_email"`
		UpdateFailed   ResponseInfo `yaml:"update_failed"`
		CreateFailed   ResponseInfo `yaml:"create_failed"`
		DeleteFailed   ResponseInfo `yaml:"delete_failed"`
		UserNotFound   ResponseInfo `yaml:"user_not_found"`
	}
	AuthResponse struct {
		UserLocked    ResponseInfo `yaml:"user_locked"`
		EmailNotExist ResponseInfo `yaml:"email_not_exist"`
		PasswordWrong ResponseInfo `yaml:"password_wrong"`
		Unauthorized  ResponseInfo `yaml:"unauthorized"`
	}
	NormalResponse struct {
		Success    ResponseInfo `yaml:"success"`
		NotFound   ResponseInfo `yaml:"not_found"`
		TimeOut    ResponseInfo `yaml:"timeout"`
		BadRequest ResponseInfo `yaml:"bad_request"`
		Internal   ResponseInfo `yaml:"internal"`
	}

	AllResponse struct {
		NormalResponse NormalResponse `yaml:"normal"`
		UserResponse   UserResponse   `yaml:"user"`
		AuthResponse   AuthResponse   `yaml:"auth"`
	}
)

func ResponseJson(w http.ResponseWriter, data interface{}, resinfo ResponseInfo) {

	res := &Response{}
	res.Code = resinfo.Code
	res.Message = resinfo.Message
	res.Data = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resinfo.Status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Load() *AllResponse {

	path := os.Getenv("STATUS_FILE_PATH")
	if path == "" {
		path = "configs/status.yml"
	}
	yamlFile, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	all := &AllResponse{}

	if err := yaml.Unmarshal(yamlFile, all); err != nil {
		panic(err)
	}
	return all
}

func Normal() NormalResponse {
	return Load().NormalResponse
}
func User() UserResponse {
	return Load().UserResponse
}
func Auth() AuthResponse {
	return Load().AuthResponse
}
