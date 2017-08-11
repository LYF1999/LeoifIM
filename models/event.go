package models

import (
	"container/list"

	"fmt"
)

const JOIN = 1
const LEAVE = 2
const MESSAGE = 3
const archiveSize = 20

var archive = list.New()

type Event struct {
	Type      int
	User      string
	Content   string
	Timestamp int
}

func NewArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

func GetEvents() []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		events = append(events, e)
		fmt.Println(e)
	}
	fmt.Println(events)
	return events
}
