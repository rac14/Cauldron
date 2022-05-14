package handler

import (
	"Cauldron/login"
	"Cauldron/player"
	"Cauldron/proxy"
	"Cauldron/upstream"
	_ "embed"
	"encoding/json"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/nbt"
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
	handleJoin(conn)

	if err := conn.WritePacket(pk.Marshal(0x38,
		// https://wiki.vg/index.php?title=Protocol&oldid=16067#Player_Position_And_Look_.28clientbound.29
		pk.Double(0), pk.Double(0), pk.Double(0), // XYZ
		pk.Float(0), pk.Float(0), // Yaw Pitch
		pk.Byte(0),        // flag
		pk.VarInt(0),      // TP ID
		pk.Boolean(false), // Dismount vehicle
	)); err != nil {
		log.Printf("Login failed on sending PlayerPositionAndLookClientbound: %v", err)
		return
	}

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

//go:embed resources/DimensionCodec.snbt
var dimensionCodecSNBT string

//go:embed resources/Dimension.snbt
var dimensionSNBT string

func handleJoin(conn net.Conn) error {
	return conn.WritePacket(pk.Marshal(0x26,
		pk.Int(0),                                          // EntityID
		pk.Boolean(false),                                  // Is hardcore
		pk.UnsignedByte(1),                                 // Gamemode
		pk.Byte(1),                                         // Previous Gamemode
		pk.Array([]pk.Identifier{"world"}),                 // World Names
		pk.NBT(nbt.StringifiedMessage(dimensionCodecSNBT)), // Dimension codec
		pk.NBT(nbt.StringifiedMessage(dimensionSNBT)),      // Dimension
		pk.Identifier("world"),                             // World Name
		pk.Long(0),                                         // Hashed Seed
		pk.VarInt(100),                                     // Max Players
		pk.VarInt(15),                                      // View Distance
		pk.Boolean(false),                                  // Reduced Debug Info
		pk.Boolean(true),                                   // Enable respawn screen
		pk.Boolean(false),                                  // Is Debug
		pk.Boolean(true),                                   // Is Flat
	))
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
