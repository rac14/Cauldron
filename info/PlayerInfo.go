package info

import "github.com/google/uuid"

type PlayerInfo struct {
	Name    string
	UUID    uuid.UUID
	OPLevel int
}
