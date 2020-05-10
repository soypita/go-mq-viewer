package handler

import (
	"encoding/json"
	"mqViewer/internals/model"
	"mqViewer/internals/services"
	"net/http"
)

// MQViewerHandler handler for http requests
type MQViewerHandler struct {
	mqService services.MQService
}

// NewMQViewerHandler create instanse of MQ viewer handler
func NewMQViewerHandler(serv services.MQService) *MQViewerHandler {
	return &MQViewerHandler{
		mqService: serv,
	}
}

// BrowseAllMessages handler return all messages in specified queue
func (h *MQViewerHandler) BrowseAllMessages(w http.ResponseWriter, r *http.Request) {
	resMsgList, err := h.mqService.GetAllMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resMsgList)
}

// CreateNewConnectionWithParams create new connection to MQ queue with provided params
func (h *MQViewerHandler) CreateNewConnectionWithParams(w http.ResponseWriter, r *http.Request) {
	var connectParams model.CreateConnectionRequest
	err := json.NewDecoder(r.Body).Decode(&connectParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.mqService.CreateConnection(&connectParams, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(connectParams)
}
