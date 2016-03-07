package listener

import (
	"log"
)

type Handler interface {
	Call(HubMessage) error
}

type Logger struct{}

func (l *Logger) Call(msg HubMessage) error {
	log.Printf("received message %v", msg)
	return nil
}

type Registry struct {
	entries []func(HubMessage) error
}

func (r *Registry) Add(h func(msg HubMessage) error) {
	r.entries = append(r.entries, h)
	return
}

func (r *Registry) Call(msg HubMessage) {
	for _, h := range r.entries {
		go h(msg)
	}
}

func MsgHandlers() Registry {
	var handlers Registry

	handlers.Add((&Logger{}).Call)
	handlers.Add((&Reloader{}).Call)
	//handlers.Add((&Mailgun{ServerConfig.Mailgun}).Call)

	return handlers
}
