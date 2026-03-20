package models

import "github.com/juevigrace/diva-server/storage/db"

type MessageType int

const (
	MESSAGE_TEXT MessageType = iota
	MESSAGE_MEDIA
	MESSAGE_SYSTEM
	MESSAGE_UNSPECIFIED
)

func (m MessageType) String() string {
	switch m {
	case MESSAGE_TEXT:
		return "TEXT"
	case MESSAGE_MEDIA:
		return "MEDIA"
	case MESSAGE_SYSTEM:
		return "SYSTEM"
	case MESSAGE_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (m MessageType) ToDB() db.MessageTypeType {
	switch m {
	case MESSAGE_TEXT:
		return db.MessageTypeTypeTEXT
	case MESSAGE_MEDIA:
		return db.MessageTypeTypeMEDIA
	case MESSAGE_SYSTEM:
		return db.MessageTypeTypeSYSTEM
	case MESSAGE_UNSPECIFIED:
		return db.MessageTypeTypeUNSPECIFIED
	default:
		return db.MessageTypeTypeUNSPECIFIED
	}
}

func MessageTypeFromString(s string) MessageType {
	switch s {
	case "TEXT":
		return MESSAGE_TEXT
	case "MEDIA":
		return MESSAGE_MEDIA
	case "SYSTEM":
		return MESSAGE_SYSTEM
	case "UNSPECIFIED":
		return MESSAGE_UNSPECIFIED
	default:
		return MESSAGE_UNSPECIFIED
	}
}
