package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//check the origin matches the frontendURL
		//for now allow all
		return true
	},
}

type TerminalHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewTerminalHandler(authService service.AuthService, userService service.UserService) *TerminalHandler {
	return &TerminalHandler{authService: authService, userService: userService}
}

func (h *TerminalHandler) RegisterRoutes(router *mux.Router) {
	terminalRouter := router.PathPrefix("/projects/{projectId}/terminal").Subrouter()

	//we dont use WithJWTAuth since we're gonna accept a ticket from wsticket in auth service and exchange with userID
	terminalRouter.HandleFunc("", h.handleTerminal).Methods("GET")

}

func (h *TerminalHandler) handleTerminal(w http.ResponseWriter, r *http.Request) {
	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		response.WriteError(w, errors.NewError("unauthorized", http.StatusUnauthorized))
		return
	}
	log.Printf("ticket: %s", ticket)

	userId, err := h.authService.ExchangeWSTicket(r.Context(), ticket)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userId)
	if err != nil {
		response.WriteError(w, errors.NewError("user not found", http.StatusUnauthorized))
		return
	}

	log.Printf("Terminal opened for user %d", user.ID)

	//upgrade HTTP connection to Websocket, this is the handshake of HTTP becomes websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade to websocket: %v", err)
	}
	defer conn.Close()

	//spawn a bash process with a PTY
	cmd := exec.Command("/bin/bash")

	//set environment variables for the shell
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color", //tell programs we support colors
		"HOME=/tmp",           //sandbox the home directory
	)

	//start the command with a PTY attached
	//ptmx is the "master" end of the PTY, we read/write through this
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("failed to start pty: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("failed to start terminal\r\n"))
		return
	}
	defer ptmx.Close()
	defer cmd.Process.Kill() //kill bash when connection closes

	//goroutine to read from PTY (bash output) then send to websocket (browser)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("pty read error: %v", err)
				}

				return
			}
			//send bash's output to the browser
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	//main loop: read from Websocket (browser input) then write to PTY (bash stdin)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			//connection closed
			log.Printf("websocket closed for user %d", user.ID)
			return
		}

		var resizeMsg struct {
			Type string `json:"type"`
			Cols uint16 `json:"cols"`
			Rows uint16 `json:"rows"`
		}

		if json.Unmarshal(msg, &resizeMsg) == nil && resizeMsg.Type == "resize" {
			pty.Setsize(ptmx, &pty.Winsize{
				Cols: resizeMsg.Cols,
				Rows: resizeMsg.Rows,
			})

			continue
		}

		//send browser's keystrokes to bash
		if _, err := ptmx.Write(msg); err != nil {
			return
		}
	}
}
