package translator

import (
	"errors"
	javapacket "github.com/Tnze/go-mc/net/packet"
	bedrockpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TranslateJava(packet javapacket.Packet) (bedrockpacket.Packet, error) {
	if _, ok := BedrockPacket[packet.ID]; ok {
		return BedrockPacket[packet.ID], nil
	}
	// TODO: HACK Golang moment
	return &bedrockpacket.Login{
		ClientProtocol:    0,
		ConnectionRequest: nil,
	}, errors.New("Packet ID not found in translator")
}

func TranslateBedrock(packet bedrockpacket.Packet) (javapacket.Packet, error) {
	if _, ok := JavaPacket[packet.ID()]; ok {
		return JavaPacket[packet.ID()], nil
	}
	// TODO: HACK Golang moment
	return javapacket.Packet{ID: int32(0x1)}, errors.New("Packet ID not found in translator")
}

//func TranslateJava(packet javapacket.Packet) (bedrockpacket.Packet, error) {
//	//Hardcoded
//	switch packet.ID {
//	case Chat:
//		return bedrockpacket.Text{TextType: bedrockpacket.TextTypeChat, NeedsTranslation: false, SourceName: packet.}
//	}
//}
//
//func TranslateBedrock(packet bedrockpacket.Packet) (javapacket.Packet, error) {
//
//}
