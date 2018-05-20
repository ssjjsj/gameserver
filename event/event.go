package event

import (
)

type EventData interface{

}
type EventHandler func (data EventData)
var eventMap = make(map[string][]EventHandler)

func AddEventListener(eventType string, handler EventHandler) {
	handlerList, exits := eventMap[eventType] 
	if exits == false{
		eventMap[eventType] = make([]EventHandler, 0)
		handlerList = eventMap[eventType]
	}

	eventMap[eventType] = append(handlerList, handler)
}


func DispatchEvent(eventType string, data EventData){
	handlerList, exits := eventMap[eventType]
	if (exits){
		for i:=0; i<len(handlerList); i++ {
			handler := handlerList[i]
			handler(data)
		}
	}
}