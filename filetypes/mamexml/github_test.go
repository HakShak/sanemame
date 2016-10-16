package mamexml

import "testing"

import "strings"

func TestGetReleases(t *testing.T) {
	releases, err := GetReleases()
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
