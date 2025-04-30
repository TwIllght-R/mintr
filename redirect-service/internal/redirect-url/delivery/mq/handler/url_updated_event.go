package handler

import (
	"context"
	"encoding/json"
	"log"
	"redirect-service/domain/entities"
	"redirect-service/domain/events"
	"redirect-service/domain/interfaces"
)

type urlUpdatedEventHandler struct {
	redirectUseCase interfaces.RedirectURLUsecase
}

func NewURLUpdatedEventHandler(redirectUseCase interfaces.RedirectURLUsecase) interfaces.EventHandler {
	return &urlUpdatedEventHandler{
		redirectUseCase: redirectUseCase,
	}
}
func (h *urlUpdatedEventHandler) HandleEvent(msg []byte) {
	var evt events.URLUpdatedEvent
	if err := json.Unmarshal(msg, &evt); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}
	_, err := h.redirectUseCase.UpdateRedirectURL(context.Background(), evt.UUID, entities.RedirectURL{
		ShortCode:   evt.ShortCode,
		OriginalURL: evt.OriginalURL,
		ExpiresAt:   evt.ExpiresAt,
		UTMSource:   evt.UTMSource,
		UTMMedium:   evt.UTMMedium,
		UTMCampaign: evt.UTMCampaign,
		UTMTerm:     evt.UTMTerm,
		UTMContent:  evt.UTMContent,
	})
	if err != nil {
		log.Println("Error updating redirect URL:", err)
		return
	}
	log.Println("Redirect URL updated successfully")
}
