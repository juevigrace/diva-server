package models

import "github.com/juevigrace/diva-server/storage/db"

type ReactionType int

const (
	REACTION_LIKE ReactionType = iota
	REACTION_COMMENT
	REACTION_BOOKMARK
	REACTION_SHARE
	REACTION_UNSPECIFIED
)

func (r ReactionType) String() string {
	switch r {
	case REACTION_LIKE:
		return "LIKE"
	case REACTION_COMMENT:
		return "COMMENT"
	case REACTION_BOOKMARK:
		return "BOOKMARK"
	case REACTION_SHARE:
		return "SHARE"
	case REACTION_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (r ReactionType) ToDB() db.ReactionTypeType {
	switch r {
	case REACTION_LIKE:
		return db.ReactionTypeTypeLIKE
	case REACTION_COMMENT:
		return db.ReactionTypeTypeCOMMENT
	case REACTION_BOOKMARK:
		return db.ReactionTypeTypeBOOKMARK
	case REACTION_SHARE:
		return db.ReactionTypeTypeSHARE
	case REACTION_UNSPECIFIED:
		return db.ReactionTypeTypeUNSPECIFIED
	default:
		return db.ReactionTypeTypeUNSPECIFIED
	}
}

func ReactionTypeFromString(s string) ReactionType {
	switch s {
	case "LIKE":
		return REACTION_LIKE
	case "COMMENT":
		return REACTION_COMMENT
	case "BOOKMARK":
		return REACTION_BOOKMARK
	case "SHARE":
		return REACTION_SHARE
	case "UNSPECIFIED":
		return REACTION_UNSPECIFIED
	default:
		return REACTION_UNSPECIFIED
	}
}
