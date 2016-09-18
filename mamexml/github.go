package mamexml

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type Asset struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	Created     string `json:"created_at"`
	Updated     string `json:"updated_at"`
	Url         string `json:"browser_download_url"`
}

type Release struct {
	Name    string  `json:"name"`
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

func GetGithubReleases() ([]Release, error) {
	resp, err := http.Get("https://api.github.com/repos/mamedev/mame/releases")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	releases := make([]Release, 0)

	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func GetGithubDownloadUrl() (string, error) {
	releases, err := GetGithubReleases()
	if err != nil {
		return "", err
	}

	if len(releases) == 0 {
		return "", nil
	}

	for _, release := range releases {
		for _, asset := range release.Assets {
			matched, _ := regexp.MatchString("mame[0-9]*lx\\.zip", asset.Name)
			if matched {
				return asset.Url, nil
			}
		}
	}

	return "", nil
}
