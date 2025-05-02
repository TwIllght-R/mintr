package handler

import (
	"analytics-service/domain/entities"
	"analytics-service/domain/events"
	"analytics-service/domain/interfaces"
	"context"
	"encoding/json"
	"log"
)

type urlClickedEventHandler struct {
	analyticsUseCase interfaces.AnalyticsUseCase
}

func NewURLClickedEventHandler(analyticsUseCase interfaces.AnalyticsUseCase) interfaces.EventHandler {
	return &urlClickedEventHandler{
		analyticsUseCase: analyticsUseCase,
	}
}

func (h *urlClickedEventHandler) HandleEvent(msg []byte) {
	var evt events.URLClickedEvent
	if err := json.Unmarshal(msg, &evt); err != nil {
		log.Println("error unmarshalling event:", err)
		return
	}
	err := h.analyticsUseCase.RecordClicked(context.Background(), entities.URLClicked{
		UrlUUID:   evt.UrlUUID,
		ShortCode: evt.ShortCode,
		Timestamp: evt.Timestamp,
		IP:        evt.IP,
		UserAgent: evt.UserAgent,
		Referer:   evt.Referer,
		OwnerUUID: evt.OwnerUUID,
	})
	if err != nil {
		log.Println("error creating analytics url:", err)
		return
	}
	log.Println("analytics url Clicked successfully")
}
