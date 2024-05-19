package model

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Uuid          string `json:"uuid"`
	PublicUuid    string `json:"publicUuid"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	PhotoPath     string `json:"photoPath"`
	BringListUUID string `json:"bringListUUID"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
}

type GetListsResponse struct {
	Lists []List `json:"lists"`
}

type List struct {
	ListUuid string `json:"listUuid"`
	Name     string `json:"name"`
	Theme    string `json:"theme"`
}

type ListItem struct {
	Specification string `json:"specification"`
	Name          string `json:"name"`
}

type AddItemBody struct {
	Uuid          string `json:"uuid"`
	Purchase      string `json:"purchase"`
	Specification string `json:"specification"`
}

type ListDetailResponse struct {
	Uuid     string     `json:"uuid"`
	Status   string     `json:"status"`
	Purchase []ListItem `json:"purchase"`
	Recently []ListItem `json:"recently"`
}
