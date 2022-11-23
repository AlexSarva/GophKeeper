package storagepg

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"

	"github.com/google/uuid"
)

func (d *PostgresDB) NewCred(cred *models.NewCred) (models.Cred, error) {
	var newCred models.Cred
	resErr := d.database.Get(&newCred, `insert into public.creds (user_id, title, login, passwd, notes)
values ($1, $2, $3, $4, $5)
returning id, title, login, passwd, notes, created, changed;`,
		cred.UserID, cred.Title, cred.Login, cred.Passwd, cred.Notes)
	if resErr != nil {
		return models.Cred{}, resErr
	}
	return newCred, nil
}

func (d *PostgresDB) AllCreds(userID uuid.UUID) ([]models.Cred, error) {
	var creds []models.Cred
	resErr := d.database.Select(&creds, `select id, title, login, passwd, notes, created, changed
from public.creds where user_id = $1 order by changed desc nulls last, created desc`,
		userID)
	if resErr != nil {
		return nil, resErr
	}
	return creds, nil
}

func (d *PostgresDB) GetCred(credID, userID uuid.UUID) (models.Cred, error) {
	var cred models.Cred
	resErr := d.database.Get(&cred, `select id, title, login, passwd, notes, created, changed
from public.creds where user_id = $1 and id = $2`,
		userID, credID)
	if resErr != nil {
		return models.Cred{}, resErr
	}
	return cred, nil
}

func (d *PostgresDB) EditCred(cred models.NewCred) (models.Cred, error) {
	var newCred models.Cred
	resErr := d.database.Get(&newCred, `update public.creds
set title = $1,
    login = $2,
    passwd = $3,
    notes = $4,
    changed = now()
where 1=1
and user_id = $5
and id = $6
returning id, title, login, passwd, notes, created, changed;`,
		cred.Title, cred.Login, cred.Passwd, cred.Notes, cred.UserID, cred.ID)
	if resErr != nil {
		return models.Cred{}, resErr
	}
	return newCred, nil
}

func (d *PostgresDB) DeleteCred(credID uuid.UUID, userID uuid.UUID) error {
	res, resErr := d.database.Exec(`delete
from public.creds where user_id = $1 and id = $2`,
		userID, credID)
	if resErr != nil {
		return resErr
	}
	affectedRows, affectedRowsErr := res.RowsAffected()
	if affectedRowsErr != nil {
		return affectedRowsErr
	}
	if affectedRows == 0 {
		return storage.ErrNoValues
	}
	return nil
}
