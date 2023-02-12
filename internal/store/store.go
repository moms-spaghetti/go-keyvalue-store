package store

import (
	"errors"
	"fmt"
	"sync"
	"task1/internal/logger"
)

var (
	ErrStoreEmpty       = errors.New("store is empty")
	ErrStoreKeyNotFound = errors.New("key not found in store")
	ErrKeyEmpty         = errors.New("key cannot be empty")
)

type StoreData map[string]interface{}

type Storage struct {
	store   StoreData
	logger  logger.Logger
	rwMutex *sync.RWMutex
}

func NewStorage(logger *logger.Logger) *Storage {
	store := make(StoreData)
	rwMutex := &sync.RWMutex{}

	return &Storage{
		store:   store,
		logger:  *logger,
		rwMutex: rwMutex,
	}
}

func (s *Storage) Get(key string) (interface{}, error) {
	if key == "" {
		s.logger.Log(ErrKeyEmpty.Error())

		return nil, ErrKeyEmpty
	}

	s.logger.Log("get store access")

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(s.store) == 0 {
		return nil, ErrStoreEmpty
	}

	value, ok := s.store[key]
	if !ok {
		return nil, ErrStoreKeyNotFound
	}

	return value, nil
}

func (s *Storage) Post(data StoreData) error {
	keys := make([]string, len(data))
	index := 0

	for key := range data {
		if key == "" {
			s.logger.Log(ErrKeyEmpty.Error())

			return ErrKeyEmpty
		}
		keys[index] = key
		index++
	}

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	for key, value := range data {
		s.store[key] = value

		s.logger.Log(fmt.Sprintf("key: %s, value: %v - added to store", key, value))
	}

	return nil
}

func (s *Storage) Delete(key string) error {
	if key == "" {
		s.logger.Log(ErrKeyEmpty.Error())

		return ErrKeyEmpty
	}

	s.logger.Log("delete store access")

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(s.store) == 0 {
		return ErrStoreEmpty
	}

	_, ok := s.store[key]
	if !ok {
		return ErrStoreKeyNotFound
	}

	delete(s.store, key)

	return nil
}
