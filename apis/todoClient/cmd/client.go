package cmd

import (
	"bytes"
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

func sendRequest(url, method, contentType string, expStatus int, body io.Reader) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	r, err := newClient().Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != expStatus {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("cannot read the body: %w", err)
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w: %s", err, msg)
	}
	return nil
}

func addItem(apiRoot, task string) error {
	u := fmt.Sprintf("%s/todo", apiRoot)
	item := struct {
		Task string `json:"task"`
	}{
		Task: task,
	}

	body := bytes.Buffer{}
	if err := json.NewEncoder(&body).Encode(item); err != nil {
		return err
	}

	return sendRequest(u, http.MethodPost, "application/json", http.StatusCreated, &body)
}

func completeItem(apiRoot string, id int) error {
	u := fmt.Sprintf("%s/todo/%d?complete", apiRoot, id)
	return sendRequest(u, http.MethodPatch, "", http.StatusNoContent, nil)
}

func deleteItem(apiRoot string, id int) error {
	u := fmt.Sprintf("%s/todo/%d", apiRoot, id)
	return sendRequest(u, http.MethodDelete, "", http.StatusNoContent, nil)
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
