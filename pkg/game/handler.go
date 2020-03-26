package game

import (
	"encoding/json"
	"log"
	"net/http"
)

type JoinGameHandler struct {
	actionChannel chan *Action
}

func NewJoinGameHandler(actionChannel chan *Action) *JoinGameHandler {
	return &JoinGameHandler{actionChannel}
}

type JoinGameRequest struct {
	Name   string `json:"name"`
	First  int    `json:"first"`
	Second int    `json:"second"`
}

type JoinGameResponse struct {
	Status int    `json:"status"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (h *JoinGameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "NG: Small Browser Based Game Server")

	request := new(JoinGameRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Join Game - Error decoding json %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := JoinGameResponse{
			Status: http.StatusBadRequest,
			Type:   "Error",
			Title:  "Invalid JSON",
			Detail: err.Error(),
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	arc := make(chan *ActionResponse)

	log.Println("Join Game Request sending to Engine")
	h.actionChannel <- &Action{
		Type:   ActionTypeJoinGame,
		Player: &Player{request.Name, request.First, request.Second},
		Reply:  arc,
	}

	ar := <-arc
	if !ar.Success {
		log.Printf("Join Game Error: %s", ar.Message)
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := JoinGameResponse{
			Status: http.StatusBadRequest,
			Type:   "Error",
			Title:  "Invalid Request",
			Detail: ar.Message,
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := JoinGameResponse{
		Status: http.StatusOK,
		Type:   "Success",
		Title:  "Joined Game",
		Detail: "Welcome to the game, player ;)",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}
