package tfe

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
)

// fetchOrganizationRunTask returns the task in an organization by name
func fetchOrganizationRunTask(name string, organization string, client *tfe.Client) (*tfe.RunTask, error) {
	options := &tfe.RunTaskListOptions{}
	for {
		list, err := client.RunTasks.List(ctx, organization, options)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving organization tasks: %v", err)
		}

		for _, task := range list.Items {
			if task != nil && task.Name == name {
				return task, nil
			}
		}

		// Exit the loop when we've seen all pages.
		if list.CurrentPage >= list.TotalPages {
			break
		}

		// Update the page number to get the next page.
		options.PageNumber = list.NextPage
	}

	return nil, fmt.Errorf("Could not find organization run task for organization %s and name %s", organization, name)
}
