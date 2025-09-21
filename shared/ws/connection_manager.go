package ws

import (
	"errors"
	"net/http"
	"sync"

	"ride-sharing/shared/contracts"

	"github.com/gorilla/websocket"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
	upgrader              = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}
)

type connWrapper struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

type ConnectionManager interface {
	InitUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
	Add(id string, con *websocket.Conn)
	Remove(id string)
	Get(id string) (*websocket.Conn, bool)
	Emit(id string, message contracts.WSMessage) error
}

type connectionManagerImpl struct {
	connectionsMap map[string]*connWrapper
	mutex          sync.RWMutex
}

func NewConnectionManager() ConnectionManager {
	return &connectionManagerImpl{
		connectionsMap: make(map[string]*connWrapper),
	}
}

func (cm *connectionManagerImpl) InitUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cm *connectionManagerImpl) Add(id string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connectionsMap[id] = &connWrapper{
		conn:  conn,
		mutex: sync.Mutex{},
	}
}

func (cm *connectionManagerImpl) Remove(id string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if wrapper, ok := cm.connectionsMap[id]; ok {
		_ = wrapper.conn.Close()
		delete(cm.connectionsMap, id)
	}
}

func (cm *connectionManagerImpl) Get(id string) (*websocket.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	wrapper, exists := cm.connectionsMap[id]
	if !exists {
		return nil, false
	}
	return wrapper.conn, true
}

func (cm *connectionManagerImpl) Emit(id string, message contracts.WSMessage) error {
	cm.mutex.RLock()
	wrapper, exists := cm.connectionsMap[id]
	cm.mutex.RUnlock()

	if !exists {
		return ErrConnectionNotFound
	}

	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()

	return wrapper.conn.WriteJSON(message)
}
