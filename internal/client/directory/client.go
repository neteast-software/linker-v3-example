package directory

import (
	"context"
	"errors"
	"net/url"

	httpclient "github.com/neteast-software/go-module/http/client"
)

type Client struct {
	http *httpclient.Client
}

func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}

func (p *Client) Badge(ctx context.Context, userID string) (Badge, error) {
	if p == nil || p.http == nil {
		return Badge{}, errors.New("用户目录 HTTP client 未初始化")
	}
	if userID == "" {
		return Badge{}, errors.New("用户 ID 不能为空")
	}
	result, err := p.http.Do(ctx, httpclient.GET("/users/"+url.PathEscape(userID)+"/badge"))
	if err != nil {
		return Badge{}, err
	}
	var response badgeResponse
	if err = result.DecodeJSON(&response); err != nil {
		return Badge{}, err
	}
	return response.Badge(), nil
}
