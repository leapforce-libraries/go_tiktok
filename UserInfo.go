package go_tiktok

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type UserInfo struct {
	OpenId          string `json:"open_id"`
	UnionId         string `json:"union_id"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarUrl100    string `json:"avatar_url_100"`
	AvatarLargeUrl  string `json:"avatar_large_url"`
	DisplayName     string `json:"display_name"`
	BioDescription  string `json:"bio_description"`
	ProfileDeepLink string `json:"profile_deep_link"`
	IsVerified      bool   `json:"is_verified"`
	FollowerCount   int64  `json:"follower_count"`
	FollowingCount  int64  `json:"following_count"`
	LikesCount      int64  `json:"likes_count"`
}

func (service *Service) GetUserInfo(fields string) (*UserInfo, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var values url.Values
	values.Set("fields", fields)

	userInfo := UserInfo{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("user/info/?%s", values.Encode())),
		ResponseModel: &userInfo,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &userInfo, nil
}
