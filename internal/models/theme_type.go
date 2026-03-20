package models

import "github.com/juevigrace/diva-server/storage/db"

type Theme int

const (
	THEME_LIGHT Theme = iota
	THEME_DARK
	THEME_SYSTEM
)

func (t Theme) String() string {
	switch t {
	case THEME_LIGHT:
		return "LIGHT"
	case THEME_DARK:
		return "DARK"
	case THEME_SYSTEM:
		return "SYSTEM"
	default:
		return "SYSTEM"
	}
}

func (t Theme) ToDB() db.ThemeType {
	switch t {
	case THEME_LIGHT:
		return db.ThemeTypeLIGHT
	case THEME_DARK:
		return db.ThemeTypeDARK
	case THEME_SYSTEM:
		return db.ThemeTypeSYSTEM
	default:
		return db.ThemeTypeSYSTEM
	}
}

func ThemeFromString(s string) Theme {
	switch s {
	case "LIGHT":
		return THEME_LIGHT
	case "DARK":
		return THEME_DARK
	case "SYSTEM":
		return THEME_SYSTEM
	default:
		return THEME_SYSTEM
	}
}
