// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package check

import (
	"fmt"
	"strings"

	"github.com/harness/gitness/types"
)

var (
	illegalRootSpaceNames = []string{"api"}

	ErrRootSpaceNameNotAllowed = &ValidationError{
		fmt.Sprintf("The following names are not allowed for a root space: %v", illegalRootSpaceNames),
	}
	ErrInvalidParentSpaceID = &ValidationError{
		"Parent space ID has to be either zero for a root space or greater than zero for a child space.",
	}
)

// Space checks the provided space and returns an error in it isn't valid.
func Space(space *types.Space) error {
	// validate name
	if err := PathName(space.PathName); err != nil {
		return err
	}

	// validate display name
	if err := Name(space.Name); err != nil {
		return err
	}

	if space.ParentID < 0 {
		return ErrInvalidParentSpaceID
	}

	// root space specific validations
	if space.ParentID == 0 {
		for _, p := range illegalRootSpaceNames {
			if strings.HasPrefix(space.PathName, p) {
				return ErrRootSpaceNameNotAllowed
			}
		}
	}

	return nil
}