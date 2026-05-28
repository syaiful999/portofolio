package utils

type Response struct {
	IsSuccess bool
	Message   string
	Status    int
	Data      interface{}
}
