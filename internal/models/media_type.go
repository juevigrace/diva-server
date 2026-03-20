package models

import "github.com/juevigrace/diva-server/storage/db"

type MediaType int

const (
	MEDIA_AUDIO MediaType = iota
	MEDIA_IMAGE
	MEDIA_VIDEO
	MEDIA_UNSPECIFIED
)

func (m MediaType) String() string {
	switch m {
	case MEDIA_AUDIO:
		return "AUDIO"
	case MEDIA_IMAGE:
		return "IMAGE"
	case MEDIA_VIDEO:
		return "VIDEO"
	case MEDIA_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (m MediaType) ToDB() db.MediaTypeType {
	switch m {
	case MEDIA_AUDIO:
		return db.MediaTypeTypeAUDIO
	case MEDIA_IMAGE:
		return db.MediaTypeTypeIMAGE
	case MEDIA_VIDEO:
		return db.MediaTypeTypeVIDEO
	case MEDIA_UNSPECIFIED:
		return db.MediaTypeTypeUNSPECIFIED
	default:
		return db.MediaTypeTypeUNSPECIFIED
	}
}

func MediaTypeFromString(s string) MediaType {
	switch s {
	case "AUDIO":
		return MEDIA_AUDIO
	case "IMAGE":
		return MEDIA_IMAGE
	case "VIDEO":
		return MEDIA_VIDEO
	case "UNSPECIFIED":
		return MEDIA_UNSPECIFIED
	default:
		return MEDIA_UNSPECIFIED
	}
}
