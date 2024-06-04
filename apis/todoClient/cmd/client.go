package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrConnection      = errors.New("connection Error")
	ErrNotFound        = errors.New("not found")
	ErrInvalidResponse = errors.New("invalid response")
	ErrInvalidData     = errors.New("invalid data")
	ErrNaN             = errors.New("not a number")
)

type item struct {
	Task        string    `json:"task"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type response struct {
	Results      []item `json:"results"`
	Date         int    `json:"date"`
	TotalResults int    `json:"total_results"`
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func getItems(url string) ([]item, error) {
	r, err := newClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrConnection)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("connot read the body: %w", err)
		}

		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
	}

	var resp response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("%w: Invalid json response", err)
	}

	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("%w: No results found", ErrNotFound)
	}

	return resp.Results, nil
}

func getAll(apiRoot string) ([]item, error) {
	u := fmt.Sprintf("%s/todo", apiRoot)
	return getItems(u)
}

func getOne(apiRoot string, id int) (item, error) {
	u := fmt.Sprintf("%s/todo/%d", apiRoot, id)
	items, err := getItems(u)
	if err != nil {
		return item{}, err
	}
	if len(items) != 1 {
		return item{}, fmt.Errorf("%w: Invalid result", ErrInvalidData)
	}
	return items[0], nil
}
