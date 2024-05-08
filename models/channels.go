package models

import "github.com/dustin/go-broadcast"

var communityChannels = make(map[string]broadcast.Broadcaster)

func OpenCommunityListener(id string) chan interface{} {
	listener := make(chan interface{})
	GetCommunityBroadcast(id).Register(listener)
	return listener
}
func CloseCommunityListener(id string, listener chan interface{}) {
	GetCommunityBroadcast(id).Unregister(listener)
	close(listener)
}
func DeleteCommunityBroadcast(id string) {
	b, ok := communityChannels[id]
	if ok {
		b.Close()
		delete(communityChannels, id)
	}
}
func GetCommunityBroadcast(id string) broadcast.Broadcaster {
	b, ok := communityChannels[id]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		communityChannels[id] = b
	}
	return b
}

var postChannels = make(map[string]broadcast.Broadcaster)
func OpenPostListener(id string) chan interface{} {
	listener := make(chan interface{})
	GetPostBroadcast(id).Register(listener)
	return listener
}
func ClosePostListener(id string, listener chan interface{}) {
	GetPostBroadcast(id).Unregister(listener)
	close(listener)
}
func DeletePostBroadcast(id string) {
	b, ok := postChannels[id]
	if ok {
		b.Close()
		delete(postChannels, id)
	}
}
func GetPostBroadcast(id string) broadcast.Broadcaster {
	b, ok := postChannels[id]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		postChannels[id] = b
	}
	return b
}

var notificationChannels = make(map[string]broadcast.Broadcaster)
func OpenNotificationListener(id string) chan interface{} {
	listener := make(chan interface{})
	GetNotificationBroadcast(id).Register(listener)
	return listener
}
func CloseNotificationListener(id string, listener chan interface{}) {
	GetPostBroadcast(id).Unregister(listener)
	close(listener)
}
func DeleteNotificationBroadcast(id string) {
	b, ok := notificationChannels[id]
	if ok {
		b.Close()
		delete(notificationChannels, id)
	}
}
func GetNotificationBroadcast(id string) broadcast.Broadcaster {
	b, ok := notificationChannels[id]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		notificationChannels[id] = b
	}
	return b
}
