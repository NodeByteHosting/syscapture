package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/sirupsen/logrus"
)

var data = []byte("SysMon Uptime and Metrics Monitoring Agent")
var interval = 2 * time.Second

var upgrader websocket.Upgrader

// InitWebSocket initializes the WebSocket upgrader with the given configuration.
func InitWebSocket(cfg *config.Config) {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
			origin := r.Header.Get("Origin")

			if len(allowedOrigins) == 1 {
				return allowedOrigins[0] == "*"
			}

			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == origin {
					return true
				}
			}

			return false
		},
	}
}

// WebSocket handles WebSocket connections and sends data at regular intervals.
func WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("Failed to set websocket upgrade: %v", err)
		return
	}

	defer conn.Close()
	for {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			logrus.Errorf("Failed to write message: %v", err)
			break
		}
		time.Sleep(interval)
	}
}
