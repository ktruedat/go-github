// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_GetAllRepositoryRulesets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"id": 21,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/21"
			  }
			}
		}]`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetAllRepositoryRulesets(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetAllRepositoryRulesets returned error: %v", err)
	}

	want := []*RepositoryRuleset{{
		ID:          Ptr(int64(21)),
		Name:        "test ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RepositoryRulesetLinks{
			Self: &RepositoryRulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/21")},
		},
	}}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetAllRepositoryRulesets returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetAllRepositoryRulesets"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetAllRepositoryRulesets(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateRepositoryRuleset_RepoNames(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_actors": [
			  {
				"actor_id": 234,
				"actor_type": "Team"
			  }
			],
			"conditions": {
			  "ref_name": {
				"include": [
				  "refs/heads/main",
				  "refs/heads/master"
				],
				"exclude": [
				  "refs/heads/dev*"
				]
			  },
			  "repository_name": {
				"include": [
				  "important_repository",
				  "another_important_repository"
				],
				"exclude": [
				  "unimportant_repository"
				],
				"protected": true
			  }
			},
			"rules": [
			  {
				"type": "creation"
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				"allowed_merge_methods": ["rebase","squash"],
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
					"do_not_enforce_on_create": true,
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Organizations.CreateRepositoryRuleset(ctx, "o", RepositoryRuleset{
		Name:        "ruleset",
		Target:      Ptr(RulesetTargetBranch),
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr(BypassActorTypeTeam),
			},
		},
		Conditions: &RepositoryRulesetConditions{
			RefName: &RepositoryRulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RepositoryRulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: &RepositoryRulesetRules{
			Creation: &EmptyRuleParameters{},
			Update: &UpdateRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			},
			Deletion:              &EmptyRuleParameters{},
			RequiredLinearHistory: &EmptyRuleParameters{},
			RequiredDeployments: &RequiredDeploymentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			},
			RequiredSignatures: &EmptyRuleParameters{},
			PullRequest: &PullRequestRuleParameters{
				AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase, PullRequestMergeMethodSquash},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			},
			RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
				DoNotEnforceOnCreate: Ptr(true),
				RequiredStatusChecks: []*RuleStatusCheck{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			},
			NonFastForward: &EmptyRuleParameters{},
			CommitMessagePattern: &PatternRuleParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			},
			CommitAuthorEmailPattern: &PatternRuleParameters{
				Operator: "contains",
				Pattern:  "github",
			},
			CommitterEmailPattern: &PatternRuleParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			},
			BranchNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			},
			TagNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			},
			CodeScanning: &CodeScanningRuleParameters{
				CodeScanningTools: []*RuleCodeScanningTool{
					{
						AlertsThreshold:         CodeScanningAlertsThresholdErrors,
						SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						Tool:                    "CodeQL",
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateRepositoryRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr(BypassActorTypeTeam),
			},
		},
		Conditions: &RepositoryRulesetConditions{
			RefName: &RepositoryRulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RepositoryRulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: &RepositoryRulesetRules{
			Creation: &EmptyRuleParameters{},
			Update: &UpdateRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			},
			Deletion:              &EmptyRuleParameters{},
			RequiredLinearHistory: &EmptyRuleParameters{},
			RequiredDeployments: &RequiredDeploymentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			},
			RequiredSignatures: &EmptyRuleParameters{},
			PullRequest: &PullRequestRuleParameters{
				AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase, PullRequestMergeMethodSquash},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			},
			RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
				DoNotEnforceOnCreate: Ptr(true),
				RequiredStatusChecks: []*RuleStatusCheck{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			},
			NonFastForward: &EmptyRuleParameters{},
			CommitMessagePattern: &PatternRuleParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			},
			CommitAuthorEmailPattern: &PatternRuleParameters{
				Operator: "contains",
				Pattern:  "github",
			},
			CommitterEmailPattern: &PatternRuleParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			},
			BranchNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			},
			TagNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			},
			CodeScanning: &CodeScanningRuleParameters{
				CodeScanningTools: []*RuleCodeScanningTool{
					{
						AlertsThreshold:         CodeScanningAlertsThresholdErrors,
						SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						Tool:                    "CodeQL",
					},
				},
			},
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Organizations.CreateRepositoryRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateRepositoryRuleset(ctx, "o", RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateRepositoryRuleset_RepoProperty(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_actors": [
			  {
				"actor_id": 234,
				"actor_type": "Team"
			  }
			],
			"conditions": {
			  "repository_property": {
				"include": [
					{
						"name": "testIncludeProp",
						"source": "custom",
						"property_values": [
							"true"
						]
					}
				],
				"exclude": [
					{
						"name": "testExcludeProp",
						"property_values": [
							"false"
						]
					}
				]
			  }
			},
			"rules": [
			  {
				"type": "creation"
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				"allowed_merge_methods": ["rebase","squash"],
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
					"do_not_enforce_on_create": true,
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Organizations.CreateRepositoryRuleset(ctx, "o", RepositoryRuleset{
		Name:        "ruleset",
		Target:      Ptr(RulesetTargetBranch),
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr(BypassActorTypeTeam),
			},
		},
		Conditions: &RepositoryRulesetConditions{
			RepositoryProperty: &RepositoryRulesetRepositoryPropertyConditionParameters{
				Include: []*RepositoryRulesetRepositoryPropertyTargetParameters{
					{
						Name:           "testIncludeProp",
						Source:         Ptr("custom"),
						PropertyValues: []string{"true"},
					},
				},
				Exclude: []*RepositoryRulesetRepositoryPropertyTargetParameters{
					{
						Name:           "testExcludeProp",
						PropertyValues: []string{"false"},
					},
				},
			},
		},
		Rules: &RepositoryRulesetRules{
			Creation: &EmptyRuleParameters{},
			Update: &UpdateRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			},
			Deletion:              &EmptyRuleParameters{},
			RequiredLinearHistory: &EmptyRuleParameters{},
			RequiredDeployments: &RequiredDeploymentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			},
			RequiredSignatures: &EmptyRuleParameters{},
			PullRequest: &PullRequestRuleParameters{
				AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase, PullRequestMergeMethodSquash},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			},
			RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
				DoNotEnforceOnCreate: Ptr(true),
				RequiredStatusChecks: []*RuleStatusCheck{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			},
			NonFastForward: &EmptyRuleParameters{},
			CommitMessagePattern: &PatternRuleParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			},
			CommitAuthorEmailPattern: &PatternRuleParameters{
				Operator: "contains",
				Pattern:  "github",
			},
			CommitterEmailPattern: &PatternRuleParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			},
			BranchNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			},
			TagNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			},
			CodeScanning: &CodeScanningRuleParameters{
				CodeScanningTools: []*RuleCodeScanningTool{
					{
						AlertsThreshold:         CodeScanningAlertsThresholdErrors,
						SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						Tool:                    "CodeQL",
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateRepositoryRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr(BypassActorTypeTeam),
			},
		},
		Conditions: &RepositoryRulesetConditions{
			RepositoryProperty: &RepositoryRulesetRepositoryPropertyConditionParameters{
				Include: []*RepositoryRulesetRepositoryPropertyTargetParameters{
					{
						Name:           "testIncludeProp",
						Source:         Ptr("custom"),
						PropertyValues: []string{"true"},
					},
				},
				Exclude: []*RepositoryRulesetRepositoryPropertyTargetParameters{
					{
						Name:           "testExcludeProp",
						PropertyValues: []string{"false"},
					},
				},
			},
		},
		Rules: &RepositoryRulesetRules{
			Creation: &EmptyRuleParameters{},
			Update: &UpdateRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			},
			Deletion:              &EmptyRuleParameters{},
			RequiredLinearHistory: &EmptyRuleParameters{},
			RequiredDeployments: &RequiredDeploymentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			},
			RequiredSignatures: &EmptyRuleParameters{},
			PullRequest: &PullRequestRuleParameters{
				AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase, PullRequestMergeMethodSquash},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			},
			RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
				DoNotEnforceOnCreate: Ptr(true),
				RequiredStatusChecks: []*RuleStatusCheck{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			},
			NonFastForward: &EmptyRuleParameters{},
			CommitMessagePattern: &PatternRuleParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			},
			CommitAuthorEmailPattern: &PatternRuleParameters{
				Operator: "contains",
				Pattern:  "github",
			},
			CommitterEmailPattern: &PatternRuleParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			},
			BranchNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			},
			TagNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			},
			CodeScanning: &CodeScanningRuleParameters{
				CodeScanningTools: []*RuleCodeScanningTool{
					{
						AlertsThreshold:         CodeScanningAlertsThresholdErrors,
						SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						Tool:                    "CodeQL",
					},
				},
			},
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Organizations.CreateRepositoryRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateRepositoryRuleset(ctx, "o", RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateRepositoryRuleset_RepoIDs(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_actors": [
			  {
				"actor_id": 234,
				"actor_type": "Team"
			  }
			],
			"conditions": {
			  "ref_name": {
				"include": [
				  "refs/heads/main",
				  "refs/heads/master"
				],
				"exclude": [
				  "refs/heads/dev*"
				]
			  },
			  "repository_id": {
					"repository_ids": [ 123, 456 ]
				}
			},
			"rules": [
			  {
				"type": "creation"
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
					"allowed_merge_methods": ["rebase","squash"],
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
					"do_not_enforce_on_create": true,
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Organizations.CreateRepositoryRuleset(ctx, "o", RepositoryRuleset{
		Name:        "ruleset",
		Target:      Ptr(RulesetTargetBranch),
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr(BypassActorTypeTeam),
			},
		},
		Conditions: &RepositoryRulesetConditions{
			RefName: &RepositoryRulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryID: &RepositoryRulesetRepositoryIDsConditionParameters{
				RepositoryIDs: []int64{123, 456},
			},
		},
		Rules: &RepositoryRulesetRules{
			Creation: &EmptyRuleParameters{},
			Update: &UpdateRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			},
			Deletion:              &EmptyRuleParameters{},
			RequiredLinearHistory: &EmptyRuleParameters{},
			RequiredDeployments: &RequiredDeploymentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			},
			RequiredSignatures: &EmptyRuleParameters{},
			PullRequest: &PullRequestRuleParameters{
				AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase, PullRequestMergeMethodSquash},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			},
			RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
				DoNotEnforceOnCreate: Ptr(true),
				RequiredStatusChecks: []*RuleStatusCheck{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			},
			NonFastForward: &EmptyRuleParameters{},
			CommitMessagePattern: &PatternRuleParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			},
			CommitAuthorEmailPattern: &PatternRuleParameters{
				Operator: "contains",
				Pattern:  "github",
			},
			CommitterEmailPattern: &PatternRuleParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			},
			BranchNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			},
			TagNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			},
			CodeScanning: &CodeScanningRuleParameters{
				CodeScanningTools: []*RuleCodeScanningTool{
					{
						AlertsThreshold:         CodeScanningAlertsThresholdErrors,
						SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						Tool:                    "CodeQL",
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateRepositoryRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr(BypassActorTypeTeam),
			},
		},
		Conditions: &RepositoryRulesetConditions{
			RefName: &RepositoryRulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryID: &RepositoryRulesetRepositoryIDsConditionParameters{
				RepositoryIDs: []int64{123, 456},
			},
		},
		Rules: &RepositoryRulesetRules{
			Creation: &EmptyRuleParameters{},
			Update: &UpdateRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			},
			Deletion:              &EmptyRuleParameters{},
			RequiredLinearHistory: &EmptyRuleParameters{},
			RequiredDeployments: &RequiredDeploymentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			},
			RequiredSignatures: &EmptyRuleParameters{},
			PullRequest: &PullRequestRuleParameters{
				AllowedMergeMethods:            []PullRequestMergeMethod{PullRequestMergeMethodRebase, PullRequestMergeMethodSquash},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			},
			RequiredStatusChecks: &RequiredStatusChecksRuleParameters{
				DoNotEnforceOnCreate: Ptr(true),
				RequiredStatusChecks: []*RuleStatusCheck{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			},
			NonFastForward: &EmptyRuleParameters{},
			CommitMessagePattern: &PatternRuleParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			},
			CommitAuthorEmailPattern: &PatternRuleParameters{
				Operator: "contains",
				Pattern:  "github",
			},
			CommitterEmailPattern: &PatternRuleParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			},
			BranchNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			},
			TagNamePattern: &PatternRuleParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			},
			CodeScanning: &CodeScanningRuleParameters{
				CodeScanningTools: []*RuleCodeScanningTool{
					{
						AlertsThreshold:         CodeScanningAlertsThresholdErrors,
						SecurityAlertsThreshold: CodeScanningSecurityAlertsThresholdHighOrHigher,
						Tool:                    "CodeQL",
					},
				},
			},
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Organizations.CreateRepositoryRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateRepositoryRuleset(ctx, "o", RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetRepositoryRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/21"
			  }
			},
			"conditions": {
				"ref_name": {
				  "include": [
					"refs/heads/main",
					"refs/heads/master"
				  ],
				  "exclude": [
					"refs/heads/dev*"
				  ]
				},
				"repository_name": {
				  "include": [
					"important_repository",
					"another_important_repository"
				  ],
				  "exclude": [
					"unimportant_repository"
				  ],
				  "protected": true
				}
			  },
			  "rules": [
				{
				  "type": "creation"
				}
			  ]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetRepositoryRuleset(ctx, "o", 21)
	if err != nil {
		t.Errorf("Organizations.GetOrganizationRepositoryRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(21)),
		Name:        "test ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RepositoryRulesetLinks{
			Self: &RepositoryRulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/21")},
		},
		Conditions: &RepositoryRulesetConditions{
			RefName: &RepositoryRulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RepositoryRulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: &RepositoryRulesetRules{Creation: &EmptyRuleParameters{}},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetRepositoryRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetRepositoryRuleset(ctx, "o", 21)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetRepositoryRulesetWithRepoPropCondition(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/21"
			  }
			},
			"conditions": {
			  "repository_property": {
				"exclude": [],
				"include": [
				  {
					"name": "testIncludeProp",
					"source": "custom",
					"property_values": [
					  "true"
					]
				  }
				]
			  }
			},
			"rules": [
			{
				"type": "creation"
			}
			]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetRepositoryRuleset(ctx, "o", 21)
	if err != nil {
		t.Errorf("Organizations.GetOrganizationRepositoryRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(21)),
		Name:        "test ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RepositoryRulesetLinks{
			Self: &RepositoryRulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/21")},
		},
		Conditions: &RepositoryRulesetConditions{
			RepositoryProperty: &RepositoryRulesetRepositoryPropertyConditionParameters{
				Include: []*RepositoryRulesetRepositoryPropertyTargetParameters{
					{
						Name:           "testIncludeProp",
						Source:         Ptr("custom"),
						PropertyValues: []string{"true"},
					},
				},
				Exclude: []*RepositoryRulesetRepositoryPropertyTargetParameters{},
			},
		},
		Rules: &RepositoryRulesetRules{Creation: &EmptyRuleParameters{}},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetRepositoryRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetRepositoryRuleset(ctx, "o", 21)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateRepositoryRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/21"
			  }
			},
			"conditions": {
				"ref_name": {
				  "include": [
					"refs/heads/main",
					"refs/heads/master"
				  ],
				  "exclude": [
					"refs/heads/dev*"
				  ]
				},
				"repository_name": {
				  "include": [
					"important_repository",
					"another_important_repository"
				  ],
				  "exclude": [
					"unimportant_repository"
				  ],
				  "protected": true
				}
			  },
			  "rules": [
				{
				  "type": "creation"
				}
			  ]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.UpdateRepositoryRuleset(ctx, "o", 21, RepositoryRuleset{
		Name:        "test ruleset",
		Target:      Ptr(RulesetTargetBranch),
		Enforcement: "active",
		Conditions: &RepositoryRulesetConditions{
			RefName: &RepositoryRulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RepositoryRulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: &RepositoryRulesetRules{Creation: &EmptyRuleParameters{}},
	})
	if err != nil {
		t.Errorf("Organizations.UpdateRepositoryRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(21)),
		Name:        "test ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RepositoryRulesetLinks{
			Self: &RepositoryRulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/21")},
		},
		Conditions: &RepositoryRulesetConditions{
			RefName: &RepositoryRulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RepositoryRulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: &RepositoryRulesetRules{Creation: &EmptyRuleParameters{}},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.UpdateRepositoryRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "UpdateRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateRepositoryRuleset(ctx, "o", 21, RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateRepositoryRulesetWithRepoProp(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/21"
			  }
			},
			"conditions": {
			  "repository_property": {
				"exclude": [],
				"include": [
				  {
					"name": "testIncludeProp",
					"source": "custom",
					"property_values": [
					  "true"
					]
				  }
				]
			  }
			  },
			  "rules": [
				{
				  "type": "creation"
				}
			  ]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.UpdateRepositoryRuleset(ctx, "o", 21, RepositoryRuleset{
		Name:        "test ruleset",
		Target:      Ptr(RulesetTargetBranch),
		Enforcement: "active",
		Conditions: &RepositoryRulesetConditions{
			RepositoryProperty: &RepositoryRulesetRepositoryPropertyConditionParameters{
				Include: []*RepositoryRulesetRepositoryPropertyTargetParameters{
					{
						Name:           "testIncludeProp",
						Source:         Ptr("custom"),
						PropertyValues: []string{"true"},
					},
				},
				Exclude: []*RepositoryRulesetRepositoryPropertyTargetParameters{},
			},
		},
		Rules: &RepositoryRulesetRules{Creation: &EmptyRuleParameters{}},
	})
	if err != nil {
		t.Errorf("Organizations.UpdateRepositoryRuleset returned error: %v", err)
	}

	want := &RepositoryRuleset{
		ID:          Ptr(int64(21)),
		Name:        "test ruleset",
		Target:      Ptr(RulesetTargetBranch),
		SourceType:  Ptr(RulesetSourceTypeOrganization),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RepositoryRulesetLinks{
			Self: &RepositoryRulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/21")},
		},
		Conditions: &RepositoryRulesetConditions{
			RepositoryProperty: &RepositoryRulesetRepositoryPropertyConditionParameters{
				Include: []*RepositoryRulesetRepositoryPropertyTargetParameters{
					{
						Name:           "testIncludeProp",
						Source:         Ptr("custom"),
						PropertyValues: []string{"true"},
					},
				},
				Exclude: []*RepositoryRulesetRepositoryPropertyTargetParameters{},
			},
		},
		Rules: &RepositoryRulesetRules{Creation: &EmptyRuleParameters{}},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.UpdateRepositoryRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "UpdateRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateRepositoryRuleset(ctx, "o", 21, RepositoryRuleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateRepositoryRulesetClearBypassActor(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"conditions": {
				"repository_name": {
					"include": [
						"important_repository",
						"another_important_repository"
					],
					"exclude": [
						"unimportant_repository"
					],
					"protected": true
				},
			  "ref_name": {
					"include": [
						"refs/heads/main",
						"refs/heads/master"
					],
					"exclude": [
						"refs/heads/dev*"
					]
				}
			},
			"rules": [
			  {
					"type": "creation"
			  }
			]
		}`)
	})

	ctx := context.Background()

	_, err := client.Organizations.UpdateRepositoryRulesetClearBypassActor(ctx, "o", 21)
	if err != nil {
		t.Errorf("Organizations.UpdateRepositoryRulesetClearBypassActor returned error: %v \n", err)
	}

	const methodName = "UpdateRepositoryRulesetClearBypassActor"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.UpdateRepositoryRulesetClearBypassActor(ctx, "o", 21)
	})
}

func TestOrganizationsService_DeleteRepositoryRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.DeleteRepositoryRuleset(ctx, "o", 21)
	if err != nil {
		t.Errorf("Organizations.DeleteRepositoryRuleset returned error: %v", err)
	}

	const methodName = "DeleteRepositoryRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.DeleteRepositoryRuleset(ctx, "0", 21)
	})
}
