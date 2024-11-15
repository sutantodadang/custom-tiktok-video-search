package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"tiktok-video-search/common/caches"
	"tiktok-video-search/common/constants"
	h "tiktok-video-search/common/request"
	"tiktok-video-search/models"
	"time"

	"github.com/rs/zerolog/log"
)

type VideoHandler struct {
	cache *caches.InMemoryCache
}

// HTTP Handlers
func NewVideoHandler(cache *caches.InMemoryCache) *VideoHandler {
	return &VideoHandler{cache: cache}
}

func (c *VideoHandler) FetchVideos(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "Keyword is required", http.StatusBadRequest)
		return
	}

	offset := r.URL.Query().Get("offset")
	if keyword == "" {
		http.Error(w, "Offset is required", http.StatusBadRequest)
		return
	}

	myCookie, err := c.cache.Get(constants.KEY_COOKIE)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Error().Err(err).Send()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	if myCookie == nil {

		cookies, err := h.HttpRequestCookie()
		if err != nil {
			log.Error().Err(err).Send()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}

		err = c.cache.Set(constants.KEY_COOKIE, cookies, constants.TIME_EXPIRED)
		if err != nil {
			return
		}

		myCookie = cookies
	}

	valueCookie, ok := myCookie.(map[string]string)
	if ok {
		ttl, err := time.Parse(time.RFC3339, valueCookie["ttl"])
		if err != nil {
			log.Error().Err(err).Send()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}

		if ttl.Before(time.Now()) {
			cookies, err := h.HttpRequestCookie()
			if err != nil {
				log.Error().Err(err).Send()
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(err)
				return
			}

			err = c.cache.Set(constants.KEY_COOKIE, cookies, constants.TIME_EXPIRED)
			if err != nil {
				return
			}

			valueCookie = cookies

		}
	}

	result, err := h.Httprequest(fmt.Sprintf("https://www.tiktok.com/api/search/general/full/?keyword=%s&offset=%s", keyword, offset), valueCookie)
	if err != nil {
		log.Error().Err(err).Send()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	vid, err := models.DecodeJson(result, c.cache)
	if err != nil {
		log.Error().Err(err).Send()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	for name, v := range valueCookie {

		if name == "ttl" {
			continue
		}

		http.SetCookie(w, &http.Cookie{Name: name, Value: v, Path: "/", Domain: "www.tiktok.com"})

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vid)
	return
}

func (c *VideoHandler) FetchVideoDetail(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	myCookie, err := c.cache.Get(constants.KEY_COOKIE)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Error().Err(err).Send()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	myVideo, err := c.cache.Get(constants.KEY_VIDEO + id)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Error().Err(err).Send()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	if myCookie == nil {

		cookies, err := h.HttpRequestCookie()
		if err != nil {
			log.Error().Err(err).Send()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}

		err = c.cache.Set(constants.KEY_COOKIE, cookies, constants.TIME_EXPIRED)
		if err != nil {
			return
		}

		myCookie = cookies
	}

	valueCookie, ok := myCookie.(map[string]string)
	if ok {
		ttl, err := time.Parse(time.RFC3339, valueCookie["ttl"])
		if err != nil {
			log.Error().Err(err).Send()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}

		if ttl.Before(time.Now()) {
			cookies, err := h.HttpRequestCookie()
			if err != nil {
				log.Error().Err(err).Send()
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(err)
				return
			}

			err = c.cache.Set(constants.KEY_COOKIE, cookies, constants.TIME_EXPIRED)
			if err != nil {
				return
			}

			valueCookie = cookies

		}
	}

	// var result []byte
	val, ok := myVideo.(models.VideoData)

	log.Info().Interface("value ", val).Send()

	if !ok {

		teast, _ := json.Marshal(val)
		log.Info().Msg(string(teast))

		log.Error().Err(errors.New("no cache")).Send()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return

	}

	// Make a request to the TikTok video URL
	resp, err := http.Get(val.URL)
	if err != nil {
		http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the headers from TikTok's response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the correct content type
	w.Header().Set("Content-Type", "video/mp4")

	// Stream the video content to the client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func (c *VideoHandler) VideoProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the video URL from the query parameter
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		http.Error(w, "Missing video URL", http.StatusBadRequest)
		return
	}

	myCookie, err := c.cache.Get(constants.KEY_COOKIE)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Error().Err(err).Send()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	if myCookie == nil {
		// log.Info().Interface("cookies nil", myCookie).Send()
		cookies, err := h.HttpRequestCookie()
		if err != nil {
			log.Error().Err(err).Send()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}

		err = c.cache.Set(constants.KEY_COOKIE, cookies, constants.TIME_EXPIRED)
		if err != nil {
			return
		}

		myCookie = cookies
	}

	valueCookie, ok := myCookie.(map[string]string)
	if ok {
		ttl, err := time.Parse(time.RFC3339, valueCookie["ttl"])
		if err != nil {
			log.Error().Err(err).Send()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}

		if ttl.Before(time.Now()) {
			cookies, err := h.HttpRequestCookie()
			if err != nil {
				log.Error().Err(err).Send()
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(err)
				return
			}

			err = c.cache.Set(constants.KEY_COOKIE, cookies, constants.TIME_EXPIRED)
			if err != nil {
				return
			}

			valueCookie = cookies

		}
	}

	// Make a request to the TikTok video URL
	req, err := http.NewRequest(http.MethodGet, videoURL, nil)
	if err != nil {
		http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
		return
	}

	// Copy the headers from TikTok's response
	for name, values := range valueCookie {

		req.Header.Add(name, values)

	}

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// Set the correct content type
	w.Header().Set("Content-Type", "video/mp4")

	// Stream the video content to the client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}
}
