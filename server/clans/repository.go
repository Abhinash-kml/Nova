package clans

import "github.com/google/uuid"

type ClansRepository interface {
	Get(uuid.UUID) Clan
	GetByName(string) Clan
	GetAll() []Clan
	Add(Clan) bool
	Delete(uuid.UUID) bool
}
