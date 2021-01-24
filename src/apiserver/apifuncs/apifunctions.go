/* Collections of functions implement basic CRUD actions */


package apifuncs			// This is the name of the folder where this package file is created and not file in which its declared.

import (
	//"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Book struct {
	Title		string	`json:"title"`			// Remember: There should be no spaces on either side of :
	Author		string	`json:"author"`
	ISBN		string	`json:"isbn"`
	Description	string	`json:"description,omitempty"`	// Again, no space on either side of "," or else, warnings will be issued.
}

// This function takes "struct" and returns in JSON in byte array format to the requester over HTTP. So it serves JSON.
func ToJSON(b Book) []byte {
	JSON_data, err := json.Marshal(b)
	if (err != nil){
		panic(err)
	}
	return JSON_data

}

// Receives JSON data as byte array and returns formatted data in Book struct (Go type). This is where json tags from incoming json are useful.
func FromJSON(data []byte) Book {
	b := Book{}
	err := json.Unmarshal(data, &b)
	if (err != nil){
		panic(err)
	}
	return b
}
// A dictionary(map) of books. Stored "in-memory"
var books = map[string]Book {
	"001": Book{Title: "Golang by Example", Author: "Rob Pike", ISBN: "001", Description: "Everything about Golang"},
	"002": Book{Title: "Python by Example", Author: "Dusty Philips", ISBN: "002", Description: "Everything about Python"},		// Note the ","
}

func WriteJSON(w http.ResponseWriter, b map[string]Book){
	JSON_data, err := json.Marshal(b)
	if (err != nil){
		panic(err)
	}
	http_writer_obj := w.Header()
	http_writer_obj.Add("Content-Type", "application/json; charset=utf-8")
	w.Write(JSON_data)
}

func CreateBook (book Book) (string, bool){			// Take struct Book as arguments. 
	_, exists := (books[book.ISBN])					// book.ISBN is the key of dictionary books.
	if exists{
		return "", false
	}
		books[book.ISBN] = book						// If doesn't exist, create a new entry in dictionary books.
		return book.ISBN, true
}


func HandleBooksRequest(w http.ResponseWriter, r *http.Request) {
	switch request_method := r.Method; request_method {		// Switch case on HTTP request methods. Cases defined: GET. POST, PUT and DELETE.
	case http.MethodGet:
		WriteJSON(w, books)
	
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if (err != nil){
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)						// Get JSON []byte array, return struct book.
		isbn, created := CreateBook(book)

		if (created){
			w.Header().Add("Location", "/api/books" + isbn)
			w.WriteHeader(http.StatusCreated)
			w.Write ([]byte("Book added to database"))
			
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([] byte("Unsupported Request Method"))		// ResponseWriter Object expects []byte array and not String
	}
	
}


func Getbook (isbn string) (Book, bool) {
	book, exists := books[isbn]
	return book, exists						//Don't put parenthesis () around return variables when they are more than one.
}

func Putbook (b []byte) {
	Book := FromJSON(b)
	books[Book.ISBN] = Book

}

// DeleteBook removes a book from the map by ISBN key
func DeleteBook(isbn string) {
	delete(books, isbn)
}


func HandleBookRequest(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Path[len("/api/books/"):]			// Filtering URL to get just ISBN number passed.
	
	switch method :=r.Method; method {
	case http.MethodGet:
		book, found := Getbook(isbn)
		if found{
			json_data := ToJSON(book)		// Book is a struct do convert it into JSON Byte array.
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.Write(json_data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if (err != nil){
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)						// Get JSON []byte array, return struct book.
		_, exists := books[book.ISBN]				// If there is an dict element of isbn passed via PUT.
		if exists {									// If yes, then Update.
			books[isbn] = book						// We know isbn is exists, so we can use it directly as dict key.
			w.Write(body)
			w.Write([]byte("Updated"))
		}


	case http.MethodDelete:
		DeleteBook(isbn)
		w.WriteHeader(http.StatusOK)
		
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))

	}

}