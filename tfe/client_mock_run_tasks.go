package tfe

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type runTaskNamesKey struct {
	organization, taskID string
}

// newMockRunTasks creates a mock runTasks implementation. Any created
// runTasks will have the id given in defaultRunTaskID.
func newMockRunTasks(options testClientOptions) *mockRunTasks {
	return &mockRunTasks{
		options:      options,
		runTaskNames: make(map[runTaskNamesKey]*tfe.RunTask),
	}
}

type mockRunTasks struct {
	options      testClientOptions
	runTaskNames map[runTaskNamesKey]*tfe.RunTask
}

func (m *mockRunTasks) Create(ctx context.Context, organization string, options tfe.RunTaskCreateOptions) (*tfe.RunTask, error) {
	task := &tfe.RunTask{
		ID:   m.options.defaultRunTaskID,
		Name: options.Name,
		Organization: &tfe.Organization{
			Name: organization,
		},
	}

	m.runTaskNames[runTaskNamesKey{organization, options.Name}] = task

	return task, nil
}

func (m *mockRunTasks) List(ctx context.Context, organization string, options *tfe.RunTaskListOptions) (*tfe.RunTaskList, error) {
	if organization != m.options.defaultOrganization {
		return nil, tfe.ErrInvalidOrg
	}

	list := tfe.RunTaskList{
		Items: make([]*tfe.RunTask, 0),
	}

	for _, task := range m.runTaskNames {
		if task.Organization.Name == m.options.defaultOrganization {
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

func (m *mockRunTasks) Read(ctx context.Context, runTaskID string) (*tfe.RunTask, error) {
	panic("not implemented")
}

func (m *mockRunTasks) ReadWithOptions(ctx context.Context, runTaskID string, options *tfe.RunTaskReadOptions) (*tfe.RunTask, error) {
	panic("not implemented")
}

func (m *mockRunTasks) Update(ctx context.Context, runTaskID string, options tfe.RunTaskUpdateOptions) (*tfe.RunTask, error) {
	panic("not implemented")
}

func (m *mockRunTasks) Delete(ctx context.Context, runTaskID string) error {
	panic("not implemented")
}

func (m *mockRunTasks) AttachToWorkspace(ctx context.Context, workspaceID, runTaskID string, enforcement tfe.TaskEnforcementLevel) (*tfe.WorkspaceRunTask, error) {
	panic("not implemented")
}
