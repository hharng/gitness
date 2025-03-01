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

package limiter

import (
	"context"

	"github.com/harness/gitness/errors"
)

var ErrMaxNumReposReached = errors.New("maximum number of repositories reached")

// ResourceLimiter is an interface for managing resource limitation.
type ResourceLimiter interface {
	// RepoCount allows the creation of a specified number of repositories.
	RepoCount(ctx context.Context, count int) error
}

var _ ResourceLimiter = Unlimited{}

type Unlimited struct {
}

// NewResourceLimiter creates a new instance of ResourceLimiter.
func NewResourceLimiter() ResourceLimiter {
	return Unlimited{}
}

//nolint:revive
func (Unlimited) RepoCount(ctx context.Context, count int) error {
	return nil
}
