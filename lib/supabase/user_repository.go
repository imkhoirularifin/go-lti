package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	
	"io/ioutil"
	"net/http"
)

type SupabaseUser struct {
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	WishStatus   string `json:"wish_status,omitempty"` 
	CanvasUserID string `json:"canvas_user_id"`
	WishLink	string `json:"wish_link,omitempty"`
}

func (s *SupabaseClient) InsertUser(user SupabaseUser) error {
	bodyBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.Url+"/rest/v1/users", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header = s.Headers

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to insert user, status: %d", resp.StatusCode)
	}

	return nil
}
func (s *SupabaseClient) GetUserByEmail(email string) (*SupabaseUser, error) {
	url := fmt.Sprintf("%s/rest/v1/users?email=eq.%s&select=*", s.Url, email)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", s.ApiKey)
	req.Header.Set("Authorization", "Bearer "+s.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request Supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		var users []SupabaseUser
		if err := json.Unmarshal(body, &users); err != nil {
			return nil, fmt.Errorf("failed to parse Supabase response: %w", err)
		}

		if len(users) > 0 {
			return &users[0], nil
		}

		return nil, nil 
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return nil, fmt.Errorf("failed to get user: status %d, body: %s", resp.StatusCode, string(body))
}

func (s *SupabaseClient) CheckUserExistsByEmail(email string) (bool, error) {
	url := fmt.Sprintf("%s/rest/v1/users?email=eq.%s", s.Url, email)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header = s.Headers

	resp, err := s.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Supabase GET failed: status %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var results []map[string]interface{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return false, err
	}

	return len(results) > 0, nil
}

