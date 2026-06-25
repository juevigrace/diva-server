package storage

import "github.com/google/uuid"

type RoleType string

const (
	RoleTypeUSER      RoleType = "USER"
	RoleTypeMODERATOR RoleType = "MODERATOR"
	RoleTypeADMIN     RoleType = "ADMIN"
)

type UserStatusType string

const (
	UserStatusTypeACTIVE    UserStatusType = "ACTIVE"
	UserStatusTypeSUSPENDED UserStatusType = "SUSPENDED"
	UserStatusTypeINACTIVE  UserStatusType = "INACTIVE"
)

type SessionStatusType string

const (
	SessionStatusTypeACTIVE  SessionStatusType = "ACTIVE"
	SessionStatusTypeEXPIRED SessionStatusType = "EXPIRED"
	SessionStatusTypeCLOSED  SessionStatusType = "CLOSED"
)

type SessionType string

const (
	SessionTypeNORMAL   SessionType = "NORMAL"
	SessionTypeTEMPORAL SessionType = "TEMPORAL"
)

type ThemeType string

const (
	ThemeTypeLIGHT  ThemeType = "LIGHT"
	ThemeTypeDARK   ThemeType = "DARK"
	ThemeTypeSYSTEM ThemeType = "SYSTEM"
)

type MediaType string

const (
	MediaTypeTypeAUDIO       MediaType = "AUDIO"
	MediaTypeTypeIMAGE       MediaType = "IMAGE"
	MediaTypeTypeVIDEO       MediaType = "VIDEO"
	MediaTypeTypeUNSPECIFIED MediaType = "UNSPECIFIED"
)

type DivaUser struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PhoneNumber  string
	PasswordHash string
	Role         RoleType
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    *int64
}

type DivaPermission struct {
	ID          uuid.UUID
	Name        string
	Description string
	Action      string
	RoleLevel   RoleType
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   *int64
}

type DivaSession struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	AccessToken     string
	RefreshToken    string
	Device          string
	Type            SessionType
	Status          SessionStatusType
	IpAddress       string
	UserAgent       string
	AccessExpiresAt int64
	RefreshExpiresAt int64
	CreatedAt       int64
	UpdatedAt       int64
}

type DivaUserState struct {
	UserID       uuid.UUID
	Verified     bool
	Status       UserStatusType
	LastActiveAt int64
	UpdatedAt    int64
}

type DivaUserProfile struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	BirthDate *int64
	Alias     string
	Bio       string
	Avatar    string
	UpdatedAt int64
}

type DivaUserPreference struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Device              string
	Theme               ThemeType
	OnboardingCompleted bool
	Language            string
	LastSyncAt          int64
	CreatedAt           int64
	UpdatedAt           int64
}

type DivaUserPermission struct {
	PermissionID uuid.UUID
	UserID       uuid.UUID
	GrantedBy    *uuid.UUID
	Granted      bool
	GrantedAt    int64
	ExpiresAt    *int64
	UpdatedAt    int64
}

type DivaAction struct {
	ID     uuid.UUID
	Name   string
	UserID uuid.UUID
}

type DivaActionVerification struct {
	ActionID  uuid.UUID
	Token     string
	Verified  bool
	ExpiresAt int64
	UsedAt    *int64
}
