package upstream

import (
	"github.com/sandertv/gophertunnel/minecraft"
	bedrocklogin "github.com/sandertv/gophertunnel/minecraft/protocol/login"
	pk "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type UpstreamConn struct {
	ClientConn *minecraft.Conn
	/* ServerConn *net.Conn */
	name string
}

func ConnectToUpstreamServer(host string, name string) *UpstreamConn {

	conn, err := minecraft.Dialer{
		IdentityData: bedrocklogin.IdentityData{
			DisplayName: name,
		},
	}.Dial("raknet", host)

	if err != nil {
		panic(err)
	}

	err = conn.DoSpawn()

	if err != nil {
		panic(err)
	}

	//for {
	//
	//	packet, err := conn.ReadPacket()
	//	if err != nil {
	//		log.Printf("Error reading packet from server: %s\n", err.Error())
	//		return nil
	//	}
	//
	//	log.Printf("Recieved packet %d from upstream server", packet.ID)
	//
	//}

	return &UpstreamConn{
		name:       name,
		ClientConn: conn,
	}
}

func (u *UpstreamConn) ReadPacket() pk.Packet {
	p, err := u.ClientConn.ReadPacket()
	if err != nil {
		return nil
	}
	return p
}
