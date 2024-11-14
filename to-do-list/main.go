package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Preetham-P/go-projects/to-do-list/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/mongo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render

var db *mgo.Database

var client *mongo.Client

var ctx context.Context

var cancel context.CancelFunc

const (
	dbhostname     string = "127.0.0.1:27017"
	dbname         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type todoModel struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Status    string        `bson:"status"`
	CreatedAt time.Time     `bson:"createdAt"`
}

type todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"time"`
}

func init() {

	client, ctx, cancel = database.Connect()
	rnd = renderer.New()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("There was an error", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {

	filter := bson.M{} // Empty filter matches all documents
	foundtodos, err := database.FindDocuments(client, ctx, dbname, collectionName, filter)
	if err != nil {
		log.Fatal(err)
	}

	todoList := []todo{}

	for _, ft := range foundtodos {
		var t todoModel
		bsonBytes, _ := bson.Marshal(ft)
		bson.Unmarshal(bsonBytes, &t)
		todoList = append(todoList, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Status:    t.Status,
			CreatedAt: t.CreatedAt,
		})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
	})
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	fmt.Println("Main function entered")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port ", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("There was an error %s", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "This id is invalid",
		})
		return
	}

	deleteFilter := bson.M{"_id": bson.ObjectIdHex(id)}
	deleteResult, err := database.DeleteDocument(client, ctx, dbname, collectionName, deleteFilter)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message":      "todo deleted successfully",
		"DeleteResult": deleteResult.DeletedCount,
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "This id is invalid",
		})
		return
	}

	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "No title please add title",
		})
		return
	}

	updateFilter := bson.M{"_id": bson.ObjectIdHex(id)}
	update := bson.M{"title": t.Title, "completed": t.Status}
	updateResult, err := database.UpdateDocument(client, ctx, dbname, collectionName, updateFilter, update)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}
	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"id":      updateResult.ModifiedCount,
	})

}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title is request",
		})
		return
	}
	tm := todoModel{
		ID:        bson.NewObjectId(),
		Title:     t.Title,
		Status:    "New",
		CreatedAt: time.Now().UTC(),
	}

	insertResult, err := database.CreateDocument(client, ctx, dbname, collectionName, tm)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "failed to save todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"id":      insertResult.InsertedID,
	})
}
