package models

import "github.com/juevigrace/diva-server/storage/db"

type CollectionType int

const (
	COLLECTION_ALBUM CollectionType = iota
	COLLECTION_PLAYLIST
	COLLECTION_MIX
	COLLECTION_FAVORITES
	COLLECTION_FEATURED
	COLLECTION_TRENDING
	COLLECTION_UNSPECIFIED
)

func (c CollectionType) String() string {
	switch c {
	case COLLECTION_ALBUM:
		return "ALBUM"
	case COLLECTION_PLAYLIST:
		return "PLAYLIST"
	case COLLECTION_MIX:
		return "MIX"
	case COLLECTION_FAVORITES:
		return "FAVORITES"
	case COLLECTION_FEATURED:
		return "FEATURED"
	case COLLECTION_TRENDING:
		return "TRENDING"
	case COLLECTION_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (c CollectionType) ToDB() db.CollectionTypeType {
	switch c {
	case COLLECTION_ALBUM:
		return db.CollectionTypeTypeALBUM
	case COLLECTION_PLAYLIST:
		return db.CollectionTypeTypePLAYLIST
	case COLLECTION_MIX:
		return db.CollectionTypeTypeMIX
	case COLLECTION_FAVORITES:
		return db.CollectionTypeTypeFAVORITES
	case COLLECTION_FEATURED:
		return db.CollectionTypeTypeFEATURED
	case COLLECTION_TRENDING:
		return db.CollectionTypeTypeTRENDING
	case COLLECTION_UNSPECIFIED:
		return db.CollectionTypeTypeUNSPECIFIED
	default:
		return db.CollectionTypeTypeUNSPECIFIED
	}
}

func CollectionTypeFromString(s string) CollectionType {
	switch s {
	case "ALBUM":
		return COLLECTION_ALBUM
	case "PLAYLIST":
		return COLLECTION_PLAYLIST
	case "MIX":
		return COLLECTION_MIX
	case "FAVORITES":
		return COLLECTION_FAVORITES
	case "FEATURED":
		return COLLECTION_FEATURED
	case "TRENDING":
		return COLLECTION_TRENDING
	case "UNSPECIFIED":
		return COLLECTION_UNSPECIFIED
	default:
		return COLLECTION_UNSPECIFIED
	}
}
