package httpv1

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"

	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

type ITransporter interface {
	GetRepo(ctx context.Context, name string) ([]domain.Binpack, error)
}

type Client struct {
	Addr   string
	client *resty.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		Addr: cfg.Address,
		client: resty.New().
			SetRetryCount(0).SetLogger(mylog.SugarLogger),
	}
}

func (c *Client) GetRepo(ctx context.Context, branch string) ([]domain.Binpack, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var result domain.RequestResult
	resp, err := c.client.R().
		SetHeader("content-type", "text/plain").
		SetContext(ctx).
		Get(c.Addr + branch)
	if err != nil {
		mylog.SugarLogger.Warnf("Cannot initiate request: %v", err)

		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &result)

	return result.Packages, err
}
