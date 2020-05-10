package services

import (
	"fmt"
	"mqViewer/internals/model"
	"strconv"

	"github.com/soypita/mq-golang-jms20/jms20subset"
	"github.com/soypita/mq-golang-jms20/mqjms"
)

const (
	messageTextConvertError = "message body cannot convert to text"
)

// MQService is a basic interface for mq service
type MQService interface {
	CreateConnection(req *model.CreateConnectionRequest, browseMode bool) error
	GetAllMessages() ([]*model.Message, error)
}

type mqServiceImpl struct {
	connFactory mqjms.ConnectionFactoryImpl
	mqCtx       jms20subset.JMSContext
	queue       jms20subset.Queue
	consumer    jms20subset.JMSConsumer
	producer    jms20subset.JMSProducer
}

// NewDefaultMQService create new MQService instance
func NewDefaultMQService() *mqServiceImpl {
	return &mqServiceImpl{}
}

// GetAllMessages return all messages from queue in browse mode
func (s *mqServiceImpl) GetAllMessages() ([]*model.Message, error) {
	if s.consumer == nil {
		return nil, fmt.Errorf("connection to queue is not configured")
	}
	if s.connFactory.BrowseMode != true {
		return nil, fmt.Errorf("connection to queue configured in non browse mode")
	}
	var resList []*model.Message
	rawMsgList, err := s.consumer.BrowseAllNoWait()
	if err != nil {
		return nil, fmt.Errorf("error during reading from queue: %w", err)
	}

	for _, rawMsg := range rawMsgList {
		resList = append(resList, s.convertToMsgResp(rawMsg))
	}
	return resList, nil
}

// CreateConnection init connection to MQ with provided browse mode
func (s *mqServiceImpl) CreateConnection(req *model.CreateConnectionRequest, browseMode bool) error {
	parsePort, err := strconv.Atoi(req.Port)
	if err != nil {
		return err
	}
	s.connFactory = mqjms.ConnectionFactoryImpl{
		QMName:      req.ManagerName,
		Hostname:    req.Host,
		PortNumber:  parsePort,
		ChannelName: req.ChannelName,
		UserName:    req.Username,
		Password:    req.Password,
		BrowseMode:  browseMode,
	}

	ctx, contextErr := s.connFactory.CreateContext()
	if contextErr != nil {
		return err
	}
	s.mqCtx = ctx
	s.queue = ctx.CreateQueue(req.Queue)

	return nil
}

func (s *mqServiceImpl) convertToMsgResp(rawMsg jms20subset.Message) *model.Message {
	resp := &model.Message{
		MessageID: rawMsg.GetJMSMessageID(),
		Timestamp: rawMsg.GetJMSTimestamp(),
	}

	switch msg := rawMsg.(type) {
	case jms20subset.TextMessage:
		resp.MessageBody = *msg.GetText()
	default:
		resp.MessageBody = messageTextConvertError
	}
	return resp
}
