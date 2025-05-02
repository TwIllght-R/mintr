package handler

import (
	"context"
	"encoding/json"
	"log"
	"redirect-service/domain/entities"
	"redirect-service/domain/events"
	"redirect-service/domain/interfaces"
)

type urlCreatedEventHandler struct {
	redirectUseCase interfaces.RedirectURLUsecase
}

func NewURLCreatedEventHandler(redirectUseCase interfaces.RedirectURLUsecase) interfaces.EventHandler {
	return &urlCreatedEventHandler{
		redirectUseCase: redirectUseCase,
	}
}

func (h *urlCreatedEventHandler) HandleEvent(msg []byte) {
	var evt events.URLCreatedEvent
	if err := json.Unmarshal(msg, &evt); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}
	_, err := h.redirectUseCase.CreateRedirectURL(context.Background(), entities.RedirectURL{
		UUID:        evt.UUID,
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
		log.Println("Error creating redirect URL:", err)
		return
	}
	log.Println("Redirect URL created successfully")
}
