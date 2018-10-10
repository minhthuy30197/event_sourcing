package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/minhthuy30197/event_sourcing/model"
)

func (c *Controller) SaveEvent(ev model.Event) error {
	tx, err := c.EventDB.Begin()
	if err != nil {
		return err
	}

	// Encode
	eventDB, err := Encode(ev)
	if err != nil {
		log.Println("Encode loi")
		return err
	}

	// TODO: Kiem tra version

	// Insert EventDB
	err = c.EventDB.Insert(&eventDB)
	if err != nil {
		log.Println(err)
		log.Println("Insert loi")
		tx.Rollback()
		tx.Commit()
		return err
	}
	tx.Rollback()
	return err

	// Apply event. Neu loi thi rollback
	err = ev.Apply(c.Config)
	if err != nil {
		tx.Rollback()
		tx.Commit()
		log.Println("Loi khi update bang read")
		return err 
	}

	return tx.Commit()
}

func BuildBaseEvent(aggregateID, userID, eventType string, data interface{}, version int32) model.Event {
	var event model.Event

	event.AggregateId = aggregateID
	event.Time = time.Now()
	event.UserID = userID
	event.Revision = 1
	event.Version = version
	event.Data = data
	event.EventType = eventType

	return event
}

func Encode(event model.Event) (model.EventSource, error) {
	ret := model.EventSource{}
	var err error

	ret.AggregateId = event.AggregateId
	ret.Time = event.Time
	ret.EventType = event.EventType
	ret.Revision = event.Revision
	ret.UserID = event.UserID
	ret.Version = event.Version

	tmp, err := json.Marshal(event.Data)
	if err != nil {
		return model.EventSource{}, err
	}
	ret.Data = []string{string(tmp)}

	return ret, nil
}

// Decode 
func Decode(event model.EventSource, tmpStruct interface{}) (model.Event, error) {
	var err error
	ret := model.Event{}

	ret.AggregateId = event.AggregateId
	ret.EventType = event.EventType
	ret.Revision = event.Revision
	ret.Time = event.Time
	ret.UserID = event.UserID
	ret.Version = event.Version

	err = json.Unmarshal([]byte(event.Data[0]), &tmpStruct)
	if err != nil {
		fmt.Println("error:", err)
	}
	ret.Data = tmpStruct

	return ret, nil
}

// Events returns **All** the persisted events
func (c *Controller) Events(aggregateID string) ([]model.Event, error) {
	events := []model.EventSource{}
	ret := []model.Event{}

	// Lay du lieu tu event sourcing
	_, err := c.EventDB.Query(&events, `SELECT * FROM es.event_source WHERE aggregate_id = ? ORDER BY time`, aggregateID)
	if err != nil {
		return []model.Event{}, err
	}

	for _, event := range events {
		ev, err := Decode(event, model.AddTeacherEvent{})
		if err != nil {
			log.Println(err)
			return []model.Event{}, err
		}
		ret = append(ret, ev)
	}

	return ret, nil
}
