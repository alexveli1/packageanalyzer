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
	GetOfficialDiff(ctx context.Context, branch1, branch2 string) (*domain.VerificationInfo, error)
}

// Client implementation of ITransporter interface serving HTTP connection with resty client
type Client struct {
	ExEndpt  string
	VrfEndpt string
	client   *resty.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		ExEndpt:  cfg.ExportEndpoint,
		VrfEndpt: cfg.VerificationEndpoint,
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
	// result := *getBranchContent(branch)
	var result domain.RequestResult
	resp, err := c.client.R().
		SetHeader("content-type", "text/plain").
		SetContext(ctx).
		Get(c.ExEndpt + branch)
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

func (c *Client) GetOfficialDiff(ctx context.Context, branch1, branch2 string) (*domain.VerificationInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	//compareResult := *getCompareResults(branch1)
	var compareResult domain.VerificationInfo
	resp, err := c.client.R().
		SetHeader("content-type", "text/plain").
		SetContext(ctx).
		Get(c.VrfEndpt + "pkgset1=" + branch1 + "&pkgset2=" + branch2)
	if err != nil {
		mylog.SugarLogger.Warnf("Cannot initiate request: %v", err)

		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &compareResult)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot unmarshall body")

		return nil, err
	}

	return &compareResult, nil
}

// getBranchContent used for testing to avoid excessive load on server
/*func getBranchContent(branch string) *domain.RequestResult {
	f, _ := os.Open("test/data/" + branch + ".json")
	data, _ := io.ReadAll(f)
	var v domain.RequestResult
	err := json.Unmarshal(data, &v)
	if err != nil {
		mylog.SugarLogger.Warnf("Cannot marshal data: %v", err)
	}

	return &v
}*/

/*func getCompareResults(branch string) *domain.VerificationInfo {
	f, _ := os.Open("test/data/packagescompare_" + branch + ".json")
	data, _ := io.ReadAll(f)
	var v domain.VerificationInfo
	err := json.Unmarshal(data, &v)
	if err != nil {
		mylog.SugarLogger.Warnf("Cannot marshal data: %v", err)
	}

	return &v
}*/
