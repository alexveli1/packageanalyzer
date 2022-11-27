// Package httpv1 gathers web data
package httpv1

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"

	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

// ITransporter interface for client, which might [add setup required] get info from different protocols, clients
type ITransporter interface {
	GetRepo(ctx context.Context, name string) ([]domain.Binpack, error)
}

// Client implementation of ITransporter interface serving HTTP connection with resty client
type Client struct {
	Addr   string
	client *resty.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		Addr: cfg.Address,
		client: resty.New().
			SetRetryCount(0).
			SetTimeout(cfg.Timeout).
			SetLogger(mylog.SugarLogger),
	}
}

// GetRepo connects to api and downloads JSON response
func (c *Client) GetRepo(ctx context.Context, branch string) ([]domain.Binpack, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	//result := *getContents(branch)
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
	if err != nil {
		mylog.SugarLogger.Warnf("cannot unmarshall body")

		return nil, err
	}

	return result.Packages, nil
}

// getContents used for testing to avoid excessive load on server
/*func getContents(branch string) *domain.RequestResult {
	f, _ := os.Open("/home/alex/Documents/GoLang/Basalt/" + branch + ".json")
	data, _ := io.ReadAll(f)
	var v domain.RequestResult
	_ = json.Unmarshal(data, &v)
	return &v
}
*/
