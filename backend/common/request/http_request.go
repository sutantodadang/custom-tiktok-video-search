package request

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"tiktok-video-search/common/constants"
	"time"

	"github.com/rs/zerolog/log"
)

var defaultHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
	"Accept-Language": "en-US,en;q=0.5",
}

func Httprequest(url string, cookiers map[string]string) (c []byte, err error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	for name, v := range cookiers {
		req.AddCookie(&http.Cookie{Name: name, Value: v})
	}

	for name, value := range defaultHeaders {
		req.Header.Set(name, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	c, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	return
}

func HttpRequestCookie() (c map[string]string, err error) {
	url := "https://www.tiktok.com/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for name, value := range defaultHeaders {
		req.Header.Set(name, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	c = make(map[string]string)

	for _, cookie := range resp.Cookies() {
		c[cookie.Name] = cookie.Value
	}

	c["ttl"] = time.Now().Add(constants.TIME_EXPIRED).Format(time.RFC3339)

	return
}
