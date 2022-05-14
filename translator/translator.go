package translator

import (
	javapacket "github.com/Tnze/go-mc/net/packet"
	bedrockpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TranslateJava(packet javapacket.Packet) bedrockpacket.Packet {
	return BedrockPacket[packet.ID]
}

func TranslateBedrock(packet bedrockpacket.Packet) javapacket.Packet {
	return JavaPacket[packet.ID()]
}
