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
	_, resRoleErr := tx.Exec("INSERT INTO public.acl (user_id, role_id, created_by) VALUES ($1, $2, $3)", user.ID, 4, user.ID)
	if resRoleErr != nil {
		return resRoleErr
	}
	return tx.Commit()
}

func (a *Admin) GetRoles() (map[string]int, error) {
	var roles []models.Roles
	err := a.database.Select(&roles, "SELECT id, role_name FROM public.roles")
	if err != nil {
		return nil, err
	}
	rolesMap := make(map[string]int)
	for i := 0; i < len(roles); i += 1 {
		rolesMap[roles[i].Name] = roles[i].ID
	}
	return rolesMap, nil
}

func (a *Admin) GetPermissions(userID uuid.UUID) (models.PermissionResponse, error) {
	var perms []string
	err := a.database.Select(&perms, `
select role_name from public.roles
where exists(select 1 from public.acl where 1=1
    and acl.role_id = roles.id
	and acl.is_del = 0
    and acl.user_id = $1)
order by id`, userID)
	if err != nil {
		return models.PermissionResponse{}, err
	}
	return models.PermissionResponse{
		ID:          userID,
		Permissions: perms,
	}, nil
}

// GrantPermission grants permissions
func (a *Admin) GrantPermission(permission models.Permission) error {
	tx := a.database.MustBegin()
	_, resErr := tx.NamedExec(`
insert into public.acl (user_id, role_id, created_by)
values (:user_id, :role_id, :created_by)
on conflict (user_id, role_id) do update
set is_del = 0,
    deleted_by = null,
    created_by = excluded.created_by;`, &permission)
	if resErr != nil {
		return resErr
	}
	return tx.Commit()
}

// RevokePermission grants permissions
func (a *Admin) RevokePermission(permission models.Permission) error {
	tx := a.database.MustBegin()
	_, resErr := tx.Exec(`
update public.acl set
is_del = 1,
deleted_by = $1
where user_id = $2 and role_id = $3;`, permission.CreatedBy, permission.ID, permission.RoleID)
	if resErr != nil {
		return resErr
	}
	return tx.Commit()
}

// Login insert new User in Databse
func (a *Admin) Login(email string) (models.PasswordCheck, error) {
	var passwd models.PasswordCheck
	err := a.database.Get(&passwd, "SELECT id, passwd FROM public.users WHERE email=$1", email)
	if err != nil {
		return models.PasswordCheck{}, err
	}
	return passwd, nil
}

// GetUserInfo get user credentials from database by username
func (a *Admin) GetUserInfo(userID uuid.UUID) (models.Token, error) {
	var userInfo models.Token
	err := a.database.Get(&userInfo, "SELECT id, username, email, is_admin, token, created FROM public.users WHERE id=$1", userID)
	if err != nil {
		return models.Token{}, err
	}
	userInfo.Token = "Bearer " + userInfo.Token
	return userInfo, nil
}

func (a *Admin) GetUserRoles(userID uuid.UUID) (map[string]bool, error) {
	var roles []string
	rolesMap := make(map[string]bool)
	err := a.database.Select(&roles, `
select r.role_name from public.acl
inner join public.roles r on r.id = acl.role_id
where user_id = $1
and is_del = 0
order by r.id;`, userID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, storage.ErrNoValues
	}
	for i := 0; i < len(roles); i += 1 {
		rolesMap[roles[i]] = true
	}
	return rolesMap, nil
}
