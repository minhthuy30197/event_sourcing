package controller

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"git.hocngay.com/hocngay/event-sourcing/model"
)

var eventRegistry = map[string]reflect.Type{}

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
		return err
	}

	// TODO: Apply event. Neu loi thi rollback
	err = ev.Apply(c.Config)
	if err != nil {
		tx.Rollback()
		log.Println("Loi khi update bang read")
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

func Register(events ...model.EventInterface) {
	for _, event := range events {
		eventRegistry[event.AggregateType()] = reflect.TypeOf(event)
	}
	log.Println(eventRegistry)
}

// Decode return a deserialized event, ready to user
func Decode(event model.EventSource) (model.Event, error) {
	// deserialize json
	var err error
	ret := model.Event{}

	log.Println("-------------- den day roi 1") 

	// reflexion magic
	dataPointer := reflect.New(eventRegistry["AddTeacherEvent"])
	log.Println("-------------- den day roi 2")
	dataValue := dataPointer.Elem()
	var data map[string]interface{}
	log.Println("-------------- den day roi 3")
	array := []byte(event.Data[0])
	err = json.Unmarshal(array, &data)
	if err != nil {
		return model.Event{}, err
	}

	log.Println("-------------- den day roi 4")
	n := dataValue.NumField()
	for i := 0; i < n; i++ {
		field := dataValue.Type().Field(i)
		jsonName := field.Tag.Get("json")
		if jsonName == "" {
			jsonName = field.Name
		}
		log.Println(jsonName)
		val := dataValue.FieldByName(field.Name)
		val.Set(reflect.ValueOf(data[jsonName]))
		log.Println(data[jsonName])
	}
	log.Println("-------------- den day roi 5")
	ret.Time = event.Time
	ret.AggregateId = event.AggregateId
	ret.EventType = event.EventType
	ret.Version = event.Version
	ret.Revision = event.Revision
	ret.UserID = event.UserID

	dataInterface := dataValue.Interface()
	ret.Data = dataInterface

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
		ev, err := Decode(event)
		if err != nil {
			log.Println(err)
			return []model.Event{}, err
		}
		ret = append(ret, ev)
	}

	return ret, nil
}