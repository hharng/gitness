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

package types

type Pipeline struct {
	ID          int64  `db:"pipeline_id"              json:"id"`
	Description string `db:"pipeline_description"     json:"description"`
	UID         string `db:"pipeline_uid"             json:"uid"`
	Disabled    bool   `db:"pipeline_disabled"        json:"disabled"`
	CreatedBy   int64  `db:"pipeline_created_by"      json:"created_by"`
	// Seq is the last execution number for this pipeline
	Seq           int64  `db:"pipeline_seq"             json:"seq"`
	RepoID        int64  `db:"pipeline_repo_id"         json:"repo_id"`
	DefaultBranch string `db:"pipeline_default_branch"  json:"default_branch"`
	ConfigPath    string `db:"pipeline_config_path"     json:"config_path"`
	Created       int64  `db:"pipeline_created"         json:"created"`
	// Execution contains information about the latest execution if available
	Execution *Execution `db:"-"                        json:"execution,omitempty"`
	Updated   int64      `db:"pipeline_updated"         json:"updated"`
	Version   int64      `db:"pipeline_version"         json:"-"`
}
