package service

import "errors"

type inMemoryRepository struct {
	states map[string]reality
}

func newInMemoryRepository() (repo *inMemoryRepository) {
	repo = &inMemoryRepository{}
	repo.states = make(map[string]reality)
	return
}

func (repo *inMemoryRepository) updateReality(gameID string, newReality reality) (err error) {
	repo.states[gameID] = newReality
	return
}

func (repo *inMemoryRepository) getReality(gameID string) (gameReality reality, err error) {
	gameReality, ok := repo.states[gameID]
	if !ok {
		err = errors.New("No such game found.")
	}
	return
}
