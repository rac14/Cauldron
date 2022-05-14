package translator

import (
	javapacket "github.com/Tnze/go-mc/net/packet"
	bedrockpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var BedrockPacket = make(map[int32]bedrockpacket.Packet)

var JavaPacket = make(map[uint32]javapacket.Packet)

func initTranslator() {
}

func RegisterJavaToBe(id int32, packet bedrockpacket.Packet) {
	BedrockPacket[id] = packet
}

func RegisterBeToJava(id uint32, packet javapacket.Packet) {
	JavaPacket[id] = packet
}
