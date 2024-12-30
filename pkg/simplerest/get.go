package simplerest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrNotOK = errors.New("response is not a HTTP OK")

func Get[T any](url string) (T, error) {
	var result T

	resp, err := http.Get(url)
	if err != nil {
		return result, fmt.Errorf("HTTP error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, ErrNotOK
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("could not read response: %w", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("invalid json in response: %w", err)
	}

	return result, nil
}
