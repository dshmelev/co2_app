package main

import (
	//	"encoding/json"
	"log"
	"net/http"
)

type test_struct struct {
	Test string
}

func main() {
	log.Println("Listen :8080")
	http.HandleFunc("/send", func(w http.ResponseWriter, req *http.Request) {
		//   decoder := json.NewDecoder(req.Body)
		//   var t test_struct
		//   err := decoder.Decode(&t)
		//   if err != nil {
		// 	panic(err)
		//   }
		//		b, err := json.Marshal(req.Body)
		//		if err != nil {
		//			panic(err)
		//		}
		//		defer req.Body.Close()
		//		log.Println(b)
		req.ParseForm()
		log.Println(req.Form)

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
