package todoCont

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	todoMod "github.com/ankush/todo/Model"
	mongodb "github.com/ankush/todo/Mongodb"
	"github.com/go-chi/chi"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render

func init() {
	rnd = renderer.New(renderer.Options{
		JSONIndent: true,
	})
}

type todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}

func fetchTodos(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var todos []todoMod.TodoModel
	cursor, err := mongodb.Collection.Find(ctx, bson.M{})
	if err != nil {
		rnd.JSON(res, http.StatusInternalServerError, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var t todoMod.TodoModel
		if err := cursor.Decode(&t); err != nil {
			rnd.JSON(res, http.StatusInternalServerError, renderer.M{
				"message": "failed to decode data",
				"error":   err,
			})
			return
		}
		todos = append(todos, t)
	}
	rnd.JSON(res, http.StatusOK, renderer.M{
		"data": todos,
	})
}

func CreateTodo(res http.ResponseWriter, req *http.Request) {
	var t todo
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		rnd.JSON(res, http.StatusInternalServerError, renderer.M{
			"status":  false,
			"message": "Cannot decode todo",
		})
		return
	}
	if t.Title == "" {
		rnd.JSON(res, http.StatusBadRequest, renderer.M{
			"status":  false,
			"message": "Title field is required",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	data := todoMod.TodoModel{
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	result, err := mongodb.Collection.InsertOne(ctx, data)
	if err != nil {
		rnd.JSON(res, http.StatusInternalServerError, renderer.M{
			"message": "Failed to create todo",
			"error":   err.Error(),
		})
		return
	}
	rnd.JSON(res, http.StatusOK, renderer.M{
		"status":  "success",
		"todo_id": result.InsertedID,
	})
}

func UpdateTodo(res http.ResponseWriter, req *http.Request) {
	id := strings.TrimSpace(chi.URLParam(req, "id"))
	var t todo
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		rnd.JSON(res, http.StatusBadRequest, renderer.M{
			"status":  false,
			"message": "Cannot decode data",
			"error":   err.Error(),
		})

	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": oid,
	}
	update := todo{
		Title:     t.Title,
		Completed: t.Completed,
	}
	fmt.Println(filter)
	fmt.Println(update)
	result, err := mongodb.Collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil || result.MatchedCount == 0 {
		rnd.JSON(res, http.StatusNotFound, renderer.M{
			"status":  false,
			"message": "TODO not found",
		})
		return
	}
	rnd.JSON(res, http.StatusOK, renderer.M{
		"status":  "success",
		"message": "Todo updated successfully",
	})

}

func DeleteTodo(res http.ResponseWriter, req *http.Request) {
	id := strings.TrimSpace(chi.URLParam(req, "id"))
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": oid,
	}
	result, err := mongodb.Collection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		rnd.JSON(res, http.StatusNotFound, renderer.M{
			"message": "Failed to delete todo or no document found",
			"error":   err.Error(),
		})
		return
	}

	rnd.JSON(res, http.StatusOK, renderer.M{
		"status":  "success",
		"message": "Todo deleted successfully",
	})

}
func GetTodos(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var todos []todo
	cur, err := mongodb.Collection.Find(ctx, bson.M{})
	if err != nil {
		rnd.JSON(res, http.StatusInternalServerError, renderer.M{
			"message": "Failed to get todos",
			"error":   err.Error(),
		})
		return
	}
	err = cur.All(ctx, &todos)
	if err != nil {
		rnd.JSON(res, http.StatusInternalServerError, renderer.M{
			"message": "Failed to get todos",
			"error":   err.Error(),
		})
		return
	}
	rnd.JSON(res, http.StatusOK, renderer.M{
		"status":  "success",
		"message": "Todos retrieved successfully",
		"todos":   todos,
	})

}
