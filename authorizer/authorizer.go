package authorizer

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage/admin"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidAccessToken = errors.New("invalid auth token")
var ErrGenerateToken = errors.New("something wrong happens when generate token")
var ErrHashPassword = errors.New("something wrong happens when hashing password")
var ErrComparePassword = errors.New("password doesnt match")
var ErrNoUserExists = errors.New("user doesnt exist in database")

type Claims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}

type Authorizer struct {
	adminDB        *admin.Admin
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthorizer(db *admin.Admin, signingKey []byte, expireDuration time.Duration) *Authorizer {
	return &Authorizer{
		adminDB:        db,
		signingKey:     signingKey,
		expireDuration: expireDuration, // 86400
	}
}

func (a *Authorizer) SignUp(user models.User) (*models.User, error) {
	// Create password hash

	hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if bcrypteErr != nil {
		return nil, ErrHashPassword
	}

	user.Password = string(hashedPassword)

	expires := time.Now().Add(a.expireDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expires),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserID: user.ID,
	})

	tokenValue, tokenErr := token.SignedString(a.signingKey)
	if tokenErr != nil {
		return nil, ErrGenerateToken
	}

	user.Token = tokenValue
	user.TokenExp = expires

	err := a.adminDB.Register(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.Token = fmt.Sprintf("Bearer %s", user.Token)

	return &user, nil
}

func (a *Authorizer) SignIn(userLogin *models.UserLogin) (*models.User, error) {
	userCred, err := a.adminDB.Login(userLogin)
	if err != nil {
		if errors.As(sql.ErrNoRows, &err) {
			return nil, ErrNoUserExists
		}
		return nil, err
	}

	cryptErr := bcrypt.CompareHashAndPassword([]byte(userCred.Password), []byte(userLogin.Password))
	if cryptErr != nil {
		return nil, ErrComparePassword
	}

	userCred.Password = ""
	userCred.Token = fmt.Sprintf("Bearer %s", userCred.Token)

	return userCred, nil
}

func (a *Authorizer) RenewToken(user models.User) (*models.User, error) {

	expires := time.Now().Add(a.expireDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expires),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserID: user.ID,
	})

	tokenValue, tokenErr := token.SignedString(a.signingKey)
	if tokenErr != nil {
		return nil, ErrGenerateToken
	}

	user.Token = tokenValue
	user.TokenExp = expires

	err := a.adminDB.RenewToken(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.Token = fmt.Sprintf("Bearer %s", user.Token)

	return &user, nil
}

func (a *Authorizer) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return uuid.UUID{}, ErrInvalidAccessToken
}
