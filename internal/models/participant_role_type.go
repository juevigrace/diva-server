package models

import "github.com/juevigrace/diva-server/storage/db"

type ParticipantRole int

const (
	PARTICIPANT_OWNER ParticipantRole = iota
	PARTICIPANT_ADMIN
	PARTICIPANT_MEMBER
	PARTICIPANT_UNSPECIFIED
)

func (p ParticipantRole) String() string {
	switch p {
	case PARTICIPANT_OWNER:
		return "OWNER"
	case PARTICIPANT_ADMIN:
		return "ADMIN"
	case PARTICIPANT_MEMBER:
		return "MEMBER"
	case PARTICIPANT_UNSPECIFIED:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

func (p ParticipantRole) ToDB() db.ParticipantRoleType {
	switch p {
	case PARTICIPANT_OWNER:
		return db.ParticipantRoleTypeOWNER
	case PARTICIPANT_ADMIN:
		return db.ParticipantRoleTypeADMIN
	case PARTICIPANT_MEMBER:
		return db.ParticipantRoleTypeMEMBER
	case PARTICIPANT_UNSPECIFIED:
		return db.ParticipantRoleTypeUNSPECIFIED
	default:
		return db.ParticipantRoleTypeUNSPECIFIED
	}
}

func ParticipantRoleFromString(s string) ParticipantRole {
	switch s {
	case "OWNER":
		return PARTICIPANT_OWNER
	case "ADMIN":
		return PARTICIPANT_ADMIN
	case "MEMBER":
		return PARTICIPANT_MEMBER
	case "UNSPECIFIED":
		return PARTICIPANT_UNSPECIFIED
	default:
		return PARTICIPANT_UNSPECIFIED
	}
}
