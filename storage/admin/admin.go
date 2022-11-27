package admin

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Admin initializing from PostgreSQL database
type Admin struct {
	database *sqlx.DB
}

// NewAdminDBConnection initializing from PostgreSQL database connection
func NewAdminDBConnection(config string) *Admin {

	db, err := sqlx.Connect("postgres", config)
	db.MustExec(ddl)
	if err != nil {
		log.Fatalln(err)
	}
	return &Admin{
		database: db,
	}
}

// Ping check availability of database
func (a *Admin) Ping() bool {
	//d.database.
	return a.database.Ping() == nil
}

// CheckUser insert new User in Database
func (a *Admin) CheckUser(userID uuid.UUID) bool {
	var user string
	resErr := a.database.Get(&user, "select email from public.users where id = $1", userID)
	if resErr != nil {
		return false
	}
	if len(user) == 0 {
		return false
	}
	return true
}

// Register insert new User in Databse
func (a *Admin) Register(user models.User) error {
	tx := a.database.MustBegin()
	resInsert, resErr := tx.NamedExec("INSERT INTO public.users (id, username, email, passwd, token, token_expires) VALUES (:id, :username, :email, :passwd, :token, :token_expires) on conflict (email) do nothing ", &user)
	if resErr != nil {
		return resErr
	}
	affectedRows, affectedRowsErr := resInsert.RowsAffected()
	if affectedRowsErr != nil {
		return affectedRowsErr
	}
	if affectedRows == 0 {
		return storage.ErrDuplicatePK
	}
	return tx.Commit()
}

// RenewToken refresh token for User in Databse
func (a *Admin) RenewToken(user models.User) error {
	tx := a.database.MustBegin()
	resInsert, resErr := tx.Exec(`
update public.users set token = $1,
                      token_expires = $2
                      where id = $3
`, user.Token, user.TokenExp, user.ID)
	if resErr != nil {
		return resErr
	}
	affectedRows, affectedRowsErr := resInsert.RowsAffected()
	if affectedRowsErr != nil {
		return affectedRowsErr
	}
	if affectedRows == 0 {
		return storage.ErrDuplicatePK
	}
	return tx.Commit()
}

// Login insert new User in Databse
func (a *Admin) Login(userLogin *models.UserLogin) (*models.User, error) {
	var user models.User
	err := a.database.Get(&user, "SELECT id, username, email, passwd, token, token_expires FROM public.users WHERE email=$1", userLogin.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserInfo get user credentials from database by username
func (a *Admin) GetUserInfo(userID uuid.UUID) (*models.User, error) {
	var userInfo models.User
	err := a.database.Get(&userInfo, "SELECT id, username, email, passwd, token, token_expires FROM public.users WHERE id=$1", userID)
	if err != nil {
		return nil, err
	}
	userInfo.Token = "Bearer " + userInfo.Token
	return &userInfo, nil
}
