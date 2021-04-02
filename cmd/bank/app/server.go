package app

import (
	"bgo-homeworks-11/cmd/bank/app/dto"
	"bgo-homeworks-11/pkg/card"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	cardSvc *card.Service
	mux     *http.ServeMux
}

func NewServer(cardSvc *card.Service, mux *http.ServeMux) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/", s.badGateway)
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/addCard", s.addCard)
}

func (s *Server) badGateway(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadGateway)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	suid := r.URL.Query().Get("uid")
	uid, err := strconv.ParseInt(suid, 10, 64)
	if err != nil || uid == 0 {
		makeResponse(w, dto.Result{Result: "Error", ErrorDescription: "Wrong uid format"})
		return
	}

	cards := s.cardSvc.All(r.Context())
	found := false
	dtos := make([]*dto.CardDTO, len(cards))

	for i, c := range cards {
		if c.Id == uid {
			found = true
			dtos[i] = &dto.CardDTO{
				Id:     c.Id,
				Number: c.Number,
				Type:   c.CardType,
			}
		}
	}

	if !found {
		makeResponse(w, dto.Result{Result: "No cards"})
		return
	}
	makeResponse(w, dto.Result{Cards: dtos})

}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		makeResponse(w, dto.Result{Result: "Error", ErrorDescription: "Wrong params"})
		return
	}

	cardType := r.PostForm.Get("type")
	suid := r.PostForm.Get("uid")
	isNewUser := r.PostForm.Get("new") == "y"

	if cardType == "" || suid == "" {
		makeResponse(w, dto.Result{Result: "Error", ErrorDescription: "Wrong params"})
		return
	}

	uid, err := strconv.ParseInt(suid, 10, 64)

	if err != nil {
		makeResponse(w, dto.Result{Result: "Error", ErrorDescription: "Wrong uid"})
		return
	}

	err = s.cardSvc.IssueCard(uid, cardType, isNewUser, r.Context())
	if err != nil {
		makeResponse(w, dto.Result{Result: "Error", ErrorDescription: err.Error()})
		return
	}

	makeResponse(w, dto.Result{Result: "Ok"})
}

func makeResponse(w http.ResponseWriter, dto dto.Result) {
	respBody, err := json.Marshal(dto)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	// по умолчанию статус 200 Ok
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}