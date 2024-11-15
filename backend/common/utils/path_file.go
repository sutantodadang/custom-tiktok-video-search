package utils

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func FileReadHelper() (filePath string, err error) {

	dir, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	log.Info().Msg(dir)

	err = os.MkdirAll(filepath.Join(dir, "cookies"), os.ModePerm)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	var file *os.File

	defer file.Close()

	filePath = filepath.Join(dir, "cookies", "cookie.json")

	return

}
