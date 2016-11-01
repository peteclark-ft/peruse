package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/peteclark-ft/peruse/structs"
)

type list struct {
	Title string     `json:"title"`
	Items []listItem `json:"items"`
}

type listItem struct {
	URL string `json:"apiUrl"`
}

type getList struct {
	client   *http.Client
	url      string
	apiKey   string
	user     string
	password string
}

const article = "http://www.ft.com/ontology/content/Article"

func (g getList) requestList(uuid string) (list, []structs.UPPContent, error) {
	var content []structs.UPPContent
	jsonList := list{}

	req, err := http.NewRequest("GET", g.url+"/lists/"+uuid, nil)
	if err != nil {
		return jsonList, content, err
	}

	req.SetBasicAuth(g.user, g.password)

	resp, err := g.client.Do(req)
	if err != nil {
		return jsonList, content, err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&jsonList)
	if err != nil {
		return jsonList, content, err
	}

	for _, item := range jsonList.Items {
		result, err := g.requestContent(item)
		if err != nil {
			//logrus.WithError(err).Error("Failed to request content!")
			continue
		}

		content = append(content, result)
	}

	return jsonList, content, nil
}

func (g getList) requestContent(item listItem) (structs.UPPContent, error) {
	var uppContent structs.UPPContent

	uri, _ := url.Parse(item.URL)

	query := uri.Query()
	query.Add("apiKey", g.apiKey)
	uri.RawQuery = query.Encode()

	resp, err := http.Get(uri.String())

	if err != nil {
		return uppContent, err
	}

	if resp.StatusCode != 200 {
		return uppContent, errors.New("Non-200 status, skipping.")
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&uppContent)

	if err != nil {
		return uppContent, err
	}

	if uppContent.Type != article {
		return uppContent, errors.New("Skipping non-article")
	}

	return uppContent, nil
}
