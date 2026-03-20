package models

import "github.com/juevigrace/diva-server/storage/db"

type ChatType int

const (
	CHAT_DIRECT ChatType = iota
	CHAT_GROUP
	CHAT_UNSPECIFIED
)

func (c ChatType) String() string {
	switch c {
	case CHAT_DIRECT:
		return "DIRECT"
	case CHAT_GROUP:
		return "GROUP"
	case CHAT_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (c ChatType) ToDB() db.ChatTypeType {
	switch c {
	case CHAT_DIRECT:
		return db.ChatTypeTypeDIRECT
	case CHAT_GROUP:
		return db.ChatTypeTypeGROUP
	case CHAT_UNSPECIFIED:
		return db.ChatTypeTypeUNSPECIFIED
	default:
		return db.ChatTypeTypeUNSPECIFIED
	}
}

func ChatTypeFromString(s string) ChatType {
	switch s {
	case "DIRECT":
		return CHAT_DIRECT
	case "GROUP":
		return CHAT_GROUP
	case "UNSPECIFIED":
		return CHAT_UNSPECIFIED
	default:
		return CHAT_UNSPECIFIED
	}
}
