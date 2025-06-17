package supabase

import (
	"fmt"
	"net/http"
	"os"
)

type SupabaseClient struct {
	Url     string
	ApiKey  string
	Client  *http.Client
	Headers http.Header
}

func NewSupabaseClient() *SupabaseClient {
	url := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")

	headers := http.Header{}
	headers.Set("apikey", apiKey)
	headers.Set("Authorization", "Bearer "+apiKey)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")

	return &SupabaseClient{
		Url:     url,
		ApiKey:  apiKey,
		Client:  &http.Client{},
		Headers: headers,
	}
}

func (s *SupabaseClient) Ping() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/v1/", s.Url), nil)
	if err != nil {
		return err
	}
	req.Header = s.Headers

	res, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode != 404 {
		return fmt.Errorf("ping failed with status: %d", res.StatusCode)
	}

	return nil
}
