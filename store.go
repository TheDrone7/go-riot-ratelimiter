package ratelimiter

// Implements basic Cache/Map behavior
type Store struct {
	data map[string]any
}

// Creates a new Store instance
func NewStore() *Store {
	return &Store{
		data: make(map[string]any),
	}
}

// Stores a key-value pair in the store
func (s *Store) Set(key string, value any) {
	s.data[key] = value
}

// Retrieves a value by key from the store
// Returns the value and a boolean indicating if the key was found
func (s *Store) Get(key string) (any, bool) {
	value, exists := s.data[key]
	return value, exists
}

// Checks if a key exists in the store
func (s *Store) Has(key string) bool {
	_, exists := s.data[key]
	return exists
}

// Removes a key-value pair from the store
// Returns true if the key was found and removed, false otherwise
func (s *Store) Remove(key string) bool {
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return true
	}
	return false
}

// Returns the number of key-value pairs in the store
func (s *Store) Size() int {
	return len(s.data)
}

// Clears all key-value pairs from the store
func (s *Store) Clear() {
	s.data = make(map[string]any)
}
