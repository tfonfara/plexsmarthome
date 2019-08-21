package plex

import (
	"encoding/json"
	"log"
	"net/http"
)

type WebhookPayload struct {
	Event   string `json:"event"`
	User    bool   `json:"user"`
	Owner   bool   `json:"owner"`
	Rating  string `json:"rating"`
	Account struct {
		ID    int    `json:"id"`
		Thumb string `json:"thumb"`
		Title string `json:"title"`
	} `json:"Account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"Server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"PublicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"Player"`
	Metadata struct {
		LibrarySectionType   string `json:"librarySectionType"`
		RatingKey            string `json:"ratingKey"`
		Key                  string `json:"key"`
		ParentRatingKey      string `json:"parentRatingKey"`
		GrandparentRatingKey string `json:"grandparentRatingKey"`
		GUID                 string `json:"guid"`
		LibrarySectionID     int    `json:"librarySectionID"`
		Type            	 string `json:"type"`
		Title                string `json:"title"`
		GrandparentKey       string `json:"grandparentKey"`
		ParentKey            string `json:"parentKey"`
		GrandparentTitle     string `json:"grandparentTitle"`
		ParentTitle          string `json:"parentTitle"`
		Summary              string `json:"summary"`
		Index                int    `json:"index"`
		ParentIndex          int    `json:"parentIndex"`
		RatingCount          int    `json:"ratingCount"`
		Thumb                string `json:"thumb"`
		Art                  string `json:"art"`
		ParentThumb          string `json:"parentThumb"`
		GrandparentThumb     string `json:"grandparentThumb"`
		GrandparentArt       string `json:"grandparentArt"`
		AddedAt              int    `json:"addedAt"`
		UpdatedAt            int    `json:"updatedAt"`
	} `json:"Metadata"`
}

type Webhook struct {
	onEvent func(p WebhookPayload, raw string)
}

func NewWebhook(fn func(p WebhookPayload, raw string)) *Webhook {
	return &Webhook{fn }
}

func (wh *Webhook) Handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(0); err != nil {
		log.Printf("Error while reading form: %v\n", err)
		return
	}

	payload, hasPayload := r.MultipartForm.Value["payload"]
	if hasPayload {
		var eventPayload WebhookPayload
		if err := json.Unmarshal([]byte(payload[0]), &eventPayload); err != nil {
			log.Printf("Error while parsing json: %v\n", err)
			return
		}
		wh.onEvent(eventPayload, payload[0])
	}
}