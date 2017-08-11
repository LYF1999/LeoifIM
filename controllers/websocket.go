package controllers

import (
	"LeoifIM/models"
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	beego.Controller
}

func (this *WebSocketController) Get() {
	uName := this.GetSession("name")
	if uName == nil || uName == "" {
		this.Data["json"] = map[string]string{"detail": "no auth"}
		this.Ctx.ResponseWriter.WriteHeader(401)
	}
	this.ServeJSON()
}

func (this *WebSocketController) Join() {
	uName := this.GetString("name")
	uNameStr := uName
	if uName == "" {

	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	Join(uNameStr, ws)

	defer Leave(uNameStr)

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if string(p) == "" {
			continue
		}
		publish <- newEvent(models.MESSAGE, uNameStr, string(p))
	}
}

func broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if event.Type == models.JOIN {
		sub := FindUser(event.User)
		ws := sub.Conn
		events, _ := json.Marshal(models.GetEvents())
		ws.WriteMessage(websocket.TextMessage, events)
	}

	if event.Type == models.MESSAGE {
		models.NewArchive(event)
	}
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}
