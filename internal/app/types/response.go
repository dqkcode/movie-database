package types

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

const (
	CodeSuccess = "0000"
	CodeFail    = "1000"
)

// func ResponseJson(w http.ResponseWriter, code string, data interface{}, err string) {
// 	json.NewEncoder(w).Encode(Response{
// 		Code:  code,
// 		Data:  data,
// 		Error: err,
// 	})
// }

func ResponseJson(w http.ResponseWriter, data interface{}, resinfo ResponseInfo) {

	res := &Response{}
	res.Code = resinfo.Code
	res.Message = resinfo.Message
	res.Data = data

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resinfo.Status)
}

func Load() *AllResponse {
	// f, err := os.Open("configs/status.yml")
	yamlFile, err := ioutil.ReadFile("configs/status.yml")

	if err != nil {
		panic(err)
	}
	all := &AllResponse{}

	if err := yaml.Unmarshal(yamlFile, all); err != nil {
		panic(err)
	}
	// if err := yaml.NewDecoder(f).Decode(all); err != nil {
	// 	panic(err)
	// }
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
