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

package protection

import (
	"context"
	"reflect"
	"testing"

	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// nolint:gocognit // it's a unit test
func TestBranch_MergeVerify(t *testing.T) {
	user := &types.Principal{ID: 42}
	admin := &types.Principal{ID: 66, Admin: true}

	tests := []struct {
		name   string
		branch Branch
		in     MergeVerifyInput
		expOut MergeVerifyOutput
		expVs  []types.RuleViolations
	}{
		{
			name:   "empty",
			branch: Branch{},
			in:     MergeVerifyInput{Actor: user},
			expOut: MergeVerifyOutput{
				DeleteSourceBranch: false,
				AllowedMethods:     enum.MergeMethods,
			},
			expVs: []types.RuleViolations{},
		},
		{
			name: "admin-no-bypass",
			branch: Branch{
				Bypass: DefBypass{},
				PullReq: DefPullReq{
					StatusChecks: DefStatusChecks{RequireUIDs: []string{"abc"}},
					Comments:     DefComments{RequireResolveAll: true},
					Merge:        DefMerge{DeleteBranch: true},
				},
			},
			in: MergeVerifyInput{
				Actor:       admin,
				AllowBypass: false,
				PullReq:     &types.PullReq{UnresolvedCount: 1},
			},
			expOut: MergeVerifyOutput{
				DeleteSourceBranch: true,
				AllowedMethods:     enum.MergeMethods,
			},
			expVs: []types.RuleViolations{
				{
					Bypassable: true,
					Bypassed:   false,
					Violations: []types.Violation{
						{Code: codePullReqCommentsReqResolveAll},
						{Code: codePullReqStatusChecksReqUIDs},
					},
				},
			},
		},
		{
			name: "user-bypass",
			branch: Branch{
				Bypass: DefBypass{UserIDs: []int64{user.ID}},
				PullReq: DefPullReq{
					StatusChecks: DefStatusChecks{RequireUIDs: []string{"abc"}},
					Comments:     DefComments{RequireResolveAll: true},
					Merge:        DefMerge{DeleteBranch: true},
				},
			},
			in: MergeVerifyInput{
				Actor:       user,
				AllowBypass: true,
				PullReq:     &types.PullReq{UnresolvedCount: 1},
			},
			expOut: MergeVerifyOutput{
				DeleteSourceBranch: true,
				AllowedMethods:     enum.MergeMethods,
			},
			expVs: []types.RuleViolations{
				{
					Bypassable: true,
					Bypassed:   true,
					Violations: []types.Violation{
						{Code: codePullReqCommentsReqResolveAll},
						{Code: codePullReqStatusChecksReqUIDs},
					},
				},
			},
		},
		{
			name: "user-no-bypass",
			branch: Branch{
				PullReq: DefPullReq{
					StatusChecks: DefStatusChecks{RequireUIDs: []string{"abc"}},
					Comments:     DefComments{RequireResolveAll: true},
					Merge:        DefMerge{DeleteBranch: true},
				},
			},
			in: MergeVerifyInput{
				Actor:       user,
				AllowBypass: true,
				PullReq:     &types.PullReq{UnresolvedCount: 1},
			},
			expOut: MergeVerifyOutput{
				DeleteSourceBranch: true,
				AllowedMethods:     enum.MergeMethods,
			},
			expVs: []types.RuleViolations{
				{
					Bypassable: false,
					Bypassed:   false,
					Violations: []types.Violation{
						{Code: codePullReqCommentsReqResolveAll},
						{Code: codePullReqStatusChecksReqUIDs},
					},
				},
			},
		},
		{
			name: "merge-methods",
			branch: Branch{
				Bypass: DefBypass{},
				PullReq: DefPullReq{
					StatusChecks: DefStatusChecks{},
					Comments:     DefComments{},
					Merge: DefMerge{
						StrategiesAllowed: []enum.MergeMethod{enum.MergeMethodRebase, enum.MergeMethodSquash},
						DeleteBranch:      false,
					},
				},
			},
			in: MergeVerifyInput{
				Actor: user,
			},
			expOut: MergeVerifyOutput{
				DeleteSourceBranch: false,
				AllowedMethods:     []enum.MergeMethod{enum.MergeMethodRebase, enum.MergeMethodSquash},
			},
			expVs: []types.RuleViolations{},
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.branch.Sanitize(); err != nil {
				t.Errorf("invalid: %s", err.Error())
				return
			}

			out, results, err := test.branch.MergeVerify(ctx, test.in)
			if err != nil {
				t.Errorf("error: %s", err.Error())
				return
			}

			if want, got := test.expOut, out; !reflect.DeepEqual(want, got) {
				t.Errorf("output: want=%+v got=%+v", want, got)
			}

			if want, got := len(test.expVs), len(results); want != got {
				t.Errorf("number of violations mismatch: want=%d got=%d", want, got)
				return
			}

			for i := range results {
				if want, got := test.expVs[i].Bypassable, results[i].Bypassable; want != got {
					t.Errorf("rule result %d, bypassable mismatch: want=%t got=%t", i, want, got)
					return
				}

				if want, got := test.expVs[i].Bypassed, results[i].Bypassed; want != got {
					t.Errorf("rule result %d, bypassed mismatch: want=%t got=%t", i, want, got)
					return
				}

				if want, got := len(test.expVs[i].Violations), len(results[i].Violations); want != got {
					t.Errorf("rule result %d, violations count mismatch: want=%d got=%d", i, want, got)
					return
				}

				for j := range results[i].Violations {
					if want, got := test.expVs[i].Violations[j].Code, results[i].Violations[j].Code; want != got {
						t.Errorf("rule result %d, violation %d, code mismatch: want=%s got=%s", i, j, want, got)
					}
				}
			}
		})
	}
}

