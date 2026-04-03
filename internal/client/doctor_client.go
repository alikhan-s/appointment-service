package client

import (
	"context"
	"net/http"
	"time"

	"github.com/alikhan-s/appointment-s/internal/model"
)

type DoctorClient interface {
	CheckDoctorExists(ctx context.Context, doctorID string) error
}

type doctorHTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewDoctorHTTPClient(baseURL string) DoctorClient {
	return &doctorHTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *doctorHTTPClient) CheckDoctorExists(ctx context.Context, doctorID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/doctors/"+doctorID, nil)
	if err != nil {
		return model.ErrDoctorServiceDown
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return model.ErrDoctorServiceDown
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return model.ErrDoctorNotExists
	}
	if resp.StatusCode != http.StatusOK {
		return model.ErrDoctorServiceDown
	}

	return nil
}
