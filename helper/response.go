package helper

type Response struct{
	Status bool `json:"status"`
	Code int `json:"code"`
	Data interface{} `json:"data"`
}