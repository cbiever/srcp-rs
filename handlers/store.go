package handlers

import (
	"srcp-rs/srcp"
)

type Store struct {
	store map[int]map[int]*srcp.GeneralLoco
}

func (store *Store) Get(bus int, address int) *srcp.GeneralLoco {
	if store.store != nil && store.store[bus] != nil {
		return store.store[bus][address]
	} else {
		return nil
	}
}

func (store *Store) Create(bus int, address int) *srcp.GeneralLoco {
	if store.store == nil {
		store.store = make(map[int]map[int]*srcp.GeneralLoco)
	}
	if store.store[bus] == nil {
		store.store[bus] = make(map[int]*srcp.GeneralLoco)
	}
	store.store[bus][address] = new(srcp.GeneralLoco)
	return store.store[bus][address]
}

func (store *Store) Delete(bus int, address int) {
	if store.store != nil && store.store[bus] != nil {
		store.store[bus][address] = nil
	}
}
