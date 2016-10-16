package mamexml

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

import "gopkg.in/cheggaaa/pb.v1"

//Asset artifacts from a release
type Asset struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	Created     string `json:"created_at"`
	Updated     string `json:"updated_at"`
	URL         string `json:"browser_download_url"`
}

//Release top level object for releases
type Release struct {
	Name    string  `json:"name"`
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

//GetReleases returns first page of releases
func GetReleases(api string, repo string) ([]Release, error) {
	startTime := time.Now()
	url := fmt.Sprintf(api, repo)
	log.Printf("Getting releases from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var releases []Release

	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}

	log.Printf("Got %d releases in %s", len(releases), time.Since(startTime))

	return releases, nil
}

//GetLatestRelease returns latest release
func GetLatestRelease(api string, repo string) (*Release, error) {
	releases, err := GetReleases(api, repo)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, nil
	}

	log.Printf("Latest release: %s", releases[0].TagName)

	return &releases[0], nil
}

//GetLatestXMLAsset returns info for most recent artifact
func GetLatestXMLAsset(api string, repo string) (*Asset, error) {
	release, err := GetLatestRelease(api, repo)
	if err != nil {
		return nil, err
	}

	for _, asset := range release.Assets {
		matched, _ := regexp.MatchString("mame[0-9]*lx\\.zip", asset.Name)
		if matched {
			log.Printf("Latest XML Asset: %s", asset.Name)
			return &asset, nil
		}
	}

	return nil, nil
}

//Download downloads url to filename
func Download(fileName string, url string, size int) error {
	log.Printf("Creating %s", fileName)
	out, err := os.Create(fileName)
	if err != nil {
		return nil
	}
	defer out.Close()

	log.Printf("Downloading %s", url)
	bar := pb.New(size).SetUnits(pb.U_BYTES).SetWidth(80).Start()
	bar.ShowSpeed = true
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	proxyReader := bar.NewProxyReader(response.Body)

	_, err = io.Copy(out, proxyReader)
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}

//ExtractAsset extracts compressed assest
func ExtractAsset(zipFileName string) (string, error) {
	zipReader, err := zip.OpenReader(zipFileName)
	if err != nil {
		return "", err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		log.Printf("Extracting file: %s", file.Name)
		bar := pb.New(int(file.UncompressedSize64)).SetUnits(pb.U_BYTES).SetWidth(80).Start()
		bar.ShowSpeed = true
		contentReader, err := file.Open()
		if err != nil {
			return "", err
		}
		defer contentReader.Close()

		proxyReader := bar.NewProxyReader(contentReader)

		extractedFile, err := os.Create(file.Name)
		if err != nil {
			return "", err
		}
		defer extractedFile.Close()

		_, err = io.Copy(extractedFile, proxyReader)
		if err != nil {
			return "", err
		}

		bar.Finish()

		return file.Name, nil
	}

	return "", nil
}

//GetLatestXMLFile finds, downloads, extracts latest release asset
func GetLatestXMLFile(api string, repo string) (string, error) {
	asset, err := GetLatestXMLAsset(api, repo)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(asset.Name); os.IsNotExist(err) {
		err := Download(asset.Name, asset.URL, asset.Size)
		if err != nil {
			return "", err
		}
	}

	fileInfo, err := os.Stat(asset.Name)
	if err != nil {
		return "", err
	}
	if fileInfo.Size() != int64(asset.Size) {
		return "", errors.New("File sizes do not match")
	}

	xmlFile, err := ExtractAsset(asset.Name)
	if err != nil {
		return "", err
	}

	return xmlFile, nil
}
