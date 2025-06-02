package handler

import (
	"log"
	"net/http"
	"sync"
	"vintage-vision-api/internal/model"
	"vintage-vision-api/internal/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var partyConnections = make(map[string][]*websocket.Conn)
var partyMu sync.RWMutex

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	partyID := r.URL.Query().Get("party_id")
	if partyID == "" {
		http.Error(w, "party_id is required", http.StatusBadRequest)
		return
	}

	senderID := r.URL.Query().Get("sender_id")
	if senderID == "" {
		http.Error(w, "sender_id (token) is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	partyMu.Lock()
	partyConnections[partyID] = append(partyConnections[partyID], conn)
	partyMu.Unlock()
	log.Printf("Nueva conexión a party %s\n", partyID)

	for {
		var msg model.WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			utils.Logger.Errorf("Error leyendo mensaje de WebSocket: %v", err)
			break
		}

		if err := utils.ValidateWebSocketMessage(msg); err != nil {
			utils.Logger.Warnf("Mensaje inválido: %v", err)
			continue
		}

		broadcastToParty(partyID, msg, conn)
	}

	removeConnection(partyID, conn)
	conn.Close()
	utils.Logger.Infof("Conexión cerrada de party %s", partyID)
}

func broadcastToParty(partyID string, msg model.WebSocketMessage, sender *websocket.Conn) {
	partyMu.RLock()
	conns := partyConnections[partyID]
	partyMu.RUnlock()

	for _, conn := range conns {
		if conn != sender {
			err := conn.WriteJSON(msg)
			if err != nil {
				utils.Logger.Errorf("Error enviando mensaje a party %s: %v", partyID, err)
			}
		}
	}
}

func removeConnection(partyID string, conn *websocket.Conn) {
	partyMu.Lock()
	defer partyMu.Unlock()

	conns := partyConnections[partyID]
	for i, c := range conns {
		if c == conn {
			partyConnections[partyID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	if len(partyConnections[partyID]) == 0 {
		delete(partyConnections, partyID)
		utils.Logger.Infof("Party %s eliminada (sin conexiones)", partyID)
	}
}
