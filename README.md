# Scoring Audit
Using for auding your score in https://scoring.vccloud.vn/home

## Install
- Clone the source code:

`git clone git@github.com:tung491/scoring-audit.git`

- Install golang
- Build 

`go build -o daudit main.go`

### Configure
- Create .jira.yaml in `$HOME`
- Configure username and token. 
```yaml
username: tung491
token: abcxyz
```

### Usage
#### Audit tasks
1. Feature:
   - Check missing field (Rate, Due date, Level, etc.)
   - Find out tasks which have been done but reviewer doesn't back them to you
2. Example:
```shell
➜  đ-audit git:(main) ✗ daudit audit
Analyzing 3542 issues
+---------+----------------------------------------+--------------------------------+
|  NAME   |                  URL                   |            PROBLEMS            |
+---------+----------------------------------------+--------------------------------+
| BKE-247 | https://jira.vccloud.vn/browse/BKE-247 | Missing Due Date               |
| CS-219  | https://jira.vccloud.vn/browse/CS-219  | Missing Due Date, Missing      |
|         |                                        | Category, Missing start date   |
| CS-502  | https://jira.vccloud.vn/browse/CS-502  | Missing Rate                   |
| CS-271  | https://jira.vccloud.vn/browse/CS-271  | Missing Category               |
+---------+----------------------------------------+--------------------------------+
```
#### Find out in-reviewed tasks
1. Feature:
   - Find out your tasks which has being reviewed
2. Example:
```shell
➜  đ-audit git:(main) ✗ daudit in-review
Analyzing 3542 issues
+-------+-----+----------+--------+
| ISSUE | URL | ASSIGNEE | STATUS |
+-------+-----+----------+--------+
+-------+-----+----------+--------+
```

### Help
```shell
➜  đ-audit git:(main) ✗ daudit  
Audit tasks which missing components

Usage:
  audit [command]

Available Commands:
  audit       Audit Task
  help        Help about any command
  in-review   List your in-review tasks

Flags:
      --config string   config file (default is $HOME/.jira.yaml)
  -h, --help            help for audit
```
