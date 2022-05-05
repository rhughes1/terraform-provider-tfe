package tfe

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type WorkspaceRunTaskNamesKey struct {
	workspaceID, taskID string
}

// newMockWorkspaceRunTasks creates a mock WorkspaceRunTasks implementation. Any created
// WorkspaceRunTasks will have the id given in defaultWorkspaceRunTaskID.
func newMockWorkspaceRunTasks(options testClientOptions) *mockWorkspaceRunTasks {
	return &mockWorkspaceRunTasks{
		options:               options,
		WorkspaceRunTaskNames: make(map[WorkspaceRunTaskNamesKey]*tfe.WorkspaceRunTask),
	}
}

type mockWorkspaceRunTasks struct {
	options               testClientOptions
	WorkspaceRunTaskNames map[WorkspaceRunTaskNamesKey]*tfe.WorkspaceRunTask
}

func (m *mockWorkspaceRunTasks) List(ctx context.Context, workspaceID string, options *tfe.WorkspaceRunTaskListOptions) (*tfe.WorkspaceRunTaskList, error) {
	if workspaceID != m.options.defaultWorkspaceID {
		return nil, tfe.ErrInvalidWorkspaceID
	}

	list := tfe.WorkspaceRunTaskList{
		Items: make([]*tfe.WorkspaceRunTask, 0),
	}

	for _, task := range m.WorkspaceRunTaskNames {
		if task.Workspace.ID == workspaceID {
			list.Items = append(list.Items, task)
		}
	}

	list.Pagination = &tfe.Pagination{
		CurrentPage: 1,
		TotalPages:  1,
		TotalCount:  len(list.Items),
	}

	return &list, nil
}

func (m *mockWorkspaceRunTasks) Read(ctx context.Context, workspaceID, workspaceTaskID string) (*tfe.WorkspaceRunTask, error) {
	panic("not implemented")
}

func (m *mockWorkspaceRunTasks) Create(ctx context.Context, workspaceID string, options tfe.WorkspaceRunTaskCreateOptions) (*tfe.WorkspaceRunTask, error) {

	task := &tfe.WorkspaceRunTask{
		ID:               m.options.defaultWorkspaceRunTaskID,
		EnforcementLevel: options.EnforcementLevel,
		RunTask:          options.RunTask,
		Workspace: &tfe.Workspace{
			ID: workspaceID,
		},
	}

	m.WorkspaceRunTaskNames[WorkspaceRunTaskNamesKey{workspaceID, task.ID}] = task

	return task, nil
}

func (m *mockWorkspaceRunTasks) Update(ctx context.Context, workspaceID, workspaceTaskID string, options tfe.WorkspaceRunTaskUpdateOptions) (*tfe.WorkspaceRunTask, error) {
	panic("not implemented")
}

func (m *mockWorkspaceRunTasks) Delete(ctx context.Context, workspaceID, workspaceTaskID string) error {
	panic("not implemented")
}
