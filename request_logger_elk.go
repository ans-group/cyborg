package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/olivere/elastic.v3"
)

type RequestLoggerELKPayload struct {
	RequestLogEntry
	Time time.Time `json:"@timestamp"`
	UUID string    `json:"uuid"`
}

type RequestLoggerELK struct {
	client        *elastic.Client
	serverAddress string
	index         string
	uuid          string
}

func NewRequestLoggerELK() (*RequestLoggerELK, error) {
	address := viper.GetString("logger.elk.address")
	if len(address) < 1 {
		return nil, errors.New("Invalid ELK address provided")
	}

	index := viper.GetString("logger.elk.index")
	if len(index) < 1 {
		return nil, errors.New("Invalid ELK index provided")
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(address))
	if err != nil {
		return nil, fmt.Errorf("Failed to create ELK client: %s", err)
	}

	return &RequestLoggerELK{
		client: client,
		index:  index,
		uuid:   uuid.New().String(),
	}, nil
}

func (l *RequestLoggerELK) Log(entry RequestLogEntry) {
	row := RequestLoggerELKPayload{
		Time:            time.Now(),
		UUID:            l.uuid,
		RequestLogEntry: entry,
	}

	_, err := l.client.Index().Index(l.index).Type("cyborgpayload").BodyJson(row).Refresh(true).Do()
	if err != nil {
		log.Printf("Error inserting ELK entry: %s", err)
	}
}
