package smsproxy

import (
	"fmt"
	"sync"
)

type repository interface {
	update(id MessageID, newStatus MessageStatus) error
	save(id MessageID) error
	get(id MessageID) (MessageStatus, error)
}

type inMemoryRepository struct {
	db   map[MessageID]MessageStatus
	lock sync.RWMutex
}

func (r *inMemoryRepository) save(id MessageID) error {
	// save given MessageID with ACCEPTED status. If given MessageID already exists, return an error
	r.lock.Lock()
	if _, found := r.db[id]; found {
		return fmt.Errorf("already saved message : %s", id)
	}
	r.db[id] = Accepted
	r.lock.Unlock()
	return nil
}

func (r *inMemoryRepository) get(id MessageID) (MessageStatus, error) {
	// return status of given message, by it's MessageID. If not found, return NOT_FOUND status
	r.lock.RLock()
	status, found := r.db[id]
	r.lock.RUnlock()
	if found {
		return status, nil
	}
	return NotFound, fmt.Errorf("not found message : %s", id)
}

func (r *inMemoryRepository) update(id MessageID, newStatus MessageStatus) error {
	// Set new status for a given message.
	// If message is not in ACCEPTED state already - return an error.
	// If current status is FAILED or DELIVERED - don't update it and return an error. Those are final statuses and cannot be overwritten.
	r.lock.Lock()
	defer r.lock.Unlock()
	status, found := r.db[id]
	if !found {
		return fmt.Errorf("not found message : %s", id)
	}
	if status == Failed || status == Delivered {
		return fmt.Errorf("already updated message : %s", id)
	}
	if status == Accepted || status == Confirmed {
		r.db[id] = newStatus
	}
	return nil
}

func newRepository() repository {
	return &inMemoryRepository{db: make(map[MessageID]MessageStatus), lock: sync.RWMutex{}}
}
