package workclient

import (
	"AlexSarva/GophKeeper/models"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
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
)

func switchTypesList(infoType string, r *gentleman.Response) (interface{}, error) {
	var res interface{}
	switch infoType {
	case "cards":
		var cards []models.Card
		if respErr := r.JSON(&cards); respErr != nil {
			return nil, respErr
		}
		res = cards
	case "notes":
		var notes []models.Note
		if respErr := r.JSON(&notes); respErr != nil {
			return nil, respErr
		}
		res = notes
	case "files":
		var files []models.File
		if respErr := r.JSON(&files); respErr != nil {
			return nil, respErr
		}
		res = files
	case "creds":
		var creds []models.Cred
		if respErr := r.JSON(&creds); respErr != nil {
			return nil, respErr
		}
		res = creds
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
	client  *gentleman.Client
	baseUrl string
}

// WorkClient initialize new client for work with service
func WorkClient(baseUrl string) (*Client, error) {
	cli := gentleman.New()
	cli.Use(timeout.Request(5 * time.Second))
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(5, 100*time.Millisecond), nil)))
	return &Client{
		client:  cli,
		baseUrl: baseUrl,
	}, nil
}

// UseToken method uses to add bearer token to client
func (c *Client) UseToken(bearer string) *Client {
	token := strings.Split(bearer, " ")
	log.Println(token[len(token)-1])
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

	result, resultErr := switchTypesList(infoType, res)
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
	if infoType != "files" {
		req.Use(body.JSON(elem))
	}
	if infoType == "files" {
		file := elem.(models.NewFile)
		req.Use(query.Set("title", file.Title))
		req.Use(query.Set("note", file.Notes))
		f := io.NopCloser(bytes.NewBuffer(file.File))
		req.Use(body.Reader(f))
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
	if infoType != "files" {
		req.Use(body.JSON(elem))
	}
	if infoType == "files" {
		file := elem.(models.NewFile)
		req.Use(query.Set("title", file.Title))
		req.Use(query.Set("note", file.Notes))
		f := io.NopCloser(bytes.NewBuffer(file.File))
		req.Use(body.Reader(f))
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
