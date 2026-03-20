package models

import "github.com/juevigrace/diva-server/storage/db"

type Visibility int

const (
	VISIBILITY_PUBLIC Visibility = iota
	VISIBILITY_PRIVATE
	VISIBILITY_FRIENDS
	VISIBILITY_UNSPECIFIED
)

func (v Visibility) String() string {
	switch v {
	case VISIBILITY_PUBLIC:
		return "PUBLIC"
	case VISIBILITY_PRIVATE:
		return "PRIVATE"
	case VISIBILITY_FRIENDS:
		return "FRIENDS"
	case VISIBILITY_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (v Visibility) ToDB() db.VisibilityType {
	switch v {
	case VISIBILITY_PUBLIC:
		return db.VisibilityTypePUBLIC
	case VISIBILITY_PRIVATE:
		return db.VisibilityTypePRIVATE
	case VISIBILITY_FRIENDS:
		return db.VisibilityTypeFRIENDS
	case VISIBILITY_UNSPECIFIED:
		return db.VisibilityTypeUNSPECIFIED
	default:
		return db.VisibilityTypeUNSPECIFIED
	}
}

func VisibilityFromString(s string) Visibility {
	switch s {
	case "PUBLIC":
		return VISIBILITY_PUBLIC
	case "PRIVATE":
		return VISIBILITY_PRIVATE
	case "FRIENDS":
		return VISIBILITY_FRIENDS
	case "UNSPECIFIED":
		return VISIBILITY_UNSPECIFIED
	default:
		return VISIBILITY_UNSPECIFIED
	}
}
