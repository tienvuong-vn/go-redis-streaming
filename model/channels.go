package model

import (
	"strings"

	"github.com/dustin/go-broadcast"
)

type Message struct {
	ChannelId string
	Text      string
}

type Listener struct {
	ChannelId string
	Chan      chan interface{}
}

type ChannelManager struct {
	channels     map[string]broadcast.Broadcaster
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message
	SseTotal     int64
	SseClosed    int64
	TotalMessage int64
	SseLive      int64
	WsTotal      int64
	WsClosed     int64
	WsLive       int64
}

func NewChannelManager() *ChannelManager {
	manager := &ChannelManager{
		channels:     make(map[string]broadcast.Broadcaster),
		open:         make(chan *Listener, 100),
		close:        make(chan *Listener, 100),
		delete:       make(chan string, 100),
		messages:     make(chan *Message, 100),
		SseTotal:     0,
		SseClosed:    0,
		TotalMessage: 0,
		SseLive:      0,
		WsTotal:      0,
		WsClosed:     0,
		WsLive:       0,
	}

	go manager.run()
	return manager
}

func (m *ChannelManager) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case channelId := <-m.delete:
			m.deleteBroadcast(channelId)
		case message := <-m.messages:
			m.TotalMessage += 1
			m.channel(message.ChannelId).Submit(message.Text)
		}
	}
}

func (m *ChannelManager) register(listener *Listener) {
	m.channel(listener.ChannelId).Register(listener.Chan)
}

func (m *ChannelManager) deregister(listener *Listener) {
	m.channel(listener.ChannelId).Unregister(listener.Chan)
	if !IsClosed(listener.Chan) {
		close(listener.Chan)
	}
}

func (m *ChannelManager) deleteBroadcast(channelId string) {
	b, ok := m.channels[channelId]
	if ok {
		b.Close()
		delete(m.channels, channelId)
	}
}

func (m *ChannelManager) channel(channelId string) broadcast.Broadcaster {
	b, ok := m.channels[channelId]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.channels[channelId] = b
	}
	return b
}

func (m *ChannelManager) OpenListener(path string, channelLink string) chan interface{} {
	// Each channel separate by ,
	s := strings.Split(channelLink, ",")
	// Add one listener for all channels
	listener := make(chan interface{})
	for i := 0; i < len(s); i++ {
		m.open <- &Listener{
			ChannelId: path + ":" + s[i],
			Chan:      listener,
		}
	}
	if path != "sse" && channelLink != "admin" {
		// Add Ping channel to every listener
		m.open <- &Listener{
			ChannelId: "PING",
			Chan:      listener,
		}
	}
	return listener
}

func (m *ChannelManager) CloseListener(path string, channelLink string, channel chan interface{}) {
	// Each channel separate by ,
	s := strings.Split(channelLink, ",")
	for i := 0; i < len(s); i++ {
		m.close <- &Listener{
			ChannelId: path + ":" + s[i],
			Chan:      channel,
		}
	}
	if path != "sse" && channelLink != "admin" {
		// Close Ping channel
		m.close <- &Listener{
			ChannelId: "PING",
			Chan:      channel,
		}
	}
}

func (m *ChannelManager) DeleteBroadcast(channelId string) {
	m.delete <- channelId
}

func (m *ChannelManager) Submit(channelId string, text string) {
	s := strings.Split(channelId, ",")
	// Send message to all listener
	for i := 0; i < len(s); i++ {
		msg := &Message{
			ChannelId: s[i],
			Text:      text,
		}
		m.messages <- msg
	}
}

func IsClosed(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
