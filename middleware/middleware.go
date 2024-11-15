package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/d4ny4z0rd/godogo/model"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func loadTheEnv(){
	err := godotenv.Load(".env")
	if err!=nil {
		log.Fatal("Error loading the dot env file")
	}
}

func createDBInstance(){
	connectionString, dbName, colName := os.Getenv("DB_URI"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err!=nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongodb")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance created")
}

func init(){
	loadTheEnv()
	createDBInstance()
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")

	payload := getAllTodos()
	json.NewEncoder(w).Encode(payload)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	var todo model.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	insertOneTodo(todo)
	json.NewEncoder(w).Encode(todo)
}

func CompleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := mux.Vars(r)
	completeTodo(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func UndoTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := mux.Vars(r)
	undoTodo(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","DELETE")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := mux.Vars(r)
	deleteTodo(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")

	count := deleteAllTodos()
	json.NewEncoder(w).Encode(count)	
}

func getAllTodos() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err!=nil {
		log.Fatal(err)
	}

	var results []primitive.M
	
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err!=nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func insertOneTodo(todo model.Todo) {
	res, err := collection.InsertOne(context.Background(), todo)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single record", res.InsertedID)
}

func completeTodo(todo string) {
	id, _ := primitive.ObjectIDFromHex(todo)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":true}}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count:",res.ModifiedCount)
}

func undoTodo(todoId string) {
	id, _ := primitive.ObjectIDFromHex(todoId)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":false}}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count:", res.ModifiedCount)
}

func deleteTodo(todoId string){
	id, _ := primitive.ObjectIDFromHex(todoId)
	filter := bson.M{"_id" : id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted document", d.DeletedCount)
}

func deleteAllTodos() int64 {
	d, err := collection.DeleteMany(context.Background(), bson.D{{}})
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted document,",d.DeletedCount)
	return d.DeletedCount
}