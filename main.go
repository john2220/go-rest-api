package main

import (
	"encoding/json"		
	"log"
	"net/http"
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
	"context"
	// "time"
	// "fmt"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Book Struct(Model)
type Book struct {
	ID 				string `json:"id"`
	Isbn 			string `json:"isbn"`
	Title 		string `json:"title"`
	Author 		*Author `json:"author"`
}

type Author struct {
	Firstname  	string `json:"firstname"`
	Lastname  	string `json:"lastname"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range books {
		if item.ID == params["id"] {		
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}


func GetClient() *mongo.Client {
    clientOptions := options.Client().ApplyURI("mongodb+srv://johnjohn:<password>@cluster0-601z6.mongodb.net/test?retryWrites=true&w=majority")
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    err = client.Connect(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    return client
}


func main() {
	// Init Router
	r := mux.NewRouter()

	// Data	
	books = append(books, Book{ID: "1", Isbn: "sdf980", Title: "Book one", Author: &Author {Firstname: "John", Lastname: "Dodo"}})
	books = append(books, Book{ID: "2", Isbn: "3240", Title: "Book two", Author: &Author {Firstname: "Anna", Lastname: "Karenina"}})

	// Routes
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://johnjohn:<password>@cluster0-601z6.mongodb.net/test?retryWrites=true&w=majority"))
	// if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// 	err = client.Connect(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer client.Disconnect(ctx)
	// 	err = client.Ping(ctx, readpref.Primary())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(databases)

  c := GetClient()
    err := c.Ping(context.Background(), readpref.Primary())
    if err != nil {
        log.Fatal("Couldn't connect to the database", err)
    } else {
        log.Println("Connected!")
    }
    
	log.Fatal(http.ListenAndServe(":8000", r))
}


