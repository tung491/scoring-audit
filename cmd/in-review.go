package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tung491/đ-audit/formatter"
)

var listIssueHeader = []string{"Issue", "URL", "Assignee", "Status"}

var listInReviewTask = &cobra.Command{
	Use:   "in-review",
	Short: "List your in-review tasks",
	Run: func(cmd *cobra.Command, args []string) {
		issues := listIssues()
		me := getUserInfo()
		var data [][]string
		channel := make(chan []string)
		for _, issue := range issues {
			go getInReviewedTask(issue, me, channel)
		}

		for count := 0; count < len(issues); count++ {
			select {
			case s := <-channel:
				if len(s) > 0 {
					data = append(data, s)
				}
			}
		}
		formatter.Output(listIssueHeader, data)
	},
}

func getInReviewedTask(issue Issue, me User, channel chan []string) {
	fields := issue.Fields
	if isReviewedTask(me, issue) && issue.Fields.Status.Name == "In Review" {
		url := "https://jira.vccloud.vn/browse/" + issue.Key
		channel <- []string{issue.Key, url, fields.Assignee.Name, fields.Status.Name}
	}
	channel <- []string{}
}

func init() {
	rootCmd.AddCommand(listInReviewTask)
}
