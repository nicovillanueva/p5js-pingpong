package server

import (
	"errors"
	"math/rand"
)

// DataStore stores stuff
type DataStore struct {
	Users   map[string]string
	Matches map[int64]*PingPongMatch
}

// NewDataStore spawns a new datastore
func NewDataStore() *DataStore {
	return &DataStore{
		Users:   make(map[string]string),
		Matches: make(map[int64]*PingPongMatch),
	}
}

// FindMatch returns a match based on the id
func (ds *DataStore) FindMatch(id int64) (*PingPongMatch, error) {
	m, ok := ds.Matches[id]
	if !ok {
		return &PingPongMatch{}, errors.New("match not found")
	}
	return m, nil
}

// SaveMatch returns the ID of the saved match
// Returns error in case of db save error
func (ds *DataStore) SaveMatch(match *PingPongMatch) error {
	i := rand.Int63()
	for {
		if _, exists := ds.Matches[i]; exists {
			logger.Error("collision generating match id: %v", i)
			i = rand.Int63()
			break
		}
		break
	}
	// logger.Printf("saving %s match of %+v with id %d", match.kind, match.players, i)
	match.matchID = i
	ds.Matches[i] = match
	return nil
}

// StoreUser saves a user
func (ds *DataStore) StoreUser(user, pass string) error {
	if ds.Users == nil {
		ds.Users = make(map[string]string, 0)
	}
	logger.Debugf("storing user: %s", user)
	if _, ok := ds.Users[user]; ok {
		return errors.New("user exists")
	}
	ds.Users[user] = pass
	return nil
}

// IsUserValid says if the user/pass combination is valid
func (ds *DataStore) IsUserValid(user, pass string) bool {
	if ds.Users == nil {
		logger.Errorf("tried to validate user with nil DataStore!")
		return false
	}
	logger.Debugf("validating user %s", user)
	if ds.Users[user] == pass {
		return true
	}
	return false
}
