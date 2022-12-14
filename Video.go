package go_tiktok

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const maxCount int32 = 20

type Video struct {
	Id               string `json:"id"`
	CreateTime       int64  `json:"create_time"`
	CoverImageUrl    string `json:"cover_image_url"`
	ShareUrl         string `json:"share_url"`
	VideoDescription string `json:"video_description"`
	Duration         int32  `json:"duration"`
	Height           int32  `json:"height"`
	Width            int32  `json:"width"`
	Title            string `json:"title"`
	EmbedHtml        string `json:"embed_html"`
	EmbedLink        string `json:"embed_link"`
	LikeCount        int32  `json:"like_count"`
	CommentCount     int32  `json:"comment_count"`
	ShareCount       int32  `json:"share_count"`
	ViewCount        int32  `json:"view_count"`
}

type VideoResponse struct {
	Data struct {
		Videos  []Video `json:"videos"`
		HasMore bool    `json:"has_more"`
		Cursor  *int64  `json:"cursor"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		LogId   string `json:"log_id"`
	} `json:"error"`
}

type ListVideosConfig struct {
	Fields string
	From   *time.Time
	To     *time.Time
}

func (service *Service) ListVideos(cfg *ListVideosConfig) (*[]Video, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var values = url.Values{}
	values.Set("fields", cfg.Fields)

	var videos []Video

	var cursor *int64 = nil

	if cfg.To != nil {
		cursor_ := cfg.To.UnixMilli()
		cursor = &cursor_
	}

	for {
		body := struct {
			Cursor   *int64 `json:"cursor"`
			MaxCount int32  `json:"max_count"`
		}{
			cursor,
			maxCount,
		}

		videoResponse := VideoResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.url(fmt.Sprintf("video/list/?%s", values.Encode())),
			BodyModel:     body,
			ResponseModel: &videoResponse,
		}
		_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		if len(videoResponse.Data.Videos) == 0 {
			break
		}

		for _, video := range videoResponse.Data.Videos {
			if cfg.From != nil {
				if cfg.From.Unix() > video.CreateTime {
					goto ready
				}
			}

			videos = append(videos, video)
		}

		if !videoResponse.Data.HasMore {
			break
		}
		if videoResponse.Data.Cursor == nil {
			return nil, errortools.ErrorMessage("Cursor is nil although HasMore is true")
		}

		cursor = videoResponse.Data.Cursor
	}
ready:
	return &videos, nil
}
