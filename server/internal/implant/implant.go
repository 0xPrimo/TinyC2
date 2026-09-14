package implant

import (
	"time"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/store"
)

type Implant struct {
	ID       string
	Meta     map[string]any
	tasks    []Task
	channels *store.Store[string, *Channel]
	seen     time.Time
}

func NewImplant(id string, meta map[string]any) *Implant {
	var (
		implant = Implant{
			ID:       id,
			channels: store.NewStore[string, *Channel](),
			Meta:     make(map[string]any),
		}
	)

	implant.MetaUpdate(meta)

	return &implant
}

func (i *Implant) Alive() bool {
	return time.Since(i.seen) <= 5*time.Second
}

func (i *Implant) TimerUpdate() {
	i.seen = time.Now()
}

func (i *Implant) TaskAdd(task Task) {
	i.tasks = append(i.tasks, task)
}

func (i *Implant) TaskPopAll() []Task {
	tasks := i.tasks
	i.tasks = []Task{}
	return tasks
}

func (i *Implant) ChannelAdd(name string, channel *Channel) {
	i.channels.Set(name, channel)
}

func (i *Implant) ChannelRemove(name string) {
	i.channels.Delete(name)
}

func (i *Implant) ChannelCurrent() string {
	var channelName string

	i.channels.ForEach(func(name string, channel *Channel) {
		if channel.InUse == true {
			channelName = channel.Name
		}
	})

	return channelName
}

func (i *Implant) ChannelList() []Channel {
	var channels []Channel

	i.channels.ForEach(func(name string, channel *Channel) {
		channels = append(channels, *channel)
	})

	return channels
}

func (i *Implant) MetaUpdate(meta map[string]any) {
	if val, ok := meta["host"]; ok {
		if host, isString := val.(string); isString {
			i.Meta["host"] = host
		}
	}

	if val, ok := meta["user"]; ok {
		if user, isString := val.(string); isString {
			i.Meta["user"] = user
		}
	}

	if val, ok := meta["domain"]; ok {
		if domain, isString := val.(string); isString {
			i.Meta["domain"] = domain
		}
	}

	if val, ok := meta["pid"]; ok {
		if pid, isString := val.(string); isString {
			i.Meta["pid"] = pid
		}
	}

	if val, ok := meta["os"]; ok {
		if os, isString := val.(string); isString {
			i.Meta["os"] = os
		}
	}
}
