package client

import "net/http"

type KickflowClient struct {
	BaseURL     string
	AccessToken string
	CallerID    string
	HTTPClient  *http.Client
}

func NewKickflowClient(baseURL string, accessToken string, callerID string) *KickflowClient {
	return &KickflowClient{
		BaseURL:     baseURL,
		AccessToken: accessToken,
		CallerID:    callerID,
		HTTPClient:  http.DefaultClient,
	}
}

func (c *KickflowClient) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Caller-ID", c.CallerID)
	req.Header.Set("Content-Type", "application/json")
	return c.HTTPClient.Do(req)
}
