package storagepg

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"

	"github.com/google/uuid"
)

// NewNote adds new note to database
func (d *PostgresDB) NewNote(note *models.NewNote) (models.Note, error) {
	var newNote models.Note
	resErr := d.database.Get(&newNote, `insert into public.notes (user_id, title, note)
values ($1, $2, $3)
returning id, title, note, created, changed;`,
		note.UserID, note.Title, note.Note)
	if resErr != nil {
		return models.Note{}, resErr
	}
	return newNote, nil
}

// AllNotes returns all notes from database by current user
func (d *PostgresDB) AllNotes(userID uuid.UUID) ([]models.Note, error) {
	var notes []models.Note
	resErr := d.database.Select(&notes, `select id, title, note, created, changed
from public.notes where user_id = $1 order by changed desc nulls last, created desc`,
		userID)
	if resErr != nil {
		return nil, resErr
	}
	return notes, nil
}

// GetNote returns note from database by current user and note ID
func (d *PostgresDB) GetNote(noteID uuid.UUID, userID uuid.UUID) (models.Note, error) {
	var note models.Note
	resErr := d.database.Get(&note, `select id, title, note, created, changed
from public.notes where user_id = $1 and id = $2`,
		userID, noteID)
	if resErr != nil {
		return models.Note{}, resErr
	}
	return note, nil
}

// EditNote changes information in database about note by current user and note ID
func (d *PostgresDB) EditNote(note models.NewNote) (models.Note, error) {
	var newNote models.Note
	resErr := d.database.Get(&newNote, `update public.notes 
set title = $1,
    note = $2,
    changed = now()
where 1=1
and user_id = $3
and id = $4
returning id, title, note, created, changed;`,
		note.Title, note.Note, note.UserID, note.ID)
	if resErr != nil {
		return models.Note{}, resErr
	}
	return newNote, nil
}

// DeleteNote deletes note from database by current user and note ID
func (d *PostgresDB) DeleteNote(noteID uuid.UUID, userID uuid.UUID) error {
	res, resErr := d.database.Exec(`delete
from public.notes where user_id = $1 and id = $2`,
		userID, noteID)
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
