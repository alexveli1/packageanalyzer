package httpv1

import (
	"context"
	"encoding/json"
	"io"
	"os"

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
	result = *getContents(branch)
	// TODO reenable after testing
	/*	req := c.client.SetBaseURL(c.Addr)
		_, err := req.R().
			SetHeader("content-type", "text/plain").
			SetContext(ctx).
			SetResult(result).
			Get(branch)
		if err != nil {
			mylog.SugarLogger.Warnf("Cannot initiate request: %v", err)

			return nil, err
		}*/

	return result.Packages, nil
}

func getContents(branch string) *domain.RequestResult {
	f, _ := os.Open("/home/alex/Documents/GoLang/Basalt/" + branch + ".json")
	data, _ := io.ReadAll(f)
	var v domain.RequestResult
	_ = json.Unmarshal(data, &v)
	return &v
}
