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

package webhook

import (
	"context"
	"fmt"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// ListExecutions returns the executions of the webhook.
func (c *Controller) ListExecutions(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	webhookUID string,
	filter *types.WebhookExecutionFilter,
) ([]*types.WebhookExecution, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoView)
	if err != nil {
		return nil, err
	}

	// get the webhook and ensure it belongs to us
	webhook, err := c.getWebhookVerifyOwnership(ctx, repo.ID, webhookUID)
	if err != nil {
		return nil, err
	}

	// get webhook executions
	webhookExecutions, err := c.webhookExecutionStore.ListForWebhook(ctx, webhook.ID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook executions for webhook %d: %w", webhook.ID, err)
	}

	return webhookExecutions, nil
}
