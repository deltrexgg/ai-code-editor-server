package terminal

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TerminalHandler(w http.ResponseWriter, r *http.Request) {
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	cmd := exec.Command("bash")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("PTY start error:", err)
		return
	}
	defer func() {
		_ = ptmx.Close()
		_ = cmd.Process.Kill()
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		_, _ = ptmx.Write(msg)
	}
}