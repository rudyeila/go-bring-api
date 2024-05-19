package bring

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/rudyeila/go-bring-client/bring/model"
)

const (
	CONTENT_TYPE_X_WWW_FORM_URL_URL_ENCODED = "application/x-www-form-urlencoded"
	CONTEXT_TYPE_JSON                       = "application/json"
)

type Bring struct {
	client  *http.Client
	baseURL string
	creds   Creds
	conf    Config
	log     *slog.Logger
}

type Creds struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
}

func New(conf Config, logger *slog.Logger) *Bring {
	client := &http.Client{
		Timeout: conf.DefaultTimeout,
	}

	return &Bring{
		conf:    conf,
		baseURL: conf.BaseURL,
		client:  client,
		log:     logger,
	}
}

func (b *Bring) Login() error {
	authURL := fmt.Sprintf("%s/bringauth", b.conf.BaseURL)

	form := url.Values{}
	form.Add("email", b.conf.User)
	form.Add("password", b.conf.Password)

	req, err := http.NewRequest(http.MethodPost, authURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", CONTENT_TYPE_X_WWW_FORM_URL_URL_ENCODED)
	req.Header.Set("X-Bring-Api-Key", b.conf.ApiKey)
	req.Header.Set("X-Bring-Client", b.conf.ClientID)

	res, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, but got %s", res.Status)
	}

	authData := &model.LoginResponse{}
	err = json.NewDecoder(res.Body).Decode(authData)
	if err != nil {
		return err
	}

	b.creds = Creds{
		UserID:       authData.Uuid,
		AccessToken:  authData.AccessToken,
		RefreshToken: authData.RefreshToken,
		TokenType:    authData.TokenType,
		ExpiresIn:    authData.ExpiresIn,
	}

	return nil
}

func (b *Bring) GetLists() (*model.GetListsResponse, error) {
	userURL, err := b.GetUserBaseUrl()
	if err != nil {
		return nil, err
	}

	listsUrl := fmt.Sprintf("%s/lists", userURL)

	req, err := b.authenticatedRequest(http.MethodGet, listsUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	lists := &model.GetListsResponse{}
	err = json.NewDecoder(res.Body).Decode(lists)
	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (b *Bring) GetList(listID string) (*model.ListDetailResponse, error) {
	listURL := fmt.Sprintf("%s/bringlists/%s", b.baseURL, listID)

	req, err := b.authenticatedRequest(http.MethodGet, listURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	listRes := &model.ListDetailResponse{}
	err = json.NewDecoder(res.Body).Decode(listRes)
	if err != nil {
		return nil, err
	}

	return listRes, nil
}

func (b *Bring) AddItem(listID string, name, sub string) error {
	listURL := fmt.Sprintf("%s/bringlists/%s", b.baseURL, listID)

	form := url.Values{}
	form.Add("uuid", listID)
	form.Add("purchase", name)
	form.Add("specification", sub)

	req, err := b.authenticatedRequest(http.MethodPut, listURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", CONTENT_TYPE_X_WWW_FORM_URL_URL_ENCODED)

	res, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request not succesful. Response status was %s. Body: %v", res.Status, res.Body)
	}

	return nil
}

func (b *Bring) GetUserBaseUrl() (string, error) {
	if b.creds.UserID == "" {
		return "", errors.New("no user ID found in credentials")
	}

	return fmt.Sprintf("%s/bringusers/%s", b.baseURL, b.creds.UserID), nil
}

func (b *Bring) authenticatedRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	err = b.handleAuth()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", b.creds.TokenType, b.creds.AccessToken))
	req.Header.Add("X-Bring-Api-Key", b.conf.ApiKey)
	req.Header.Add("X-Bring-Client", b.conf.ClientID)

	return req, nil
}

func (b *Bring) handleAuth() error {
	if b.creds.AccessToken == "" {
		err := b.Login()
		if err != nil {
			return err
		}
	}

	// TODO check if token is expired and refresh it

	// TODO handle case where refresh token is also expired

	return nil
}
