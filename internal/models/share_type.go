package models

import "github.com/juevigrace/diva-server/storage/db"

type ShareType int

const (
	SHARE_DIRECT ShareType = iota
	SHARE_EXTERNAL
	SHARE_EMBED
	SHARE_DOWNLOAD
	SHARE_UNSPECIFIED
)

func (s ShareType) String() string {
	switch s {
	case SHARE_DIRECT:
		return "DIRECT"
	case SHARE_EXTERNAL:
		return "EXTERNAL"
	case SHARE_EMBED:
		return "EMBED"
	case SHARE_DOWNLOAD:
		return "DOWNLOAD"
	case SHARE_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (s ShareType) ToDB() db.ShareTypeType {
	switch s {
	case SHARE_DIRECT:
		return db.ShareTypeTypeDIRECT
	case SHARE_EXTERNAL:
		return db.ShareTypeTypeEXTERNAL
	case SHARE_EMBED:
		return db.ShareTypeTypeEMBED
	case SHARE_DOWNLOAD:
		return db.ShareTypeTypeDOWNLOAD
	case SHARE_UNSPECIFIED:
		return db.ShareTypeTypeUNSPECIFIED
	default:
		return db.ShareTypeTypeUNSPECIFIED
	}
}

func ShareTypeFromString(s string) ShareType {
	switch s {
	case "DIRECT":
		return SHARE_DIRECT
	case "EXTERNAL":
		return SHARE_EXTERNAL
	case "EMBED":
		return SHARE_EMBED
	case "DOWNLOAD":
		return SHARE_DOWNLOAD
	case "UNSPECIFIED":
		return SHARE_UNSPECIFIED
	default:
		return SHARE_UNSPECIFIED
	}
}
