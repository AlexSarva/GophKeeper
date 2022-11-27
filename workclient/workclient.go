package workclient

import (
	"AlexSarva/GophKeeper/crypto"
	"AlexSarva/GophKeeper/crypto/cryptorsa"
	"AlexSarva/GophKeeper/crypto/symmetric"
	"AlexSarva/GophKeeper/models"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/eapache/go-resiliency.v1/retrier"
	retry "gopkg.in/h2non/gentleman-retry.v2"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/auth"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/query"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

var (
	ErrUserExist      = errors.New("login is already taken")
	ErrCreds          = errors.New("invalid login/password pair")
	ErrToken          = errors.New("invalid token")
	ErrInternalServer = errors.New("an internal server error")
	ErrReqFormat      = errors.New("invalid request format")
	ErrNoData         = errors.New("no info in DB")
	ErrTokenExpired   = errors.New("unauthorized: token is expired")
)

func InitCryptorizer(ketsPath string, size int) *crypto.Cryptorizer {
	cryptorizer := cryptorsa.InitRSAConf(ketsPath, size)
	return &crypto.Cryptorizer{
		Cryptorizer: cryptorizer,
	}
}

func (c *Client) switchTypesList(infoType string, r *gentleman.Response) (interface{}, error) {
	var res interface{}
	switch infoType {
	case "cards":
		var cards []models.Card
		if respErr := r.JSON(&cards); respErr != nil {
			return nil, respErr
		}
		var decrCards []models.Card
		for _, card := range cards {
			if decryptErr := card.Decrypt(c.cryptorizer); decryptErr != nil {
				return nil, decryptErr
			}
			decrCards = append(decrCards, card)
		}
		res = decrCards
		break
	case "notes":
		var notes []models.Note
		if respErr := r.JSON(&notes); respErr != nil {
			return nil, respErr
		}
		var descrNotes []models.Note
		for _, note := range notes {
			if decryptErr := note.Decrypt(c.cryptorizer); decryptErr != nil {
				return nil, decryptErr
			}
			descrNotes = append(descrNotes, note)
		}
		res = descrNotes
		break
	case "files":
		var files []models.File
		if respErr := r.JSON(&files); respErr != nil {
			return nil, respErr
		}
		var descrFiles []models.File
		for _, file := range files {
			if symDecrErr := file.SymDecrypt(c.symmCrypto); symDecrErr != nil {
				return nil, symDecrErr
			}
			descrFiles = append(descrFiles, file)
		}
		res = descrFiles
		break
	case "creds":
		var creds []models.Cred
		if respErr := r.JSON(&creds); respErr != nil {
			return nil, respErr
		}
		var descrCreds []models.Cred
		for _, cred := range creds {
			if decryptErr := cred.Decrypt(c.cryptorizer); decryptErr != nil {
				return nil, decryptErr
			}
			descrCreds = append(descrCreds, cred)
		}
		res = descrCreds
		break
	}
	return res, nil
}

func switchType(infoType string, r *gentleman.Response) (interface{}, error) {
	var res interface{}
	switch infoType {
	case "cards":
		var card models.Card
		if respErr := r.JSON(&card); respErr != nil {
			return nil, respErr
		}
		res = card
	case "notes":
		var note models.Note
		if respErr := r.JSON(&note); respErr != nil {
			return nil, respErr
		}
		res = note
	case "files":
		var file models.File
		if respErr := r.JSON(&file); respErr != nil {
			return nil, respErr
		}
		res = file
	case "creds":
		var cred models.Cred
		if respErr := r.JSON(&cred); respErr != nil {
			return nil, respErr
		}
		res = cred
	}
	return res, nil
}

// Client custom type of work client
type Client struct {
	client      *gentleman.Client
	baseUrl     string
	cryptorizer *crypto.Cryptorizer
	symmCrypto  *symmetric.SymmetricCrypto
}

// InitClient initialize new client for work with service
func InitClient(cfg *models.GUIConfig) (*Client, error) {
	cli := gentleman.New()
	cli.Use(timeout.Request(5 * time.Second))
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(5, 100*time.Millisecond), nil)))
	cryptorizer := InitCryptorizer(cfg.KeysPath, cfg.KeysSize)
	symmCrypto := symmetric.SymmCrypto(cfg.Secret)
	return &Client{
		client:      cli,
		baseUrl:     cfg.ServerAddress,
		cryptorizer: cryptorizer,
		symmCrypto:  symmCrypto,
	}, nil
}

