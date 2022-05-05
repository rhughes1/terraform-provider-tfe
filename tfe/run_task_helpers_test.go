package tfe

import (
	"testing"

	"github.com/hashicorp/go-tfe"
)

func TestFetchOrganizationRunTask(t *testing.T) {
	tests := map[string]struct {
		taskName     string
		org          string
		expectExists bool
		err          bool
	}{
		"non exisiting organization": {
			"a-task",
			"not-an-org",
			false,
			true,
		},
		"non exisiting task": {
			"not-a-task",
			"hashicorp",
			false,
			true,
		},
		"existing task": {
			"a-task",
			"hashicorp",
			true,
			false,
		},
	}

	client := testTfeClient(t, testClientOptions{defaultRunTaskID: "task-123"})
	client.RunTasks.Create(nil, "hashicorp", tfe.RunTaskCreateOptions{
		Name: "a-task",
		URL:  runTasksUrl(),
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := fetchOrganizationRunTask(test.taskName, test.org, client)

			if (err != nil) != test.err {
				t.Fatalf("expected error is %t, got %v", test.err, err)
			}

			if test.expectExists {
				if got == nil || got.Name != test.taskName {
					t.Fatalf("wrong result\ngot: %#v\nwant: %#v", got, nil)
				}

			} else {
				if got != nil {
					t.Fatalf("wrong result\ngot: %#v\nwant: %#v", got, nil)
				}
			}
		})
	}
}

func TestFetchWorkspaceRunTask(t *testing.T) {
	org_name := "hashicorp"
	ws_name := "a-workspace"
	task_name := "a-task"

	tests := map[string]struct {
		org          string
		workspace    string
		taskName     string
		expectExists bool
		err          bool
	}{
		"non exisiting organization": {
			"not-an-org",
			ws_name,
			task_name,
			false,
			true,
		},
		"non exisiting workspace": {
			org_name,
			"not-a-workspace",
			task_name,
			false,
			true,
		},
		"non exisiting run task": {
			org_name,
			ws_name,
			"not-a-task",
			false,
			true,
		},
		"an existing workspace run task": {
			org_name,
			ws_name,
			task_name,
			true,
			false,
		},
	}

	client := testTfeClient(t, testClientOptions{defaultRunTaskID: "task-123"})
	task, _ := client.RunTasks.Create(nil, org_name, tfe.RunTaskCreateOptions{
		Name: task_name,
		URL:  runTasksUrl(),
	})
	ws, _ := client.Workspaces.Create(nil, org_name, tfe.WorkspaceCreateOptions{
		Name: &ws_name,
	})
	client.WorkspaceRunTasks.Create(nil, ws.ID, tfe.WorkspaceRunTaskCreateOptions{
		EnforcementLevel: tfe.Mandatory,
		RunTask:          task,
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := fetchWorkspaceRunTask(test.taskName, test.workspace, test.org, client)

			if (err != nil) != test.err {
				t.Fatalf("expected error is %t, got %v", test.err, err)
			}

			if test.expectExists {
				if got == nil || got.RunTask.Name != test.taskName {
					t.Fatalf("wrong result\ngot: %#v\nwant: %#v", got, nil)
				}

			} else {
				if got != nil {
					t.Fatalf("wrong result\ngot: %#v\nwant: %#v", got, nil)
				}
			}
		})
	}
}
