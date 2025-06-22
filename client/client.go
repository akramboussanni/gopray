package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type MawaqitClient struct {
	AuthToken  string
	BaseURL    string
	HTTPClient *http.Client
}

const baseURL = "https://mawaqit.net/api"

func newClient() *MawaqitClient {
	return &MawaqitClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func NewWithToken(token string) *MawaqitClient {
	cli := newClient()
	cli.AuthToken = token
	return cli
}

func NewWithCredentials(username, password string) (*MawaqitClient, error) {
	cli := newClient()
	data := username + ":" + password
	auth64 := base64.StdEncoding.EncodeToString([]byte(data))

	req, err := http.NewRequest("POST", cli.BaseURL+"/2.0/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+auth64)
	response, err := cli.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", "auth failed: "+response.Status)
	}

	var result struct {
		ApiAccessToken string `json:"apiAccessToken"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	cli.AuthToken = result.ApiAccessToken
	return cli, nil
}

/*
   async def login(self) -> None:
       """Log into the MAWAQIT website."""

       if (self.username is None) or (self.password is None):
           raise MissingCredentials("Please provide a MAWAQIT login and password.")

       auth = aiohttp.BasicAuth(self.username, self.password)

       endpoint_url = LOGIN_URL

       async with await self.session.post(endpoint_url, auth=auth) as response:
           if response.status == 401:
               raise BadCredentialsException(
                   "Authentication failed. Please check your MAWAQIT credentials."
               )
           elif response.status != 200:
               raise NotAuthenticatedException("Authentication failed. Please retry.")

           data = await response.text()

           self.token = json.loads(data)["apiAccessToken"]
*/
