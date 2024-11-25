package game

import "errors"

type Realm string

const (
	RealmHuman  Realm = "human"
	RealmCyborg Realm = "cyborg"
)

var (
	ErrorInvalidRealm = errors.New("invalid realm")
)

func IsValidRealm(realm Realm) bool {
	return realm == RealmHuman ||
		realm == RealmCyborg
}
