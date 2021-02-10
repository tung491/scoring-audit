package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tung491/đ-audit/formatter"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var (
	auditHeader = []string{"Name", "URL", "Problems"}
)

const (
	defaultMaxResults = 50
)

type Issue struct {
	ID     string `json:"id"`
	URL    string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		IssueType  Field  `json:"issueType"`
		Project    Field  `json:"project"`
		Resolution Field  `json:"resolution"`
		Level      Field  `json:"customfield_10226,omitempty"`
		Category   Field  `json:"customfield_10303,omitempty"`
		TaskType   Field  `json:"customfield_10304,omitempty"`
		Rate       Field  `json:"customfield_10227,omitempty"`
		Assignee   Field  `json:"assignee"`
		Status     Field  `json:"status"`
		DueDate    string `json:"duedate,omitempty"`
		FinishDate string `json:"customfield_10210,omitempty"`
		StartDate  string `json:"customfield_10209,omitempty"`
	} `json:"fields"`
	ChangeLog ChangeLog `json:"changelog,omitempty"`
}

type Field struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChangeLog struct {
	Histories []Log `json:"histories"`
}

type Log struct {
	ID     string `json:"id"`
	Author Author `json:"author"`
	Items  []Item `json:"items"`
}

type Author struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}

type Item struct {
	Field      string `json:"field"`
	FieldType  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}

type User struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}

func getUserInfo() User {
	url := "https://jira.vccloud.vn/rest/api/2/user?username=" + userName

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(userName, token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}
	return user
}

func getSoftwareProjects() []string {
	client := &http.Client{}
	url := "https://jira.vccloud.vn/rest/api/2/project"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(userName, token)
	var projects []struct {
		Name        string `json:"name"`
		ProjectType string `json:"projectTypeKey"`
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&projects); err != nil {
		log.Fatal(err)
	}
	softwareProject := make([]string, 0)
	for _, project := range projects {
		if project.ProjectType == "software" {
			softwareProject = append(softwareProject, project.Name)
		}
	}
	return softwareProject
}

type IssuesResponse struct {
	Total      int     `json:"total"`
	MaxResults int     `json:"maxResults"`
	Issues     []Issue `json:"issues"`
}

func listIssues() []Issue {
	projects := getSoftwareProjects()
	resp := getIssues(projects, 0)
	issues := resp.Issues
	channel := make(chan IssuesResponse)
	pageCount := int(math.Ceil(float64(resp.Total/resp.MaxResults))) + 1
	for page := 2; page <= pageCount; page++ {
		go getIssuesChannel(projects, defaultMaxResults*(page-1), channel)
	}

	for count := 0; count < pageCount-1; count++ {
		select {
		case resp := <-channel:
			issues = append(issues, resp.Issues...)
		}
	}
	fmt.Printf("Analyzing %d issues\n", len(issues))
	return issues
}

func getIssues(projects []string, startAt int) IssuesResponse {
	url := "https://jira.vccloud.vn/rest/api/2/search"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	query := req.URL.Query()
	query.Add("startAt", strconv.Itoa(startAt))
	jqlQuery := fmt.Sprintf("project in (\"%s\")", strings.Join(projects, "\",\""))
	query.Add("jql", jqlQuery)
	query.Add("expand", "changelog")
	req.URL.RawQuery = query.Encode()
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(userName, token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var issues struct {
		Total      int     `json:"total"`
		MaxResults int     `json:"maxResults"`
		Issues     []Issue `json:"issues"`
	}
	if err := json.NewDecoder(res.Body).Decode(&issues); err != nil {
		log.Fatal(err)
	}
	return issues
}

func getIssuesChannel(projects []string, startAt int, channel chan IssuesResponse) {
	channel <- getIssues(projects, startAt)
}

func isReviewedTask(me User, issue Issue) bool {
	for _, history := range issue.ChangeLog.Histories {
		items := history.Items
		if len(items) == 2 {
			item1, item2 := items[0], items[1]
			if item1.Field == "assignee" && item1.From == me.Key &&
				item2.FromString == "In Progress" && item2.ToString == "In Review" {
				return true
			}
		}
	}
	return false
}

func auditTask(me User, issue Issue) []string {
	var problems []string
	fields := issue.Fields
	if fields.Resolution.Name == "Done" && fields.IssueType.Name == "Task" {
		if fields.Assignee.Name == userName {
			if fields.DueDate == "" {
				problems = append(problems, "Missing Due Date")
			}
			if fields.TaskType == (Field{}) || fields.Category == (Field{}) {
				problems = append(problems, "Missing Category")
			}
			if fields.Rate == (Field{}) {
				problems = append(problems, "Missing Rate")
			}
			if fields.FinishDate == "" {
				problems = append(problems, "Missing finish date")
			}
			if fields.StartDate == "" {
				problems = append(problems, "Missing start date")
			}
			if fields.Level == (Field{}) {
				problems = append(problems, "Missing Level")
			}
		} else if isReviewedTask(me, issue) {
			problems = append(problems, "Done but doesn't back to do-er")
		}

		url := "https://jira.vccloud.vn/browse/" + issue.Key
		if len(problems) > 0 {
			return []string{issue.Key, url, strings.Join(problems, ", ")}
		}
	}
	return []string{}

}

func auditTaskChannel(me User, issue Issue, channel chan []string) {
	channel <- auditTask(me, issue)
}

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit Task",
	Run: func(cmd *cobra.Command, args []string) {
		issues := listIssues()
		me := getUserInfo()
		var data [][]string
		channel := make(chan []string)
		for _, issue := range issues {
			go auditTaskChannel(me, issue, channel)
		}
		count := 0
		for count < len(issues) {
			select {
			case s := <-channel:
				if len(s) > 0 {
					data = append(data, s)
				}
				count++
			}
		}
		formatter.Output(auditHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}
