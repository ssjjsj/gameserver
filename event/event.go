package event

import (
)

type EventData interface{

}
type EventHandler func (data EventData)
eventMap dict[int, EventHandler]

func AddEventListener(eventType int, handler EventHandler) {

}


func DispatchEvent(eventType int, )