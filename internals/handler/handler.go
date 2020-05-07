package handler

import (
	"encoding/json"
	"mqViewer/internals/model"
	"net/http"
	"strconv"

	"github.com/soypita/mq-golang-jms20/mqjms"
)

type mQViewerHandler struct {
	connFactory *mqjms.ConnectionFactoryImpl
}

func NewMQViewerHandler() *mQViewerHandler {
	return &mQViewerHandler{}
}

func (h *mQViewerHandler) CreateNewConnectionWithParams(w http.ResponseWriter, r *http.Request) {
	var connectParams model.CreateConnectionRequest
	err := json.NewDecoder(r.Body).Decode(&connectParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsePort, err := strconv.Atoi(connectParams.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.connFactory = &mqjms.ConnectionFactoryImpl{
		QMName:      connectParams.ManagerName,
		Hostname:    connectParams.Host,
		PortNumber:  parsePort,
		ChannelName: connectParams.ChannelName,
		UserName:    connectParams.Username,
		Password:    connectParams.Password,
		BrowseMode:  true,
	}
}
