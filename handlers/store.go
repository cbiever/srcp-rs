package handlers

import (
	"srcp-rs/model"
	"srcp-rs/srcp"
)

type Store interface {
	SetSrcpEndpoint(string)
	GetSrcpEndpoint() string
	SaveConnection(int, srcp.SrcpConnection)
	GetConnection(int) srcp.SrcpConnection
	CreateGL(int, int) *model.GeneralLoco
	SaveGL(int, int, *model.GeneralLoco)
	GetGL(int, int) *model.GeneralLoco
	DeleteGL(int, int)
	GetGLS() map[int]map[int]*model.GeneralLoco
}

type MemoryStore struct {
	srcpEndpoint    string
	srcpConnections map[int]srcp.SrcpConnection
	gls             map[int]map[int]*model.GeneralLoco
}

var store Store

func init() {
	store = new(MemoryStore)
}

func GetStore() Store {
	return store
}

func (store *MemoryStore) SetSrcpEndpoint(srcpEndpoint string) {
	store.srcpEndpoint = srcpEndpoint
}

func (store *MemoryStore) GetSrcpEndpoint() string {
	return store.srcpEndpoint
}

func (store *MemoryStore) SaveConnection(sessionID int, connection srcp.SrcpConnection) {
	if store.srcpConnections == nil {
		store.srcpConnections = make(map[int]srcp.SrcpConnection)
	}
	store.srcpConnections[sessionID] = connection
}

func (store *MemoryStore) GetConnection(sessionID int) srcp.SrcpConnection {
	return store.srcpConnections[sessionID]
}

func (store *MemoryStore) CreateGL(bus int, address int) *model.GeneralLoco {
	gl := new(model.GeneralLoco)
	store.SaveGL(bus, address, gl)
	return gl
}

func (store *MemoryStore) SaveGL(bus int, address int, gl *model.GeneralLoco) {
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

func (store *MemoryStore) GetGL(bus int, address int) *model.GeneralLoco {
	if store.gls != nil && store.gls[bus] != nil {
		return store.gls[bus][address]
	} else {
		return nil
	}
}

func (store *MemoryStore) DeleteGL(bus int, address int) {
	if store.gls != nil && store.gls[bus] != nil {
		store.gls[bus][address] = nil
	}
}

func (store *MemoryStore) GetGLS() map[int]map[int]*model.GeneralLoco {
	return store.gls
}
