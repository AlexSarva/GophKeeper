package storagepg

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"log"

	"github.com/google/uuid"
)

// NewFile adds new file to database
func (d *PostgresDB) NewFile(file *models.NewFile) (models.File, error) {
	var newFile models.File
	resErr := d.database.Get(&newFile, `insert into public.files (user_id, title, file_name, file, notes)
values ($1, $2, $3, $4, $5)
returning id, title, file, file_name, notes, created, changed;`,
		file.UserID, file.Title, file.FileName, file.File, file.Notes)
	if resErr != nil {
		return models.File{}, resErr
	}
	return newFile, nil
}

// AllFiles returns all files from database by current user
func (d *PostgresDB) AllFiles(userID uuid.UUID) ([]models.File, error) {
	var files []models.File
	resErr := d.database.Select(&files, `select id, title, file_name, file, notes, created, changed
from public.files where user_id = $1 order by changed desc nulls last, created desc`,
		userID)
	if resErr != nil {
		return nil, resErr
	}
	return files, nil
}

// GetFile returns file from database by current user and file ID
func (d *PostgresDB) GetFile(cardID uuid.UUID, userID uuid.UUID) (models.File, error) {
	var file models.File
	resErr := d.database.Get(&file, `select id, title, file_name, file, notes, created, changed
from public.files where user_id = $1 and id = $2`,
		userID, cardID)
	if resErr != nil {
		return models.File{}, resErr
	}
	return file, nil
}

// EditFile changes information in database about file by current user and file ID
func (d *PostgresDB) EditFile(file *models.NewFile) (models.File, error) {
	var newFile models.File
	log.Printf("%+v\n", file)
	resErr := d.database.Get(&newFile, `update public.files 
set title = $1,
    file = $2,
    file_name = $3,
    notes = $4,
    changed = now()
where 1=1
and user_id = $5
and id = $6
returning id, title, file_name, file, notes, created, changed;`,
		file.Title, file.File, file.FileName, file.Notes, file.UserID, file.ID)
	if resErr != nil {
		return models.File{}, resErr
	}
	return newFile, nil
}

// DeleteFile deletes file from database by current user and file ID
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
