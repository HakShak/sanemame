package mamexml

import "testing"

import "strings"

func TestGetGithubReleases(t *testing.T) {
	releases, err := GetGithubReleases()
	if err != nil {
		t.FailNow()
	}

	if len(releases) == 0 {
		t.FailNow()
	}

	for _, release := range releases {
		if !strings.Contains(release.TagName, "mame") {
			t.FailNow()
		}
	}

	if len(releases[0].Assets) == 0 {
		t.FailNow()
	}
}

func TestGetGithubDownloadUrl(t *testing.T) {
	url, err := GetGithubDownloadUrl()
	if err != nil {
		t.FailNow()
	}

	if url == "" {
		t.FailNow()
	}
}
