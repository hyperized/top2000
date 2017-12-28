package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

// https://github.com/ChimeraCoder/gojson
// curl -s http://radiobox2.omroep.nl/data/radiobox2/nowonair/2.json\?npo_cc_skip_wall\=1 | gojson -name=OnAir

const (
	aapjeURL string = "http://radiobox2.omroep.nl/data/radiobox2/nowonair/2.json?npo_cc_skip_wall=1"
)

type OnAir struct {
	Results []struct {
		Channel  int64  `json:"channel"`
		Date     string `json:"date"`
		ID       int64  `json:"id"`
		Songfile struct {
			References struct {
				Channel string `json:"channel"`
			} `json:"_references"`
			ReferencesSsl struct {
				Channel string `json:"channel"`
			} `json:"_references_ssl"`
			Artist      string `json:"artist"`
			BumaID      string `json:"buma_id"`
			DaletID     string `json:"dalet_id"`
			Hidden      int64  `json:"hidden"`
			ID          int64  `json:"id"`
			LastUpdated string `json:"last_updated"`
			Rb1id       int64  `json:"rb1id"`
			SongID      int64  `json:"song_id"`
			Songversion struct {
				ID    int64 `json:"id"`
				Image []struct {
					AllowedToUse   int64  `json:"allowed_to_use"`
					Created        string `json:"created"`
					Deleted        int64  `json:"deleted"`
					Filename       string `json:"filename"`
					Hash           string `json:"hash"`
					ID             int64  `json:"id"`
					Name           string `json:"name"`
					OriginalHeight int64  `json:"original_height"`
					OriginalWidth  int64  `json:"original_width"`
					Replaced       int64  `json:"replaced"`
					Source         string `json:"source"`
					Updated        string `json:"updated"`
					URL            string `json:"url"`
					URLSsl         string `json:"url_ssl"`
				} `json:"image"`
			} `json:"songversion"`
			Title string `json:"title"`
		} `json:"songfile"`
		Startdatetime string `json:"startdatetime"`
		Stopdatetime  string `json:"stopdatetime"`
	} `json:"results"`
}

func getAapjeData() (OnAir, error) {
	req, err := http.Get(aapjeURL)

	/*if err != nil {
		return nil, err
	}*/

	body, err := ioutil.ReadAll(req.Body)
	/*if err != nil {
		return nil, err
	}*/

	var dat OnAir
	err = json.Unmarshal(body, &dat)

	return dat, err
}

func getCurrentTrack(air OnAir) string {
	var buffer bytes.Buffer

	buffer.WriteString("Nu op Top2000: " + air.Results[0].Songfile.Title + " -  " + air.Results[0].Songfile.Artist)

	return buffer.String()
}