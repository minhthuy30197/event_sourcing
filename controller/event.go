package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/minhthuy30197/event_sourcing/model"
)

func (c *Controller) SaveEvent(ev model.Event) error {
	tx, err := c.EventDB.Begin()
	if err != nil {
		return err
	}

	// Kiem tra version
	var version int32
	_, err = c.EventDB.Query(&version, `SELECT version FROM es.event_source WHERE aggregate_id = ? ORDER BY time DESC LIMIT 1`, ev.AggregateId)
	if err != nil {
		return err
	}
	if (version + 1) != ev.Version {
		return errors.New("Nội dung này đang được chỉnh sửa bởi một người khác. Vui lòng thử lại sau.")
	}

	// Encode
	eventDB, err := Encode(ev)
	if err != nil {
		log.Println("Encode loi")
		return err
	}

	// Insert EventDB
	err = tx.Insert(&eventDB)
	if err != nil {
		log.Println(err)
		log.Println("Insert loi")
		tx.Rollback()
		return err
	}

	// Apply event. Neu loi thi rollback
	err = ev.SaveReadDB(c.Config)
	if err != nil {
		tx.Rollback()
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
func Decode(event model.EventSource) (model.Event, error) {
	var err error
	ret := model.Event{}

	ret.AggregateId = event.AggregateId
	ret.EventType = event.EventType
	ret.Revision = event.Revision
	ret.Time = event.Time
	ret.UserID = event.UserID
	ret.Version = event.Version

	switch event.EventType {
	case "TeacherAdded":
		var tmp model.AddTeacherEvent
		err = json.Unmarshal([]byte(event.Data[0]), &tmp)
		if err != nil {
			fmt.Println("error:", err)
		}
		ret.Data = tmp
	case "TeacherRemoved":
		var tmp model.RemoveTeacherEvent
		err = json.Unmarshal([]byte(event.Data[0]), &tmp)
		if err != nil {
			fmt.Println("error:", err)
		}
		ret.Data = tmp
	}

	return ret, nil
}

// Events returns **All** the persisted events
func (c *Controller) Events(aggregateID string, startTime, endTime time.Time) ([]model.Event, error) {
	events := []model.EventSource{}
	ret := []model.Event{}

	// Lay du lieu tu event sourcing
	c.EventDB.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})

	_, err := c.EventDB.Query(&events, `SELECT * FROM es.event_source WHERE aggregate_id = ? AND time >= ? AND time <= ? ORDER BY time`,
		aggregateID, startTime, endTime)
	if err != nil {
		return []model.Event{}, err
	}
	for _, event := range events {
		ev, err := Decode(event)
		if err != nil {
			log.Println(err)
			return []model.Event{}, err
		}
		ret = append(ret, ev)
	}

	return ret, nil
}

func (c *Controller) EventsByVersion(aggregateID string, startVersion int32, endVersion int32) ([]model.Event, error) {
	ret := []model.Event{}
	
	return ret, nil
}

func (c *Controller) CreateSnapshot(aggregateID string, version int32, aggregate model.Aggregate) (model.Snapshot, error) {
	var newSnapshot, latestSnapshot model.Snapshot
	// Lay snapshot moi nhat
	_, err := c.DB.Query(&latestSnapshot, `
		SELECT * FROM es.snapshot
		WHERE aggregate_id = ? ORDER BY version DESC LIMIT 1`)
	if err != nil {
		return model.Snapshot{}, nil
	}

	// Lay danh sach event sau snapshot
	rs, err := c.EventsByVersion(aggregateID, latestSnapshot.Version, version)
	if err != nil {
		return model.Snapshot{}, nil
	}

	for _, event := range rs {
		aggregate.Apply(event)
	}

	newSnapshot.Version = version

	return newSnapshot, nil
}

func SaveSnapshot() error {
	return nil
}
