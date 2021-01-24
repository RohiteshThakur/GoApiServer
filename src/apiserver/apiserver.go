/*
Before you begin, update $GOPATH to /home/ubuntu/GoAPIServer in ~/.profile
If you want to create a library to be called by the main program (e.g. apifuncs.go), create it under src/ folder and not under pkg/

To execute:
# ./apiserver/src/apiserver/apiserver
http://127.0.0.1:8000/api/book

*/

package main

import (
	"fmt"
	//"encoding/json"
	"apiserver/apifuncs"	// This local package will be searched under $GOROOT/src and $GOPATH/src , so make sure you start the path after src/ while importing the package. 
	                    	// Also Project's (apiserver) pkg folder will be automatically created with apifuncs.a shared object file.
	"net/http"		// HTTP's "compiled" package can be found under: $GOROOT/pkg
	"os"
)

func main () {

	http.HandleFunc("/", index)
	http.HandleFunc("/api/books", apifuncs.HandleBooksRequest)
	http.HandleFunc("/api/books/", apifuncs.HandleBookRequest)

	http.ListenAndServe(port(), nil)
}

func port() string{
	port := os.Getenv("PORTNUM")
	if (len(port) == 0){
		port = "8000"
	}
	return (":" + port)
	
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to My Books' Library")
	
}