package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

type VersionType string

const (
	Main  VersionType = "main"
	Major VersionType = "major"
	Minor VersionType = "minor"
)

func (v *VersionType) String() string {
	return string(*v)
}

func (v *VersionType) Set(s string) error {
	switch s {
	case string(Main), string(Major), string(Minor):
		*v = VersionType(s)
		return nil
	default:
		return fmt.Errorf("invalid version type: %s. Available: main, major, or minor", s)
	}
}

type TagVersion struct {
	Main  int
	Major int
	Minor int
}

func (t *TagVersion) String() string {
	return fmt.Sprintf(
		"v%d.%d.%d",
		t.Main,
		t.Major,
		t.Minor,
	)
}

func (t *TagVersion) Bump(bumpType VersionType) {
	switch bumpType {
	case Main:
		t.Main++
		t.Major = 0
		t.Minor = 0
	case Major:
		t.Major++
		t.Minor = 0
	case Minor:
		t.Minor++
	}
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func NewTag(s string) *TagVersion {
	s = strings.TrimSpace(s)
	_p := strings.Split(s, "v")
	versions := strings.Split(_p[1], ".")
	return &TagVersion{
		Main:  mustAtoi(versions[0]),
		Major: mustAtoi(versions[1]),
		Minor: mustAtoi(versions[2]),
	}
}

func gitTag(tag TagVersion) {
	cmd := exec.Command(
		"git", "tag", tag.String(),
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"failed to git tag",
			"err", err,
		)
		panic(err)
	}
}

func gitTagsPush() {
	cmd := exec.Command(
		"git", "push", "origin", "--tags",
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"failed to git push tags",
			"err", err,
		)
		panic(err)
	}
}

func getTags() []TagVersion {
	var tags []TagVersion
	cmd := exec.Command(
		"git", "tag",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"getting tags failed",
			"err", err,
		)
		panic(err)
	}
	tagsRaw := string(out)
	for _, line := range strings.Split(tagsRaw, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		tags = append(
			tags,
			*NewTag(line),
		)
	}
	return tags
}

func getLatestTag(tags []TagVersion) TagVersion {
	var latest TagVersion
	highestMain := -1
	for _, t := range tags {
		if t.Main > highestMain {
			highestMain = t.Main
			latest = t
		}
	}
	highestMajor := -1
	for _, t := range tags {
		if highestMain > t.Main {
			continue
		}
		if t.Major > highestMajor {
			highestMajor = t.Major
			latest = t
		}
	}
	highestMinor := -1
	for _, t := range tags {
		if highestMain > t.Main || highestMajor > t.Major {
			continue
		}
		if t.Minor > highestMinor {
			highestMinor = t.Minor
			latest = t
		}
	}
	return latest
}

func main() {
	var version VersionType = Minor
	flag.Var(&version, "b", "What to bump: main, major, or minor (default: minor)")
	isDryRun := flag.Bool("dry", false, "If true, only logs will be shown, an actual git tag and push won't be executed")
	flag.Parse()
	if *isDryRun {
		slog.Info("dry run found, won't git tag nor git push")
	}
	tags := getTags()
	latest := getLatestTag(tags)
	slog.Debug(
		"git information",
		"found tags", len(tags),
		"latest tag", latest.String(),
	)
	latest.Bump(version)
	slog.Debug(
		"git tagging",
		"version", latest.String(),
	)
	if !*isDryRun {
		gitTag(latest)
	}
	slog.Debug(
		"git pushing",
	)
	if !*isDryRun {
		gitTagsPush()
	}
}
