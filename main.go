package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Note struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID   string             `json:"userid,omitempty" bson:"userid,omitempty"`
	Title    string             `json:"title,omitempty" bson:"title,omitempty"`
	NoteBody string             `json:"notebody,omitempty" bson:"notebody,omitempty"`
}

func main() {
	fmt.Println("Application starts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/createuser", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/allusers", GetAllUsers).Methods("GET")
	router.HandleFunc("/signin/{username}/{password}", SiginInEndpoint).Methods("GET")
	router.HandleFunc("/savenote", SaveNoteEndpoint).Methods("POST")
	router.HandleFunc("/get_all_notes", GetAllNotes).Methods("GET")
	router.HandleFunc("/getAllNotesForMember/{userid}", GetAllNotesForMemberEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}

func CreateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	collection := client.Database("note-taking-app").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(w).Encode(result)
}

func GetAllUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var users []User
	collection := client.Database("note-taking-app").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(users)

}

func SiginInEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	username := params["username"]
	password := params["password"]
	var user User
	collection := client.Database("note-taking-app").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{Username: username, Password: password}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func SaveNoteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var note Note
	_ = json.NewDecoder(request.Body).Decode(&note)
	collection := client.Database("note-taking-app").Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, note)
	json.NewEncoder(response).Encode(result)
}

func GetAllNotes(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var notes []Note
	collection := client.Database("note-taking-app").Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var note Note
		cursor.Decode(&note)
		notes = append(notes, note)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(notes)

}

func GetAllNotesForMemberEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	userId := params["userid"]
	var notes []Note
	collection := client.Database("note-taking-app").Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"userid": userId})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var note Note
		cursor.Decode(&note)
		notes = append(notes, note)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(notes)

}
