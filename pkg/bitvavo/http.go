package bitvavo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/crypto"
	"github.com/larscom/bitvavo-go/internal/util"
)

var (
	ErrHeaderNoValue = func(h string) error { return fmt.Errorf("header: %s didn't contain a value", h) }

	ErrNOKResponse = func(code int, b []byte) error {
		return fmt.Errorf("did not get OK response, code=%d, body=%s", code, string(b))
	}
)

func httpDelete[T any](
	ctx context.Context,
	url string,
	params url.Values,
	updateRateLimit func(ratelimit int64),
	updateRateLimitResetAt func(resetAt time.Time),
	config *privateConfig,
) (T, error) {
	req, _ := http.NewRequestWithContext(ctx, "DELETE", createRequestUrl(url, params), nil)
	return httpDo[T](req, make([]byte, 0), updateRateLimit, updateRateLimitResetAt, config)
}

func httpGet[T any](
	ctx context.Context,
	url string,
	params url.Values,
	updateRateLimit func(ratelimit int64),
	updateRateLimitResetAt func(resetAt time.Time),
	config *privateConfig,
) (T, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", createRequestUrl(url, params), nil)
	return httpDo[T](req, make([]byte, 0), updateRateLimit, updateRateLimitResetAt, config)
}

func httpPost[T any](
	ctx context.Context,
	url string,
	body any,
	params url.Values,
	updateRateLimit func(ratelimit int64),
	updateRateLimitResetAt func(resetAt time.Time),
	config *privateConfig,
) (T, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		var empty T
		return empty, err
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", createRequestUrl(url, params), bytes.NewBuffer(payload))
	return httpDo[T](req, payload, updateRateLimit, updateRateLimitResetAt, config)
}

func httpPut[T any](
	ctx context.Context,
	url string,
	body any,
	params url.Values,
	updateRateLimit func(ratelimit int64),
	updateRateLimitResetAt func(resetAt time.Time),
	config *privateConfig,
) (T, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		var empty T
		return empty, err
	}

	req, _ := http.NewRequestWithContext(ctx, "PUT", createRequestUrl(url, params), bytes.NewBuffer(payload))
	return httpDo[T](req, payload, updateRateLimit, updateRateLimitResetAt, config)
}

func httpDo[T any](
	request *http.Request,
	body []byte,
	updateRateLimit func(ratelimit int64),
	updateRateLimitResetAt func(resetAt time.Time),
	config *privateConfig,
) (T, error) {

	var empty T
	if err := setHeaders(request, body, config); err != nil {
		return empty, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return empty, err
	}
	defer response.Body.Close()

	if err := updateRateLimits(response, updateRateLimit, updateRateLimitResetAt); err != nil {
		return empty, err
	}

	if response.StatusCode > http.StatusIMUsed {
		return empty, unwrapErr(response)
	}

	return unwrapBody[T](response)
}

func unwrapBody[T any](response *http.Response) (T, error) {
	var data T
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return data, err
	}

	return data, nil
}

func unwrapErr(response *http.Response) error {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var apiError *ApiError
	if err := json.Unmarshal(b, &apiError); err != nil {
		return ErrNOKResponse(response.StatusCode, b)
	}
	return apiError
}

func updateRateLimits(
	response *http.Response,
	updateRateLimit func(ratelimit int64),
	updateRateLimitResetAt func(resetAt time.Time),
) error {
	for key, value := range response.Header {
		if key == headerRatelimit {
			if len(value) == 0 {
				return ErrHeaderNoValue(headerRatelimit)
			}
			updateRateLimit(util.MustInt64(value[0]))
		}
		if key == headerRatelimitResetAt {
			if len(value) == 0 {
				return ErrHeaderNoValue(headerRatelimitResetAt)
			}
			updateRateLimitResetAt(time.UnixMilli(util.MustInt64(value[0])))
		}
	}
	return nil
}

func setHeaders(request *http.Request, body []byte, config *privateConfig) error {
	if config == nil {
		return nil
	}

	timestamp := time.Now().UnixMilli()

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(headerAccessKey, config.apiKey)
	request.Header.Set(headerAccessSignature, crypto.CreateSignature(request.Method, strings.Replace(request.URL.String(), bitvavoURL, "", 1), body, timestamp, config.apiSecret))
	request.Header.Set(headerAccessTimestamp, fmt.Sprint(timestamp))
	request.Header.Set(headerAccessWindow, fmt.Sprint(config.windowTime))

	return nil
}

func createRequestUrl(url string, params url.Values) string {
	return util.IfOrElse(len(params) > 0, func() string { return fmt.Sprintf("%s?%s", url, params.Encode()) }, url)
}
