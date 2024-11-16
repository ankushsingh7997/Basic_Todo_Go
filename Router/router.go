package todoRouter

import (
	"net/http"

	todoCont "github.com/ankush/todo/Controller"
	"github.com/go-chi/chi"
)

func TodoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", todoCont.GetTodos)
		r.Post("/", todoCont.CreateTodo)
		r.Put("/{id}", todoCont.UpdateTodo)
		r.Delete("/{id}", todoCont.DeleteTodo)
	})
	return rg
}
