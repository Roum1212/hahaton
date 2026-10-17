package updates

import (
	"sync"

	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Hub struct {
	mu          sync.Mutex
	subscribers map[model_user_id.UserID]map[chan struct{}]struct{}
}

func (x *Hub) Notify(users []model_user_id.UserID) {
	x.mu.Lock()
	defer x.mu.Unlock()

	for _, userID := range users {
		for channel := range x.subscribers[userID] {
			// Signals invalidate current data; one pending signal covers any number of changes.
			select {
			case channel <- struct{}{}:
			default:
			}
		}
	}
}

func (x *Hub) NotifyAll() {
	x.mu.Lock()
	defer x.mu.Unlock()

	for _, subscribers := range x.subscribers {
		for channel := range subscribers {
			select {
			case channel <- struct{}{}:
			default:
			}
		}
	}
}

func (x *Hub) Subscribe(userID model_user_id.UserID) (<-chan struct{}, func()) {
	x.mu.Lock()
	defer x.mu.Unlock()

	channel := make(chan struct{}, 1)
	if x.subscribers[userID] == nil {
		x.subscribers[userID] = make(map[chan struct{}]struct{})
	}

	x.subscribers[userID][channel] = struct{}{}

	return channel, func() {
		x.mu.Lock()
		defer x.mu.Unlock()

		delete(x.subscribers[userID], channel)

		if len(x.subscribers[userID]) == 0 {
			delete(x.subscribers, userID)
		}
	}
}

func NewHub() *Hub {
	return &Hub{subscribers: make(map[model_user_id.UserID]map[chan struct{}]struct{})}
}
