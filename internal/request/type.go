package request

import (
	"database/sql"
)

const (
	TypeCharacterApplication string = "CharacterApplication"
	IsTypeQuery              string = "SELECT IF(requests.type = ?, true, false) FROM requests WHERE id = ?;"
)

func IsTypeDB(db *sql.DB, t string, rid int64) (bool, error) {
	i := 0
	r := db.QueryRow(IsTypeQuery, t, rid)
	err := r.Scan(&i)
	if err != nil {
		return false, err
	}

	return i == 1, nil
}

func IsTypeTx(tx *sql.Tx, t string, rid int64) (bool, error) {
	i := 0
	r := tx.QueryRow(IsTypeQuery, t, rid)
	err := r.Scan(&i)
	if err != nil {
		return false, err
	}

	return i == 1, nil
}

func IsCharacterApplicationDB(db *sql.DB, rid int64) (bool, error) {
	return IsTypeDB(db, TypeCharacterApplication, rid)
}

func IsCharacterApplicationTx(tx *sql.Tx, rid int64) (bool, error) {
	return IsTypeTx(tx, TypeCharacterApplication, rid)
}
