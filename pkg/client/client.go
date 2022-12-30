package client

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	Url string
}

func (c Client) Get() ([]byte, error) {
	resp, err := http.Get(c.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c Client) GetMetric(metric string) (*string, error) {
	body, err := c.Get()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), metric) {
			val := strings.TrimLeft(scanner.Text(), metric)
			return &val, nil
		}
	}

	return nil, errors.New("unable to find the requested metric")
}
