// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package enum

import (
	"strings"
)

// Defines repo attributes that can be used for sorting and filtering.
type RepoAttr int

// Order enumeration.
const (
	RepoAttrNone RepoAttr = iota
	RepoAttrUID
	RepoAttrCreated
	RepoAttrUpdated
)

// ParseRepoAttr parses the repo attribute string
// and returns the equivalent enumeration.
func ParseRepoAttr(s string) RepoAttr {
	switch strings.ToLower(s) {
	case uid:
		return RepoAttrUID
	case created, createdAt:
		return RepoAttrCreated
	case updated, updatedAt:
		return RepoAttrUpdated
	default:
		return RepoAttrNone
	}
}

// String returns the string representation of the attribute.
func (a RepoAttr) String() string {
	switch a {
	case RepoAttrUID:
		return uid
	case RepoAttrCreated:
		return created
	case RepoAttrUpdated:
		return updated
	case RepoAttrNone:
		return ""
	default:
		return undefined
	}
}