// UseToken method uses to add bearer token to client
func (c *Client) UseToken(bearer string) *Client {
	token := strings.Split(bearer, " ")
	c.client.Use(auth.Bearer(token[len(token)-1]))
	return c
}

// Register provides registration in service
func (c *Client) Register(userInfo *models.UserRegister) (*models.User, error) {
	var user *models.User
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/register", c.baseUrl))
	req.Method("POST")
	req.SetHeader("Content-Type", "application/json")
	req.Use(body.JSON(userInfo))
	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if res.StatusCode == 409 {
			return nil, ErrUserExist
		}
		if res.StatusCode == 417 {
			return nil, errors.New(res.String())
		}
		if res.StatusCode == 500 {
			return nil, ErrInternalServer
		}
		return nil, ErrReqFormat
	}
	if respErr := res.JSON(&user); respErr != nil {
		return nil, respErr
	}

	if res.Ok {
		c.UseToken(user.Token)
	}

	return user, nil
}

// Login provides log in service
func (c *Client) Login(userInfo *models.UserLogin) (*models.User, error) {
	var user *models.User
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/login", c.baseUrl))
	req.Method("POST")
	req.SetHeader("Content-Type", "application/json")
	req.Use(body.JSON(userInfo))
	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if res.StatusCode == 401 {
			return nil, ErrCreds
		}
		if res.StatusCode == 403 {
			return nil, ErrTokenExpired
		}
		if res.StatusCode == 500 {
			return nil, ErrInternalServer
		}
		return nil, ErrReqFormat
	}
	if respErr := res.JSON(&user); respErr != nil {
		return nil, respErr
	}

	if res.Ok {
		c.UseToken(user.Token)
	}

	return user, nil
}

// Me provides get information about user
func (c *Client) Me() (*models.User, error) {
	var user *models.User
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/users/me", c.baseUrl))
	req.Method("GET")
	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if res.StatusCode == 401 {
			return nil, ErrToken
		}
		if res.StatusCode == 500 {
			return nil, ErrInternalServer
		}
		return nil, ErrReqFormat
	}
	if respErr := res.JSON(&user); respErr != nil {
		return nil, respErr
	}

	return user, nil
}

// ElementList returns list of elements of selected type
func (c *Client) ElementList(infoType string) (interface{}, error) {
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/info/%s", c.baseUrl, infoType))
	req.Method("GET")
	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if res.StatusCode == 401 {
			return nil, ErrToken
		}
		if res.StatusCode == 500 {
			return nil, ErrInternalServer
		}
		return nil, ErrReqFormat
	}

	result, resultErr := c.switchTypesList(infoType, res)
	if resultErr != nil {
		return nil, resultErr
	}

	return result, nil
}

// Element returns element of selected type and id
func (c *Client) Element(infoType string, id uuid.UUID) (interface{}, error) {
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/info/%s/%s", c.baseUrl, infoType, id))
	req.Method("GET")
	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if res.StatusCode == 401 {
			return nil, ErrToken
		}
		if res.StatusCode == 409 {
			return nil, ErrNoData
		}
		if res.StatusCode == 500 {
			return nil, ErrInternalServer
		}
		return nil, ErrReqFormat
	}

	result, resultErr := switchType(infoType, res)
	if resultErr != nil {
		return nil, resultErr
	}

	return result, nil
}

