package sse

import (
	"encoding/json"
	"github.com/kylegk/sse-rest-server/config"
	"github.com/kylegk/sse-rest-server/db"
	"github.com/kylegk/sse-rest-server/models"
	"github.com/r3labs/sse"
)

// IngestData initializes the connection to the SSE server and populate data store
func IngestData(url string) {
	client := CreateClientConnection(url)
	Subscribe(client, insertScore)
}

// Insert event (score) data into the data store
func insertScore(msg *sse.Event) {
	score := models.StudentExam{}
	err := json.Unmarshal(msg.Data, &score)
	if err != nil {
		panic(err)
	}

	err = db.UpsertRow(config.ScoreTable, score)
	if err != nil {
		panic(err)
	}
}
