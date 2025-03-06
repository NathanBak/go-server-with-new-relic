package storage

import "sync"

// MapStorage can be used to store string keys and any type of value.
type MapStorage[T any] struct {
	data map[string]T
	lck  sync.RWMutex
}

// Set will store the specified key/value combination.
func (s *MapStorage[T]) Set(key string, val T) error {
	s.lck.Lock()
	defer s.lck.Unlock()
	if s.data == nil {
		s.data = make(map[string]T)
	}
	s.data[key] = val
	return nil
}

// Get returns the value for the specified key.
func (s *MapStorage[T]) Get(key string) (v T, ok bool, err error) {
	if s.data != nil {
		s.lck.RLock()
		defer s.lck.RUnlock()
		v, ok = s.data[key]
	}
	return
}

// Delete will delete the specified key and its associated value from storage.  Returns true if the
// key existed (and was deleted) and false if the key did not exist.  If deleting many elements from
// storage, the Shrink() function may be called afterwards to free up storage by reducing the size
// of the map.
func (s *MapStorage[T]) Delete(key string) (v T, ok bool, err error) {
	if s.data != nil {
		s.lck.Lock()
		defer s.lck.Unlock()

		if v, ok = s.data[key]; ok {
			delete(s.data, key)
		}
	}
	return
}

// Keys returns a slice of all the keys in storage.
func (s *MapStorage[T]) Keys() ([]string, error) {
	list := []string{}

	if s.data != nil {
		s.lck.RLock()
		defer s.lck.RUnlock()
		for k := range s.data {
			list = append(list, k)
		}
	}

	return list, nil
}

// Shrink will free up memory if many elements have been removed from the map.  This function blocks
// and, depending on the amount of data stored, may take awhile to complete.
func (s *MapStorage[T]) Shrink() {
	if s.data == nil {
		return
	}

	newData := make(map[string]T)

	s.lck.Lock()
	defer s.lck.Unlock()

	for k, v := range s.data {
		newData[k] = v
	}

	s.data = newData
}
