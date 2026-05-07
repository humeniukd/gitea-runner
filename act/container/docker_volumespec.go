// Copyright 2026 The Gitea Authors. All rights reserved.
// Copyright 2025 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !(WITHOUT_DOCKER || !(linux || darwin || windows || netbsd))

package container

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/moby/moby/api/types/mount"
)

const endOfVolumeSpec = rune(0)

// volumeSpec is a local copy of the small subset of docker/cli's volume parser
// that we need when translating `-v/--volume` flags into binds vs volumes.
type volumeSpec struct {
	Type     string
	Source   string
	Target   string
	ReadOnly bool
}

func parseVolumeSpec(spec string) (volumeSpec, error) {
	v := volumeSpec{}

	switch len(spec) {
	case 0:
		return v, errors.New("invalid empty volume spec")
	case 1, 2:
		v.Target = spec
		v.Type = string(mount.TypeVolume)
		return v, nil
	}

	buffer := []rune{}
	for _, char := range spec + string(endOfVolumeSpec) {
		switch {
		case isWindowsDrive(buffer, char):
			buffer = append(buffer, char)
		case char == ':' || char == endOfVolumeSpec:
			if err := populateVolumeField(char, buffer, &v); err != nil {
				populateVolumeType(&v)
				return v, fmt.Errorf("invalid spec: %s: %w", spec, err)
			}
			buffer = []rune{}
		default:
			buffer = append(buffer, char)
		}
	}

	populateVolumeType(&v)
	return v, nil
}

func isWindowsDrive(buffer []rune, char rune) bool {
	return char == ':' && len(buffer) == 1 && unicode.IsLetter(buffer[0])
}

func populateVolumeField(char rune, buffer []rune, v *volumeSpec) error {
	strBuffer := string(buffer)
	switch {
	case len(buffer) == 0:
		return errors.New("empty section between colons")
	case v.Source == "" && char == endOfVolumeSpec:
		v.Target = strBuffer
		return nil
	case v.Source == "":
		v.Source = strBuffer
		return nil
	case v.Target == "":
		v.Target = strBuffer
		return nil
	case char == ':':
		return errors.New("too many colons")
	}
	for option := range strings.SplitSeq(strBuffer, ",") {
		switch option {
		case "ro":
			v.ReadOnly = true
		case "rw":
			v.ReadOnly = false
		}
	}
	return nil
}

func populateVolumeType(v *volumeSpec) {
	switch {
	case v.Source == "":
		v.Type = string(mount.TypeVolume)
	case isFilePath(v.Source):
		v.Type = string(mount.TypeBind)
	default:
		v.Type = string(mount.TypeVolume)
	}
}

func isFilePath(source string) bool {
	switch source[0] {
	case '.', '/', '~':
		return true
	}
	if len([]rune(source)) == 1 {
		return false
	}
	if strings.HasPrefix(source, `\\`) {
		return true
	}

	first, nextIndex := utf8.DecodeRuneInString(source)
	return isWindowsDrive([]rune{first}, rune(source[nextIndex]))
}
