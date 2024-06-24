package daisysms

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	http.Client
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) Balance() (float64, error) {
	params := url.Values{
		"api_key": {c.apiKey},
		"action":  {"getBalance"},
	}
	res, err := c.Get("https://daisysms.com/stubs/handler_api.php?" + params.Encode())
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(strings.TrimPrefix(string(body), "ACCESS_BALANCE:"), 64)
}

func (c *Client) GetNumber(service Service, maxPrice ...float64) (string, string, error) {
	params := url.Values{
		"api_key": {c.apiKey},
		"action":  {"getNumber"},
		"service": {string(service)},
	}
	if len(maxPrice) > 0 {
		params.Set("maxPrice", fmt.Sprintf("%.2f", maxPrice[0]))
	}
	res, err := c.Get("https://daisysms.com/stubs/handler_api.php?" + params.Encode())
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(string(body), ":")
	if len(parts) == 1 {
		switch parts[0] {
		case "NO_NUMBERS":
			return "", "", ErrNoNumbers
		case "NO_BALANCE":
			return "", "", ErrNoMoney
		case "MAX_PRICE_EXCEEDED":
			return "", "", ErrMaxPriceExceeded
		case "TOO_MANY_ACTIVE_RENTALS":
			return "", "", ErrTooManyActiveRentals
		default:
			return "", "", fmt.Errorf("unknown error: %s", parts[0])
		}
	}

	return parts[1], parts[2], nil
}

func (c *Client) Wait(ctx context.Context, id string) (string, error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			params := url.Values{
				"api_key": {c.apiKey},
				"action":  {"getStatus"},
				"id":      {id},
			}
			res, err := c.Get("https://daisysms.com/stubs/handler_api.php?" + params.Encode())
			if err != nil {
				return "", err
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return "", err
			}

			parts := strings.Split(string(body), ":")
			if len(parts) == 1 {
				switch parts[0] {
				case "STATUS_WAIT_CODE":
					continue
				case "NO_ACTIVATION":
					return "", ErrWrongID
				case "STATUS_CANCEL":
					return "", ErrRentalCanceled
				default:
					return "", fmt.Errorf("unknown error: %s", parts[0])
				}
			}

			return parts[1], nil
		}
	}
}

func (c *Client) Done(id string) error {
	params := url.Values{
		"api_key": {c.apiKey},
		"action":  {"setStatus"},
		"id":      {id},
		"status":  {"6"},
	}
	res, err := c.Get("https://daisysms.com/stubs/handler_api.php?" + params.Encode())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	switch string(body) {
	case "ACCESS_ACTIVATION":
		return nil
	case "NO_ACTIVATION":
		return ErrWrongID
	default:
		return fmt.Errorf("unknown error: %s", body)
	}
}

func (c *Client) Cancel(id string) error {
	params := url.Values{
		"api_key": {c.apiKey},
		"action":  {"setStatus"},
		"id":      {id},
		"status":  {"8"},
	}
	res, err := c.Get("https://daisysms.com/stubs/handler_api.php?" + params.Encode())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	switch string(body) {
	case "ACCESS_CANCEL":
		return nil
	case "NO_ACTIVATION":
		return ErrWrongID
	default:
		return fmt.Errorf("unknown error: %s", body)
	}
}
