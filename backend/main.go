package main

import (
	"net/http"
	"tiktok-video-search/common/caches"

	h "tiktok-video-search/api/http"
	"tiktok-video-search/middleware"

	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.With().Caller().Logger()

}

func main() {
	rdb := caches.NewInMemoryCache()

	videoHandler := h.NewVideoHandler(rdb)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/videos", videoHandler.FetchVideos)
	mux.HandleFunc("/api/videos/detail/{id}", videoHandler.FetchVideoDetail)
	mux.HandleFunc("/api/videos/video-proxy", videoHandler.VideoProxyHandler)
	log.Info().Msg("Server started at :5555")

	handler := middleware.CorsMiddleware(mux)

	log.Fatal().Err(http.ListenAndServe(":5555", handler)).Send()
}
