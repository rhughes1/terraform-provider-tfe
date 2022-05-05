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
