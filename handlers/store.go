package handlers

import (
	"srcp-rs/model"
	"srcp-rs/srcp"
)

type Store struct {
	srcpEndpoint    string
	srcpConnections map[int]*srcp.SrcpConnection
	gls             map[int]map[int]*model.GeneralLoco
}

var store Store

func GetStore() *Store {
	return &store
}

func (store *Store) SetSrcpEndpoint(srcpEndpoint string) {
	store.srcpEndpoint = srcpEndpoint
}

func (store *Store) GetSrcpEndpoint() string {
	return store.srcpEndpoint
}

func (store *Store) SaveConnection(sessionID int, connection *srcp.SrcpConnection) {
	if store.srcpConnections == nil {
		store.srcpConnections = make(map[int]*srcp.SrcpConnection)
	}
	store.srcpConnections[sessionID] = connection
}

func (store *Store) GetConnection(sessionID int) *srcp.SrcpConnection {
	return store.srcpConnections[sessionID]
}

func (store *Store) GetGL(bus int, address int) *model.GeneralLoco {
	if store.gls != nil && store.gls[bus] != nil {
		return store.gls[bus][address]
	} else {
		return nil
	}
}

func (store *Store) CreateGL(bus int, address int) *model.GeneralLoco {
	gl := new(model.GeneralLoco)
	store.SaveGL(bus, address, gl)
	return gl
}

func (store *Store) SaveGL(bus int, address int, gl *model.GeneralLoco) {
	if store.gls == nil {
		store.gls = make(map[int]map[int]*model.GeneralLoco)
	}
	if store.gls[bus] == nil {
		store.gls[bus] = make(map[int]*model.GeneralLoco)
	}
	gl.Bus = bus
	gl.Address = address
	store.gls[bus][address] = gl
}

func (store *Store) DeleteGL(bus int, address int) {
	if store.gls != nil && store.gls[bus] != nil {
		store.gls[bus][address] = nil
	}
}

func (store *Store) GetGLS() map[int]map[int]*model.GeneralLoco {
	return store.gls
}
