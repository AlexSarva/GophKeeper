package storagepg

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"

	"github.com/google/uuid"
)

func (d *PostgresDB) NewFile(file *models.NewFile) (models.File, error) {
	var newFile models.File
	resErr := d.database.Get(&newFile, `insert into public.files (user_id, title, file, notes)
values ($1, $2, $3, $4)
returning id, title, file, notes, created, changed;`,
		file.UserID, file.Title, file.File, file.Notes)
	if resErr != nil {
		return models.File{}, resErr
	}
	return newFile, nil
}

func (d *PostgresDB) AllFiles(userID uuid.UUID) ([]models.File, error) {
	var files []models.File
	resErr := d.database.Select(&files, `select id, title, file, notes, created, changed
from public.files where user_id = $1 order by changed desc nulls last, created desc`,
		userID)
	if resErr != nil {
		return nil, resErr
	}
	return files, nil
}

func (d *PostgresDB) GetFile(cardID uuid.UUID, userID uuid.UUID) (models.File, error) {
	var file models.File
	resErr := d.database.Get(&file, `select id, title, file, notes, created, changed
from public.files where user_id = $1 and id = $2`,
		userID, cardID)
	if resErr != nil {
		return models.File{}, resErr
	}
	return file, nil
}

func (d *PostgresDB) EditFile(file models.NewFile) (models.File, error) {
	var newFile models.File
	resErr := d.database.Get(&newFile, `update public.cards 
set title = $1,
    file = $2,
    notes = $3,
    changed = now()
where 1=1
and user_id = $4
and id = $5
returning id, title, file, notes, created, changed;`,
		file.Title, file.File, file.Notes, file.UserID, file.ID)
	if resErr != nil {
		return models.File{}, resErr
	}
	return newFile, nil
}

func (d *PostgresDB) DeleteFile(fileID uuid.UUID, userID uuid.UUID) error {
	res, resErr := d.database.Exec(`delete
from public.files where user_id = $1 and id = $2`,
		userID, fileID)
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
