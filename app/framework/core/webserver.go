package core

import (
	"coralscale/app/framework/utility"
	"fmt"
	"log"
	"net/http"
)

func InitWebService(handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc("/", handler)
	utility.LoadEnv()
	port := utility.GetEnv("PORT", "8080")
	fmt.Println("Running server on port: " + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
