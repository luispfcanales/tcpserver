package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/luispfcanales/tcpserver/model"
)

type HTTPActor struct {
	port    string
	mailBox chan<- model.MessageTCP
}

func NewHTTPActor(port string, mailBox chan<- model.MessageTCP) *HTTPActor {
	s := &HTTPActor{
		port:    port,
		mailBox: mailBox,
	}
	return s
}

func (s *HTTPActor) Run() error {
	http.HandleFunc("/", s.documentation)
	http.HandleFunc("/notify", CorsMiddle(s.handlePostNotification))

	log.Println("[ Start service HTTP: ]", s.port)
	return http.ListenAndServe(s.port, nil)
}

func (s *HTTPActor) documentation(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi guys!"))
}

func (s *HTTPActor) handlePostNotification(w http.ResponseWriter, r *http.Request) {
	log.Println("ok")
	var resMsg model.ResponseMessage
	var reqMsg model.RequestMessage

	if r.Method != http.MethodPost {
		resMsg.Status = http.StatusMethodNotAllowed
		resMsg.Message = "HTTP Method not found"

		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(&resMsg)
		return
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&reqMsg); err != nil {
		resMsg.Status = http.StatusBadRequest
		resMsg.Message = "bad request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resMsg)
		return
	}

	payload := fmt.Sprintf("%s;%s", reqMsg.EventName, reqMsg.Message)
	s.mailBox <- model.MessageTCP{
		From:    "ubicacion-local",
		Payload: []byte(payload),
	}

	resMsg.Status = http.StatusCreated
	resMsg.Message = "Ok"
	resMsg.Data = "Sent Event"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resMsg)
}
