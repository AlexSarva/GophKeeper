package storagepg

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"

	"github.com/google/uuid"
)

// NewCard adds new credit card to database
func (d *PostgresDB) NewCard(card *models.NewCard) (models.Card, error) {
	var newCard models.Card
	resErr := d.database.Get(&newCard, `insert into public.cards (user_id, title, card_number,
card_owner, card_exp, notes)
values ($1, $2, $3, $4, $5, $6)
returning id, title, card_number,
card_owner, card_exp, notes, created, changed;`,
		card.UserID, card.Title, card.CardNumber, card.CardOwner, card.CardExp, card.Notes)
	if resErr != nil {
		return models.Card{}, resErr
	}
	return newCard, nil
}

// AllCards returns all credit cards from database by current user
func (d *PostgresDB) AllCards(userID uuid.UUID) ([]models.Card, error) {
	var cards []models.Card
	resErr := d.database.Select(&cards, `select id, title, card_number,
card_owner, card_exp, notes, created, changed
from public.cards where user_id = $1 order by changed desc nulls last, created desc`,
		userID)
	if resErr != nil {
		return nil, resErr
	}
	return cards, nil
}

// GetCard returns credit card from database by current user and credit card ID
func (d *PostgresDB) GetCard(cardID uuid.UUID, userID uuid.UUID) (models.Card, error) {
	var card models.Card
	resErr := d.database.Get(&card, `select id, title, card_number,
card_owner, card_exp, notes, created, changed
from public.cards where user_id = $1 and id = $2`,
		userID, cardID)
	if resErr != nil {
		return models.Card{}, resErr
	}
	return card, nil
}

// EditCard changes information in database about credit card by current user and credit card ID
func (d *PostgresDB) EditCard(card models.NewCard) (models.Card, error) {
	var newCard models.Card
	resErr := d.database.Get(&newCard, `update public.cards 
set title = $1,
    card_number = $2,
    card_owner = $3,
    card_exp = $4,
    notes = $5,
    changed = now()
where 1=1
and user_id = $6
and id = $7
returning id, title, card_number,
card_owner, card_exp, notes, created, changed;`,
		card.Title, card.CardNumber, card.CardOwner, card.CardExp, card.Notes, card.UserID, card.ID)
	if resErr != nil {
		return models.Card{}, resErr
	}
	return newCard, nil
}

// DeleteCard deletes credit card from database by current user and credit card ID
func (d *PostgresDB) DeleteCard(cardID uuid.UUID, userID uuid.UUID) error {
	res, resErr := d.database.Exec(`delete
from public.cards where user_id = $1 and id = $2`,
		userID, cardID)
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
