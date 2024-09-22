package controllers

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"io"
	"log/slog"
	"sync"
)

type user struct {
	Id string
	Ws *websocket.Conn
	Mu sync.Mutex
}

var mu sync.RWMutex
var connected = make(map[string]*user)

func Signal(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		var userId string
		err := websocket.Message.Receive(ws, &userId)
		if err != nil {
			slog.Error("error receiving user id:", "error", err.Error())
			return
		}

		newUser := &user{
			Id: userId,
			Ws: ws,
		}

		mu.Lock()
		connected[userId] = newUser
		mu.Unlock()

		slog.Info("user connected", "userId", userId)
		for {
			var incoming JSON
			err := websocket.JSON.Receive(ws, &incoming)
			if err != nil {
				slog.Error("receive error", "error", err.Error())

				if err == io.EOF {
					break
				}

				continue
			}

			msgType, ok := incoming["type"].(string)
			if !ok {
				continue
			}
			to, ok := incoming["to"].(string)
			if !ok {
				continue
			}
			switch msgType {
			case "ice":
				mu.RLock()
				recipient := connected[to]
				mu.RUnlock()

				if recipient == nil {
					continue
				}

				candidate := incoming["candidate"].(map[string]interface{})
				recipient.Mu.Lock()
				websocket.JSON.Send(recipient.Ws, JSON{
					"type":      "ice",
					"from":      userId,
					"candidate": candidate,
				})
				recipient.Mu.Unlock()
			case "offer":
				mu.RLock()
				recipient := connected[to]
				mu.RUnlock()

				if recipient == nil {
					continue
				}

				offer := incoming["offer"].(map[string]interface{})
				recipient.Mu.Lock()
				websocket.JSON.Send(recipient.Ws, JSON{
					"type":  "offer",
					"from":  userId,
					"offer": offer,
				})
				recipient.Mu.Unlock()
			case "answer":
				mu.RLock()
				recipient := connected[to]
				mu.RUnlock()

				if recipient == nil {
					continue
				}

				answer := incoming["answer"].(map[string]interface{})
				recipient.Mu.Lock()
				websocket.JSON.Send(recipient.Ws, JSON{
					"type":   "answer",
					"from":   userId,
					"answer": answer,
				})
			}
		}

		slog.Info("user disconnected", "userId", userId)
		mu.Lock()
		delete(connected, userId)
		mu.Unlock()
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
