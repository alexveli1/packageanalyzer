package httpv1

import (
	"context"

	"github.com/go-resty/resty/v2"

	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

type ITransporter interface {
	GetRepo(ctx context.Context, name string) ([]byte, error)
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

func (c *Client) GetRepo(ctx context.Context, name string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	req := c.client.SetBaseURL(c.Addr)
	_, err := req.R().
		SetHeader("content-type", "text/plain").
		SetContext(ctx).
		Get(name)
	if err != nil {
		mylog.SugarLogger.Warnf("Cannot initiate request: %v", err)

		return nil, err
	}
	return nil, nil
}
