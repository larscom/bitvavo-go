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
	httpConfg *httpConfig,
	authConfig *authConfig,
) (T, error) {
	req, _ := http.NewRequestWithContext(ctx, "DELETE", createRequestUrl(url, params), nil)
	return httpDo[T](req, make([]byte, 0), httpConfg, authConfig)
}

func httpGet[T any](
	ctx context.Context,
	url string,
	params url.Values,
	httpConfg *httpConfig,
	authConfig *authConfig,
) (T, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", createRequestUrl(url, params), nil)
	return httpDo[T](req, make([]byte, 0), httpConfg, authConfig)
}

func httpPost[T any](
	ctx context.Context,
	url string,
	body any,
	params url.Values,
	httpConfig *httpConfig,
	authConfig *authConfig,
) (T, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		var empty T
		return empty, err
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", createRequestUrl(url, params), bytes.NewBuffer(payload))
	return httpDo[T](req, payload, httpConfig, authConfig)
}

func httpPut[T any](
	ctx context.Context,
	url string,
	body any,
	params url.Values,
	httpConfig *httpConfig,
	authConfig *authConfig,
) (T, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		var empty T
		return empty, err
	}

	req, _ := http.NewRequestWithContext(ctx, "PUT", createRequestUrl(url, params), bytes.NewBuffer(payload))
	return httpDo[T](req, payload, httpConfig, authConfig)
}

func httpDo[T any](
	request *http.Request,
	body []byte,
	httpConfig *httpConfig,
	authConfig *authConfig,
) (T, error) {
	var empty T
	if err := setHeaders(request, body, authConfig); err != nil {
		return empty, err
	}

	debug(httpConfig.printer, fmt.Sprint("http request ", request.Method, " url=", request.URL.String()))

	response, err := httpConfig.client.Do(request)
	if err != nil {
		debug(httpConfig.printer, fmt.Sprint("http response error: ", err))
		return empty, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if err := updateRateLimits(response, httpConfig); err != nil {
		return empty, err
	}

	if response.StatusCode > http.StatusIMUsed {
		return empty, unwrapErr(response, httpConfig.printer)
	}

	return unwrapBody[T](response, httpConfig.printer)
}

func unwrapBody[T any](response *http.Response, printer DebugPrinter) (T, error) {
	var data T
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	debug(printer, fmt.Sprint("http response body: ", string(b), " statusCode=", response.StatusCode))

	if err := json.Unmarshal(b, &data); err != nil {
		return data, err
	}

	return data, nil
}

func unwrapErr(response *http.Response, printer DebugPrinter) error {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	debug(printer, fmt.Sprint("http response body: ", string(b), " statusCode=", response.StatusCode))

	var apiError *ApiError
	if err := json.Unmarshal(b, &apiError); err != nil {
		return ErrNOKResponse(response.StatusCode, b)
	}
	return apiError
}

func updateRateLimits(
	response *http.Response,
	httpConfig *httpConfig,
) error {
	for key, value := range response.Header {
		if key == headerRatelimit {
			if len(value) == 0 {
				return ErrHeaderNoValue(headerRatelimit)
			}
			rateLimit := util.MustInt64(value[0])
			debug(httpConfig.printer, fmt.Sprint("http rate limit is currently: ", rateLimit))
			httpConfig.updateRateLimit(rateLimit)
		}
		if key == headerRatelimitResetAt {
			if len(value) == 0 {
				return ErrHeaderNoValue(headerRatelimitResetAt)
			}
			httpConfig.updateRateLimitResetAt(time.UnixMilli(util.MustInt64(value[0])))
		}
	}
	return nil
}

func setHeaders(request *http.Request, body []byte, authConfig *authConfig) error {
	if authConfig == nil {
		return nil
	}

	timestamp := time.Now().UnixMilli()

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(headerAccessKey, authConfig.apiKey)
	request.Header.Set(headerAccessSignature, crypto.CreateSignature(request.Method, strings.Replace(request.URL.String(), apiURL, "", 1), body, timestamp, authConfig.apiSecret))
	request.Header.Set(headerAccessTimestamp, fmt.Sprint(timestamp))
	request.Header.Set(headerAccessWindow, fmt.Sprint(authConfig.windowTime))

	return nil
}

func createRequestUrl(url string, params url.Values) string {
	return util.IfOrElse(len(params) > 0, func() string { return fmt.Sprintf("%s?%s", url, params.Encode()) }, url)
}
