package handler

import (
	"context"
	"encoding/json"
	"log"
	"redirect-service/domain/events"
	"redirect-service/domain/interfaces"
)

type urlDeletedEventHandler struct {
	redirectUseCase interfaces.RedirectURLUsecase
}

func NewURLDeletedEventHandler(redirectUseCase interfaces.RedirectURLUsecase) interfaces.EventHandler {
	return &urlDeletedEventHandler{
		redirectUseCase: redirectUseCase,
	}
}
func (h *urlDeletedEventHandler) HandleEvent(msg []byte) {
	var evt events.URLDeletedEvent
	if err := json.Unmarshal(msg, &evt); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}
	err := h.redirectUseCase.DeleteRedirectURL(context.Background(), evt.UUID)
	if err != nil {
		log.Println("Error deleting redirect URL:", err)
		return
	}
	log.Println("Redirect URL deleting successfully")
}
