package ulid

import (
	oklog "github.com/oklog/ulid"
	"math/rand"
	"time"
)

// Erzeuge eine ULID
func GenerateULID() oklog.ULID {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	newID, _ := oklog.New(oklog.Timestamp(t), entropy)
	return newID
}

// Erzeugt einen ULID String
func GenerateStringULID() string {
	u := GenerateULID()
	return u.String()
}
