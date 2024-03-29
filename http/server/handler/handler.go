package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"game/internal/service"
	"game/pkg/life"
	"net/http"
	"os"
)

// создадим новый тип для добавления middleware к обработчикам
type Decorator func(http.Handler) http.Handler

// объект для хранения состояния игры
type LifeStates struct {
	service.LifeService
}

func newGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	f, err := os.ReadFile("static/index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(f)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func New(ctx context.Context,
	lifeService service.LifeService,
) (http.Handler, error) {
	serveMux := http.NewServeMux()

	lifeState := LifeStates{
		LifeService: lifeService,
	}

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newWorldGenerate, _ := life.NewWorld(100, 100)
		newWorldGenerate.RandInit(7)
		lifeState.LifeService.CurrentWorld = newWorldGenerate
		newGame(w, r)
	})
	serveMux.HandleFunc("/nextstate", lifeState.nextState)

	return serveMux, nil
}

// функция добавления middleware
func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

// получение очередного состояния игры
func (ls *LifeStates) nextState(w http.ResponseWriter, r *http.Request) {
	worldState := ls.LifeService.NewState()

	err := json.NewEncoder(w).Encode(worldState.Cells)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
