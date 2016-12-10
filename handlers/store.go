package handlers

import (
	"srcp-rs/srcp"
)

type Store struct {
	connections map[int]*srcp.SrcpConnection
	gls         map[int]map[int]*srcp.GeneralLoco
}

var store Store

func (store *Store) SaveConnection(sessionID int, connection *srcp.SrcpConnection) {
	if store.connections == nil {
		store.connections = make(map[int]*srcp.SrcpConnection)
	}
	store.connections[sessionID] = connection
}

func (store *Store) GetConnection(sessionID int) *srcp.SrcpConnection {
	return store.connections[sessionID]
}

func (store *Store) GetGL(bus int, address int) *srcp.GeneralLoco {
	if store.gls != nil && store.gls[bus] != nil {
		return store.gls[bus][address]
	} else {
		return nil
	}
}

func (store *Store) CreateGL(bus int, address int) *srcp.GeneralLoco {
	if store.gls == nil {
		store.gls = make(map[int]map[int]*srcp.GeneralLoco)
	}
	if store.gls[bus] == nil {
		store.gls[bus] = make(map[int]*srcp.GeneralLoco)
	}
	store.gls[bus][address] = new(srcp.GeneralLoco)
	return store.gls[bus][address]
}

func (store *Store) DeleteGL(bus int, address int) {
	if store.gls != nil && store.gls[bus] != nil {
		store.gls[bus][address] = nil
	}
}

func (store *Store) GetGLS() map[int]map[int]*srcp.GeneralLoco {
	return store.gls;
}
