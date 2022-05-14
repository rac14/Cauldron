package handler

import (
	"Cauldron/login"
	"Cauldron/player"
	"Cauldron/proxy"
	"Cauldron/upstream"
	"encoding/json"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/google/uuid"
	"log"
)

type ServerListResponse struct {
	Version     ServerListVersion `json:"version"`
	Players     ServerListPlayers `json:"players"`
	Description chat.Message      `json:"description"`
	FavIcon     string            `json:"favicon,omitempty"`
}

type ServerListVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type ServerListPlayers struct {
	Max    int                      `json:"max"`
	Online int                      `json:"online"`
	Sample []ServerListPlayerEntity `json:"sample"`
}

type ServerListPlayerEntity struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

var defaultResponse = ""

var listPacket pk.Packet

func InitPing() {

	a := ServerListResponse{
		Version: ServerListVersion{
			Name:     "Cauldron Proxy",
			Protocol: 758,
		},
		Players: ServerListPlayers{
			Max:    10,
			Online: 0,
			Sample: []ServerListPlayerEntity{},
		},
		Description: chat.Message{
			Text:       "Cauldron Proxy",
			Color:      "#c00f33",
			UnderLined: true,
		},
	}
	b, _ := json.Marshal(a)
	defaultResponse = string(b)

	listPacket = pk.Marshal(0x00, pk.String(getListResponse()))
}

func HandleLogin(conn net.Conn) {
	// login, get player info
	info, err := login.AcceptLogin(conn)
	if err != nil {
		log.Print("Login failed")
		return
	}

	upstreams := upstream.ConnectToUpstreamServer("127.0.0.1:19132", info.Name)
	login.LoginSuccess(conn, info.Name, info.UUID)

	pl := &player.Player{
		Name:          info.Name,
		Uuid:          info.UUID,
		Conn:          &conn,
		Upstream:      upstreams,
		CurrentServer: "lobby",
	}
	proxy.Players[info.Name] = pl

	pl.Handle()
}

func HandleListPing(conn net.Conn) {
	var p pk.Packet
	for i := 0; i < 2; i++ {
		err := conn.ReadPacket(&p)
		if err != nil {
			return
		}

		switch p.ID {
		case 0x00:
			err = conn.WritePacket(listPacket)
		case 0x01:
			err = conn.WritePacket(p)
		}
		if err != nil {
			return
		}
	}
}

func getListResponse() string {
	return defaultResponse
}
