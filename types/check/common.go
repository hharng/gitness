// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package check

import (
	"fmt"
	"regexp"
)

const (
	minPathNameLength = 1
	maxPathNameLength = 64
	pathNameRegex     = "^[a-z][a-z0-9\\-\\_]*$"

	minNameLength = 1
	maxNameLength = 256
	nameRegex     = "^[a-zA-Z][a-zA-Z0-9\\-\\_ ]*$"

	minUIDLength = 2
	maxUIDLength = 64
	uidRegex     = "^[a-z][a-z0-9\\-\\_]*$"
)

var (
	ErrPathNameLength = &ValidationError{
		fmt.Sprintf("Path name has to be between %d and %d in length.", minPathNameLength, maxPathNameLength),
	}
	ErrPathNameRegex = &ValidationError{"Path name has to start with a letter and only contain the following [a-z0-9-_]."}

	ErrNameLength = &ValidationError{
		fmt.Sprintf("Name has to be between %d and %d in length.",
			minNameLength, maxNameLength),
	}
	ErrNameRegex = &ValidationError{
		"Name has to start with a letter and only contain the following [a-zA-Z0-9-_ ].",
	}

	ErrUIDLength = &ValidationError{
		fmt.Sprintf("UID has to be between %d and %d in length.",
			minUIDLength, maxUIDLength),
	}
	ErrUIDRegex = &ValidationError{
		"UID has to start with a letter and only contain the following [a-z0-9-_].",
	}
)

// PathName checks the provided name and returns an error in it isn't valid.
func PathName(pathName string) error {
	l := len(pathName)
	if l < minPathNameLength || l > maxPathNameLength {
		return ErrPathNameLength
	}

	if ok, _ := regexp.Match(pathNameRegex, []byte(pathName)); !ok {
		return ErrPathNameRegex
	}

	return nil
}

// Name checks the provided name and returns an error in it isn't valid.
func Name(name string) error {
	l := len(name)
	if l < minNameLength || l > maxNameLength {
		return ErrNameLength
	}

	if ok, _ := regexp.Match(nameRegex, []byte(name)); !ok {
		return ErrNameRegex
	}

	return nil
}

// UID checks the provided uid and returns an error in it isn't valid.
func UID(uid string) error {
	l := len(uid)
	if l < minUIDLength || l > maxUIDLength {
		return ErrUIDLength
	}

	if ok, _ := regexp.Match(uidRegex, []byte(uid)); !ok {
		return ErrUIDRegex
	}

	return nil
}