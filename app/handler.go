package app

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"github.com/tfonfara/plexsmarthome/helper"
	"github.com/tfonfara/plexsmarthome/plex"
)

type EventHandler struct {
	configuration *Configuration
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		NewConfiguration(),
	}
}

func (h *EventHandler) WebhookHandler(p plex.WebhookPayload, raw string) {
	log.Printf("Received webhook for %s (%s) from %s (Player: %s, Account: %d)", p.Metadata.Type, p.Event, p.Player.PublicAddress, p.Player.UUID, p.Account.ID)

	if helper.Contains(h.configuration.Players, p.Player.UUID) {
		if helper.Contains(h.configuration.Hue.MediaTypes, p.Metadata.Type) {
			h.handleHue(p)
		}

		if helper.Contains(h.configuration.LaMetric.MediaTypes, p.Metadata.Type) {
			h.handleLaMetric(p)
		}
	}

	for _, f := range h.configuration.Forwards {
		if f.Account == p.Account.ID {
			h.handleForward(f.Destination, raw)
		}
	}
}

func (h *EventHandler) handleHue(p plex.WebhookPayload) {
	if p.Event != "media.play" && p.Event != "media.resume" && p.Event != "media.pause" && p.Event != "media.stop" {
		return
	}

	cfg := h.configuration.Hue
	if cfg.IpAddress == "" || cfg.ApiKey == "" {
		return
	}

	client := &http.Client{}
	url := fmt.Sprintf("http://%s/api/%s/groups/%s", cfg.IpAddress, cfg.ApiKey, cfg.GroupId)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error while calling hue bridge: %v\n", err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	group := HueGroup{}
	err = json.Unmarshal(body, &group)
	if err != nil {
		log.Printf("Error while parsing api response: %v\n", err)
	}

	if group.State.AllOn {
		url += "/action"
		var data HueRequestData

		if p.Event == "media.play" || p.Event == "media.resume" {
			log.Printf(" ---> Activating play scene")
			data = HueRequestData{Scene: cfg.SceneIdPlay}
		} else if p.Event == "media.pause" || p.Event == "media.stop" {
			log.Printf(" ---> Activating pause scene")
			data = HueRequestData{Scene: cfg.SceneIdPause}
		}

		jsonValue, _ := json.Marshal(data)
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error while sending put request: %v\n", err)
		}

		defer res.Body.Close()
	} else {
		log.Printf(" ---> Skipping scene change\n")
	}
}

func (h *EventHandler) handleLaMetric(p plex.WebhookPayload) {
	if p.Event != "media.play" && p.Event != "media.resume" {
		log.Printf(" ---> Skip sending LaMetric notification\n")
		return
	}

	cfg := h.configuration.LaMetric
	if cfg.IpAddress == "" || cfg.ApiKey == "" {
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://%s:4343/api/v2/device/notifications", cfg.IpAddress)

	data := LaMetricRequestData{
		Priority: "info",
		IconType: "none",
		Model: LaMetricModel{
			Frames: []LaMetricFrame{
				{
					Text: p.Metadata.Title,
					Icon: 7788,
				},
			},
			Cycles: 2,
		},
	}
	jsonValue, _ := json.Marshal(data)

	auth := fmt.Sprintf("dev:%s", cfg.ApiKey)
	basic := base64.StdEncoding.EncodeToString([]byte(auth))

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basic))
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error while sending put request: %v\n", err)
	}

	defer res.Body.Close()
}

func (h *EventHandler) handleForward(url string, raw string) {
	log.Printf("Forwarding webhook to url: %s", url)
	client := &http.Client{}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("payload", raw)
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error while forwarding webhook: %v\n", err)
	} else {
		log.Printf("Finished forwarding webhook: %d", res.StatusCode)
	}

	defer res.Body.Close()
}
