# Scoring Audit

Using for auding your score

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
➜  d-audit git:(main) ✗ ./daudit audit --username <username>
Analyzing 3542 issues
+---------+----------------------------------------+--------------------------------+
|  NAME   |                  URL                   |            PROBLEMS            |
+---------+----------------------------------------+--------------------------------+
| CS-502  | https://jira.vccloud.vn/browse/CS-502  | Missing Rate                   |
| GATE-17 | https://jira.vccloud.vn/browse/GATE-17 | Done but doesn't back to do-er |
| BK-310  | https://jira.vccloud.vn/browse/BK-310  | Done but doesn't back to do-er |
| BKE-247 | https://jira.vccloud.vn/browse/BKE-247 | Missing Due Date               |
| CS-271  | https://jira.vccloud.vn/browse/CS-271  | Missing Category               |
| CS-219  | https://jira.vccloud.vn/browse/CS-219  | Missing Due Date, Missing      |
|         |                                        | Category, Missing start date   |
+---------+----------------------------------------+--------------------------------+
```

#### Find out in-reviewed tasks

1. Feature:
   - Find out your tasks which has being reviewed
2. Example:

```shell
➜  d-audit git:(main) ✗ ./daudit in-review --username <username> 
Analyzing 3542 issues
+--------+---------------------------------------+----------+-----------+
| ISSUE  |                  URL                  | ASSIGNEE |  STATUS   |
+--------+---------------------------------------+----------+-----------+
| CS-508 | https://jira.vccloud.vn/browse/CS-508 | sapd     | In Review |
+--------+---------------------------------------+----------+-----------+
```

### Help

```shell
➜  d-audit git:(main) ✗ daudit  
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
