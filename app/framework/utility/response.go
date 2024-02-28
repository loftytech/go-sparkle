package utility

import "fmt"

type Response struct {

}

func Send(response map[string]any, status_code int64) {
	fmt.Println("This is a request bruv")
}