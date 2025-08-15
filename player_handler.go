package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

type PlayerHandler struct {
	template            *template.Template
	streamStatusManager *StreamStatusManager
}

func NewPlayerHandler(streamStatusManager *StreamStatusManager) PlayerHandler {
	return PlayerHandler{
		readTemplate("player.gohtml"),
		streamStatusManager,
	}
}

func (handler *PlayerHandler) Handle(writer http.ResponseWriter, request *http.Request, mappingResult PathMappingResult) {
	streamInfo := handler.streamStatusManager.StreamInfo(mappingResult.CalculatedPath)

	isStream := request.URL.Query()["stream"] != nil
	if isStream {
		streamStatus := streamInfo.DominantStatusCode()
		switch streamStatus {
		case NoStream:
		case StreamTranscodingFailed:
		case StreamInPreparation:
			RelativeRedirect(writer, request, "?stream", http.StatusSeeOther)
			return
		}
	}

	writer.Header().Add("Content-Type", "text/html; charset=utf-8")

	encodedUrl := (&url.URL{Path: mappingResult.UrlPath}).String()

	playbackUrl := encodedUrl
	if isStream {
		playbackUrl += "?stream&playlist"
	}

	dir, file := filepath.Split(mappingResult.UrlPath)
	if err := handler.template.Execute(writer, struct {
		Dir         string
		File        string
		Url         string
		PlaybackUrl string
	}{
		dir,
		file,
		encodedUrl,
		playbackUrl,
	}); err != nil {
		log.Printf("Template-Formatting failed: %s", err)
	}
}
