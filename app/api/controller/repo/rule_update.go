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

package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/gitness/app/api/usererror"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/services/protection"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/check"
	"github.com/harness/gitness/types/enum"
)

type RuleUpdateInput struct {
	UID         string              `json:"uid"`
	State       *enum.RuleState     `json:"state"`
	Description *string             `json:"description"`
	Pattern     *protection.Pattern `json:"pattern"`
	Definition  *json.RawMessage    `json:"definition"`
}

// sanitize validates and sanitizes the update rule input data.
func (in *RuleUpdateInput) sanitize() error {
	if in.UID != "" {
		if err := check.UID(in.UID); err != nil {
			return err
		}
	}

	if in.State != nil {
		state, ok := in.State.Sanitize()
		if !ok {
			return usererror.BadRequest("rule state is invalid")
		}

		in.State = &state
	}

	if in.Pattern != nil {
		if err := in.Pattern.Validate(); err != nil {
			return usererror.BadRequestf("invalid pattern: %s", err)
		}
	}

	if in.Definition != nil && len(*in.Definition) == 0 {
		return usererror.BadRequest("rule definition missing")
	}

	return nil
}

func (in *RuleUpdateInput) isEmpty() bool {
	return in.UID == "" && in.State == nil && in.Description == nil && in.Pattern == nil && in.Definition == nil
}

// RuleUpdate updates an existing protection rule for a repository.
func (c *Controller) RuleUpdate(ctx context.Context,
	session *auth.Session,
	repoRef string,
	uid string,
	in *RuleUpdateInput,
) (*types.Rule, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoEdit, false)
	if err != nil {
		return nil, err
	}

	err = in.sanitize()
	if err != nil {
		return nil, err
	}

	r, err := c.ruleStore.FindByUID(ctx, nil, &repo.ID, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get a repository rule by its uid: %w", err)
	}

	if in.isEmpty() {
		r.Users, err = c.getRuleUsers(ctx, r)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

	if in.UID != "" {
		r.UID = in.UID
	}
	if in.State != nil {
		r.State = *in.State
	}
	if in.Description != nil {
		r.Description = *in.Description
	}
	if in.Pattern != nil {
		r.Pattern = in.Pattern.JSON()
	}
	if in.Definition != nil {
		r.Definition, err = c.protectionManager.SanitizeJSON(r.Type, *in.Definition)
		if err != nil {
			return nil, usererror.BadRequestf("invalid rule definition: %s", err.Error())
		}
	}

	r.Users, err = c.getRuleUsers(ctx, r)
	if err != nil {
		return nil, err
	}

	err = c.ruleStore.Update(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to update repository-level protection rule: %w", err)
	}

	return r, nil
}