// nolint:gocognit // it's a unit test
func TestBranch_RefChangeVerify(t *testing.T) {
	user := &types.Principal{ID: 42}
	admin := &types.Principal{ID: 66, Admin: true}

	tests := []struct {
		name   string
		branch Branch
		in     RefChangeVerifyInput
		expVs  []types.RuleViolations
	}{
		{
			name: "empty",
			branch: Branch{
				Bypass:    DefBypass{},
				Lifecycle: DefLifecycle{},
			},
			in: RefChangeVerifyInput{
				Actor: user,
			},
			expVs: []types.RuleViolations{},
		},
		{
			name: "admin-no-bypass",
			branch: Branch{
				Bypass:    DefBypass{},
				Lifecycle: DefLifecycle{DeleteForbidden: true},
			},
			in: RefChangeVerifyInput{
				Actor:       admin,
				AllowBypass: false,
				RefAction:   RefActionDelete,
				RefType:     RefTypeBranch,
				RefNames:    []string{"abc"},
			},
			expVs: []types.RuleViolations{
				{
					Bypassable: true,
					Bypassed:   false,
					Violations: []types.Violation{
						{Code: codeLifecycleDelete},
					},
				},
			},
		},
		{
			name: "owner-bypass",
			branch: Branch{
				Bypass:    DefBypass{RepoOwners: true},
				Lifecycle: DefLifecycle{DeleteForbidden: true},
			},
			in: RefChangeVerifyInput{
				Actor:       user,
				AllowBypass: true,
				IsRepoOwner: true,
				RefAction:   RefActionDelete,
				RefType:     RefTypeBranch,
				RefNames:    []string{"abc"},
			},
			expVs: []types.RuleViolations{
				{
					Bypassable: true,
					Bypassed:   true,
					Violations: []types.Violation{
						{Code: codeLifecycleDelete},
					},
				},
			},
		},
		{
			name: "user-no-bypass",
			branch: Branch{
				Bypass:    DefBypass{RepoOwners: true},
				Lifecycle: DefLifecycle{DeleteForbidden: true},
			},
			in: RefChangeVerifyInput{
				Actor:       user,
				AllowBypass: true,
				IsRepoOwner: false,
				RefAction:   RefActionDelete,
				RefType:     RefTypeBranch,
				RefNames:    []string{"abc"},
			},
			expVs: []types.RuleViolations{
				{
					Bypassable: false,
					Bypassed:   false,
					Violations: []types.Violation{
						{Code: codeLifecycleDelete},
					},
				},
			},
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.branch.Sanitize(); err != nil {
				t.Errorf("invalid: %s", err.Error())
				return
			}

			results, err := test.branch.RefChangeVerify(ctx, test.in)
			if err != nil {
				t.Errorf("error: %s", err.Error())
				return
			}

			if want, got := len(test.expVs), len(results); want != got {
				t.Errorf("number of violations mismatch: want=%d got=%d", want, got)
				return
			}

			for i := range results {
				if want, got := test.expVs[i].Bypassable, results[i].Bypassable; want != got {
					t.Errorf("rule result %d, bypassable mismatch: want=%t got=%t", i, want, got)
					return
				}

				if want, got := test.expVs[i].Bypassed, results[i].Bypassed; want != got {
					t.Errorf("rule result %d, bypassed mismatch: want=%t got=%t", i, want, got)
					return
				}

				if want, got := len(test.expVs[i].Violations), len(results[i].Violations); want != got {
					t.Errorf("rule result %d, violations count mismatch: want=%d got=%d", i, want, got)
					return
				}

				for j := range results[i].Violations {
					if want, got := test.expVs[i].Violations[j].Code, results[i].Violations[j].Code; want != got {
						t.Errorf("rule result %d, violation %d, code mismatch: want=%s got=%s", i, j, want, got)
					}
				}
			}
		})
	}
}
