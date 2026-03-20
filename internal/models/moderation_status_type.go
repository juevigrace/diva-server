package models

import "github.com/juevigrace/diva-server/storage/db"

type ModerationStatus int

const (
	MODERATION_PENDING ModerationStatus = iota
	MODERATION_APPROVED
	MODERATION_REJECTED
	MODERATION_HIDDEN
	MODERATION_UNSPECIFIED
)

func (m ModerationStatus) String() string {
	switch m {
	case MODERATION_PENDING:
		return "PENDING"
	case MODERATION_APPROVED:
		return "APPROVED"
	case MODERATION_REJECTED:
		return "REJECTED"
	case MODERATION_HIDDEN:
		return "HIDDEN"
	case MODERATION_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (m ModerationStatus) ToDB() db.ModerationStatusType {
	switch m {
	case MODERATION_PENDING:
		return db.ModerationStatusTypePENDING
	case MODERATION_APPROVED:
		return db.ModerationStatusTypeAPPROVED
	case MODERATION_REJECTED:
		return db.ModerationStatusTypeREJECTED
	case MODERATION_HIDDEN:
		return db.ModerationStatusTypeHIDDEN
	case MODERATION_UNSPECIFIED:
		return db.ModerationStatusTypeUNSPECIFIED
	default:
		return db.ModerationStatusTypeUNSPECIFIED
	}
}

func ModerationStatusFromString(s string) ModerationStatus {
	switch s {
	case "PENDING":
		return MODERATION_PENDING
	case "APPROVED":
		return MODERATION_APPROVED
	case "REJECTED":
		return MODERATION_REJECTED
	case "HIDDEN":
		return MODERATION_HIDDEN
	case "UNSPECIFIED":
		return MODERATION_UNSPECIFIED
	default:
		return MODERATION_UNSPECIFIED
	}
}
