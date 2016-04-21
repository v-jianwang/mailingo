package utils

import (
	"sync"
)


type State struct {
	name string
	items map[string][]byte
	Locker sync.Mutex
}

var (
	states []State
	mutx sync.Mutex
)


func Stated(name string) (*State, bool) {
	mutx.Lock()
	defer mutx.Unlock()
	
	found := true
	for _, st := range states {
		if st.name == name {
			return &st, found
		}
	}

	found = false
	st := State{
		name: name,
		items: make(map[string][]byte),
	}
	states = append(states, st)
	return &st, found
}


func (st *State) Item(key string, defaul []byte) ([]byte, bool) {
	b, exists := st.items[key]
	if !exists {
		st.items[key] = defaul
		b = defaul
	}
	return b, exists
}


func (st *State) SetItem(key string, b []byte) {
	st.items[key] = b
}
