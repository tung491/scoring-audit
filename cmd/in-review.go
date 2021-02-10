package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tung491/đ-audit/formatter"
	"sync"
)

var listIssueHeader = []string{"Issue", "URL", "Assignee", "Status"}

var listInReviewTask = &cobra.Command{
	Use:   "in-review",
	Short: "List your in-review tasks",
	Run: func(cmd *cobra.Command, args []string) {
		issues := listIssues()
		me := getUserInfo()
		var data [][]string
		wg := new(sync.WaitGroup)
		wg.Add(len(issues))
		for _, issue := range issues {
			go getInReviewedTask(wg, &data, issue, me)
		}
		wg.Wait()
		formatter.Output(listIssueHeader, data)
	},
}

func getInReviewedTask(wg *sync.WaitGroup, data *[][]string, issue Issue, me User) {
	defer wg.Done()
	fields := issue.Fields
	if isReviewedTask(me, issue) && issue.Fields.Status.Name == "In Review" {
		url := "https://jira.vccloud.vn/browse/" + issue.Key
		*data = append(*data, []string{issue.Key, url, fields.Assignee.Name, fields.Status.Name})
	}
}


func init() {
	rootCmd.AddCommand(listInReviewTask)
}
