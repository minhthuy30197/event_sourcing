package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/minhthuy30197/event_sourcing/constant"
	"github.com/minhthuy30197/event_sourcing/model"
)

func (c *Controller) SaveEvent(ev model.Event, agg model.Aggregate) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	// Check version
	var version int32
	_, err = c.DB.Query(&version, `SELECT version FROM es.event_source WHERE aggregate_id = ? ORDER BY version DESC LIMIT 1`, ev.AggregateId)
	if err != nil {
		return err
	}
	if (version + 1) != ev.Version {
		return errors.New("Nội dung này đang được chỉnh sửa bởi một người khác. Vui lòng thử lại sau.")
	}

	// Encode
	eventDB, err := Encode(ev)
	if err != nil {
		log.Println("Encode lỗi")
		return err
	}

	// Insert EventDB
	err = tx.Insert(&eventDB)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Tao snapshot
	if ev.Version % constant.CountVersionPerSnapshot == 0 {
		snapshot, err := c.CreateSnapshot(agg, ev)
		if err != nil {
			tx.Rollback()
			log.Println("Lỗi khi tạo snapshot")
			return err
		}

		err = tx.Insert(&snapshot)
		if err != nil {
			tx.Rollback()
			log.Println("Lỗi khi insert snapshot")
			return err
		}
	}

	// Apply event. Neu loi thi rollback
	err = ev.SaveReadDB(c.Config)
	if err != nil {
		tx.Rollback()
		log.Println("Lỗi khi update bảng read")
		return err
	}

	tx.Commit()
	return nil
}

func BuildBaseEvent(aggregateID, userID, eventType string, data model.EventInterface, version int32) model.Event {
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

	_, err := c.DB.Query(&events, `SELECT * FROM es.event_source WHERE aggregate_id = ? AND time >= ? AND time <= ? ORDER BY time`,
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
	events := []model.EventSource{}
	ret := []model.Event{}

	var query string
	if endVersion != -1 {
		query = `SELECT * FROM es.event_source WHERE aggregate_id = ? AND version > ? AND version < ? ORDER BY version ASC`
	} else {
		query = `SELECT * FROM es.event_source WHERE aggregate_id = ? AND version > ? ORDER BY version ASC`
	}
	_, err := c.DB.Query(&events, query, aggregateID, startVersion, endVersion)
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

func (c *Controller) GetLatestSnapshot(aggregateID string, agg model.Aggregate) (error, int32) {
	var latestSnapshot model.Snapshot
	_, err := c.DB.Query(&latestSnapshot, `
		SELECT * FROM es.snapshot
		WHERE aggregate_id = ? ORDER BY version DESC LIMIT 1`, aggregateID)
	if err != nil {
		log.Println(err.Error())
		return err, 0
	}

	if latestSnapshot.Version > 0 {
		err = json.Unmarshal([]byte(latestSnapshot.Data[0]), &agg)
		if err != nil {
			fmt.Println("error:", err)
			return err, 0
		}
	}

	return nil, latestSnapshot.Version
}

func (c *Controller) CreateSnapshot(agg model.Aggregate, event model.Event) (model.Snapshot, error) {
	var newSnapshot model.Snapshot

	// Lay snapshot moi nhat
	err, latestVersion := c.GetLatestSnapshot(event.AggregateId, agg)
	if err != nil {
		return model.Snapshot{}, err
	}

	// Lay danh sach event sau snapshot
	rs, err := c.EventsByVersion(event.AggregateId, latestVersion, event.Version)
	if err != nil {
		return model.Snapshot{}, err
	}
	for _, ev := range rs {
		ev.Apply(agg)
	}

	// Apply event moi them
	event.Apply(agg)

	// Tao snapshot moi
	data, err := json.Marshal(agg)
	if err != nil {
		return newSnapshot, err
	}
	newSnapshot.Version = event.Version
	newSnapshot.AggregateId = event.AggregateId
	newSnapshot.Time = time.Now()
	newSnapshot.Data = []string{string(data)}

	return newSnapshot, nil
}
