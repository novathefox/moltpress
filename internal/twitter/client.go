package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrTweetNotFound = errors.New("tweet not found")
	ErrInvalidURL    = errors.New("invalid tweet url")
	ErrFetchFailed   = errors.New("failed to fetch tweet")
)

type Tweet struct {
	ID             string
	Text           string
	AuthorUsername string
}

type syndicationTweet struct {
	IDStr    string `json:"id_str"`
	Text     string `json:"text"`
	FullText string `json:"full_text"`
	User     struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

func FetchTweet(tweetURL string) (*Tweet, error) {
	ID, err := parseTweetID(tweetURL)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	endpoint := fmt.Sprintf("https://cdn.syndication.twimg.com/tweet-result?id=%s&token=a", ID)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("failed to build tweet request", "error", err, "tweet_id", ID)
		return nil, fmt.Errorf("%w", ErrFetchFailed)
	}

	response, err := client.Do(request)
	if err != nil {
		slog.Error("failed to fetch tweet", "error", err, "tweet_id", ID)
		return nil, fmt.Errorf("%w", ErrFetchFailed)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrTweetNotFound
	}

	if response.StatusCode != http.StatusOK {
		slog.Error("unexpected tweet response status", "status", response.StatusCode, "tweet_id", ID)
		return nil, fmt.Errorf("%w", ErrFetchFailed)
	}

	var payload syndicationTweet
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		slog.Error("failed to decode tweet response", "error", err, "tweet_id", ID)
		return nil, fmt.Errorf("%w", ErrFetchFailed)
	}

	if payload.IDStr == "" || payload.User.ScreenName == "" {
		return nil, ErrTweetNotFound
	}

	text := payload.Text
	if text == "" {
		text = payload.FullText
	}

	return &Tweet{
		ID:             payload.IDStr,
		Text:           text,
		AuthorUsername: payload.User.ScreenName,
	}, nil
}

func parseTweetID(tweetURL string) (string, error) {
	parsed, err := url.Parse(tweetURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	host := strings.ToLower(parsed.Host)
	if host == "" || !isTwitterHost(host) {
		return "", ErrInvalidURL
	}

	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	for i, segment := range segments {
		if segment == "status" && i+1 < len(segments) {
			id := segments[i+1]
			if id == "" || !isNumeric(id) {
				return "", ErrInvalidURL
			}
			return id, nil
		}
	}

	return "", ErrInvalidURL
}

func isTwitterHost(host string) bool {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	if host == "twitter.com" || strings.HasSuffix(host, ".twitter.com") {
		return true
	}

	return host == "x.com" || strings.HasSuffix(host, ".x.com")
}

func isNumeric(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return value != ""
}
