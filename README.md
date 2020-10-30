
# GLMR: Gitlab Merge Requests

Cli tool for creating merge requests for gitlab.

## Installation & Setup

### Step 1: Install
#### Method 1: Download pre-built binary
Download the [latest release](https://gitlab.com/mintel/personal-dev/apage/gitlab-mr-cli/-/releases)
for your operating system.
```sh
cd ~/Downloads
chmod +x glmr_X_Y
sudo mv glmr_X_Y /usr/local/bin/glmr
```

#### Method 2: Build from source
Clone this repo to your machine and build locally with make:
```sh
git clone --depth=1 git@gitlab.com:mintel/personal-dev/apage/gitlab-mr-cli.git
cd gitlab-mr-cli
make install
```
---
Once installed, you can now invoke it with `glmr`:
```
~ $ glmr
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

### Step 2: Setup

Generate an API key for your gitlab user, give it the `api` scope: [gitlab token settings](https://gitlab.com/profile/personal_access_tokens) Save that token along with other user settings in `$HOME/.glmr.yaml`:
```yaml
APIToken: your-api-token-here
Editor: vim
DescriptionTemplate: |
  # MR description
  What do you want the reviewers to review?
```

An example configuration exists in `examples/example.glmr.yaml`. 

## Usage

GLMR creates merge requests via the Gitlab API. It supports the following MR options:

* `--title` - MR title
* `--description` - MR description
* `--source` - Branch name of merge source
* `--target` - Branch name of merge target

When those flags aren't provided, GLMR sources their values in the following ways:

| Option             | Source                                                       |
| ------------------ | ------------------------------------------------------------ |
| **MR Title**       | Prompts the user                                             |
| **MR Description** | Prompts the user. Uses .glmr.yaml.DescriptionTemplate as a template. |
| **Source Branch**  | Uses the currently checked-out branch                        |
| **Target Branch**  | Defaults to "master"                                         |
| **Assignee**       | Prompts for whether or not to assign the token holder to the MR |
| **Project Slug**   | Parses the project's remote and extracts the Gitlab project-slug from the remote URL. ex. `https://gitlab.com/mintel/everest/mdl-planner.git -> mintel/everest/mdl-planner` |

### Dry Runs

Use the `--dry-run` flag to force Glmr to print the request and exit. Glmr will still perform api *reads* but will refrain from any API *writes*. Example, it will `GET /user` to get the token holder's ID, but it will **not** `POST /project/100/merge_requests`.
