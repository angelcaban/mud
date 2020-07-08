package model

import (
	"github.com/gofrs/uuid"
)

type Registration struct {
	Id       uuid.UUID `stbl:"id, PRIMARY_KEY"`
	Name     string    `stbl:"name"`
	Email    string    `stbl:"email"`
	TimeZone string    `stbl:"timezone"`
	Password []byte    `stbl:"password"`
	ShortBio string    `stbl:"shortbio"`
}
