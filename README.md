
# GLMR: Gitlab Merge Requests

Cli tool for creating merge requests for gitlab.

## Installation & Setup

Download the [latest release](https://gitlab.com/mintel/personal-dev/apage/gitlab-mr-cli/-/releases)
for your operating system. `chmod +x` and copy the executable to somewhere on your PATH
as a file called `glmr`. You can now invoke it from your shell:

```
~/everest $ glmr
GLMR is a cli tool for creating merge requests from the CLI.

Usage:
  glmr [command]

Available Commands:
  create      Create a new MR in this repository.
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.glmr.yaml)
  -h, --help            help for glmr

Use "glmr [command] --help" for more information about a command.
```

Generate an API key for your gitlab user, give it the `api` scope: [gitlab token settings](https://gitlab.com/profile/personal_access_tokens)

Save that token along with other user settings in `$HOME/.glmr.yaml`:
```yaml
APIToken: your-api-token-here
Editor: vim
DescriptionTemplate: |
  # MR description
  What do you want the reviewers to review?
```

## Configuration
GLMR looks for `.glmr.yaml` in either `$HOME` or your current directory.

| Name              | Descirption                                            | Notes                                                        |
| ----------------- | ------------------------------------------------------ | ------------------------------------------------------------ |
| `APIToken`        | Gitlab API token.                                      | Should have `api` scope to be able to write to the gitlab API. |
| `Editor`      | Editor to use for large user inputs                  | Supported editors: `vim`, `nano`, `vscode`, `typora`         |
| `DescriptionTemplate` | Template to show when prompting for the MR description |                                                              |

## Usage

Glmr creates merge requests via the Gitlab API. It supports the following MR options:

* `--title` - MR title
* `--description` - MR description
* `--source` - Branch name of merge source
* `--target` - Branch name of merge target


### Default Behaviors
When not provided, Glmr has default methods for sourcing each field:

**Title** - Prompts the user

**Description** - Prompts the user. Uses .glmr.yaml.DescriptionTemplate as a template.

**Source** - scrapes `.git/` for the current active branch

**Target** - always defaults to "master"

**Assignee** - Prompts the user whether or not the API token holder should be assigned.

**Project Slug** - Scrapes `.git/config` for the remote URL, and uses the path as the slug (ex. `https://gitlab.com/mintel/everest/mdl-planner.git > mintel/everest/mdl-planner`)

```
 ~/everest/ops/glmr (lotta-stuff) $ glmr create --dry-run
Using config file: /home/apage/.glmr.yaml
Using mintel/personal-dev/apage/gitlab-mr-cli as project ID
Using lotta-stuff as source branch
Using master as target branch

# prompt user for title....

# prompt user for description...

Assign to user apage1 (token holder)? [Y/n] y

 -- DRY RUN -- 
POST https://gitlab.com/api/v4/projects/mintel%2Fpersonal-dev%2Fapage%2Fgitlab-mr-cli/merge_requests
Content-Type: [application/json]
Authorization: [Bearer (----)]
{
  "assignee_id": 00000000,
  "description": "## Goal\nWhat is the goal of the MR?\n\n## Changes\nExplain the diff as a list of high-level changes\n\n## Notes\nAdditional notes \u0026 concerns for the reviewer\n",
  "source_branch": "lotta-stuff",
  "target_branch": "master",
  "title": "Please provide an MR title"
}
```

### Dry Runs

Use the `--dry-run` flag to force Glmr to print the request and exit. Glmr will still perform api *reads* but will refrain from any API *writes*. Example, it will `GET /user` to get the token holder's ID, but it will **not** `POST /project/100/merge_requests`.