// AddElement provides add element of selected type in service
func (c *Client) AddElement(infoType string, elem interface{}) (interface{}, error) {
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/info/%s", c.baseUrl, infoType))
	req.Method("POST")

	switch infoType {
	case "cards":
		card := elem.(*models.NewCard)
		if checkErr := card.CheckValid(); checkErr != nil {
			return nil, checkErr
		}
		if cryptoErr := card.Encrypt(c.cryptorizer); cryptoErr != nil {
			return nil, cryptoErr
		}
		req.Use(body.JSON(card))
		break
	case "creds":
		cred := elem.(*models.NewCred)
		if cryptoErr := cred.Encrypt(c.cryptorizer); cryptoErr != nil {
			return nil, cryptoErr
		}
		req.Use(body.JSON(cred))
		break
	case "notes":
		note := elem.(*models.NewNote)
		if cryptoErr := note.Encrypt(c.cryptorizer); cryptoErr != nil {
			return nil, cryptoErr
		}
		req.Use(body.JSON(note))
		break
	case "files":
		file := elem.(*models.NewFile)
		file.SymEncrypt(c.symmCrypto)
		req.Use(query.Set("title", file.Title))
		req.Use(query.Set("notes", file.Notes))
		req.Use(query.Set("filename", file.FileName))
		f := io.NopCloser(bytes.NewBuffer(file.File))
		req.Use(body.Reader(f))
		break
	default:
		return nil, errors.New("wrong info type parameter")
	}

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if res.StatusCode == 400 {
			return nil, ErrReqFormat
		}
		if res.StatusCode == 401 {
			return nil, ErrToken
		}
		if res.StatusCode == 500 {
			return nil, ErrInternalServer
		}
		return nil, ErrReqFormat
	}

	result, resultErr := switchType(infoType, res)
	if resultErr != nil {
		return nil, resultErr
	}

	return result, nil
}

// EditElement provides edit element of selected type in service
func (c *Client) EditElement(infoType string, elem interface{}, id uuid.UUID) (interface{}, error) {
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/info/%s/%s", c.baseUrl, infoType, id))
	req.Method("PATCH")
	switch infoType {
	case "cards":
		card := elem.(*models.NewCard)
		if checkErr := card.CheckValid(); checkErr != nil {
			return nil, checkErr
		}
		if cryptoErr := card.Encrypt(c.cryptorizer); cryptoErr != nil {
			return nil, cryptoErr
		}
		req.Use(body.JSON(card))
		break
	case "creds":
		cred := elem.(*models.NewCred)
		if cryptoErr := cred.Encrypt(c.cryptorizer); cryptoErr != nil {
			return nil, cryptoErr
		}
		req.Use(body.JSON(cred))
		break
	case "notes":
		note := elem.(*models.NewNote)
		if cryptoErr := note.Encrypt(c.cryptorizer); cryptoErr != nil {
			return nil, cryptoErr
		}
		req.Use(body.JSON(note))
		break
	case "files":
		file := elem.(*models.NewFile)
		file.SymEncrypt(c.symmCrypto)
		req.Use(query.Set("title", file.Title))
		req.Use(query.Set("filename", file.FileName))
		req.Use(query.Set("notes", file.Notes))
		f := io.NopCloser(bytes.NewBuffer(file.File))
		req.Use(body.Reader(f))
		break
	default:
		return nil, errors.New("wrong info type parameter")
	}

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if res.StatusCode == 400 {
			return nil, ErrReqFormat
		}
		if res.StatusCode == 401 {
			return nil, ErrToken
		}
		if res.StatusCode == 409 {
			return nil, ErrNoData
		}
		if res.StatusCode == 500 {
			return nil, ErrInternalServer
		}
		return nil, ErrReqFormat
	}

	result, resultErr := switchType(infoType, res)
	if resultErr != nil {
		return nil, resultErr
	}

	return result, nil
}

// Delete removes element from service by selected type and id
func (c *Client) Delete(infoType string, id uuid.UUID) (bool, error) {
	req := c.client.Request()
	req.URL(fmt.Sprintf("%s/info/%s/%s", c.baseUrl, infoType, id))
	req.Method("DELETE")
	res, err := req.Send()
	if err != nil {
		return false, err
	}
	if !res.Ok {
		if res.StatusCode == 401 {
			return false, ErrToken
		}
		if res.StatusCode == 409 {
			return false, ErrNoData
		}
		if res.StatusCode == 500 {
			return false, ErrInternalServer
		}
		return false, ErrReqFormat
	}

	if res.StatusCode != 200 {
		return false, errors.New("something went wrong")
	}

	return true, nil
}
