
# GLMR: Gitlab Merge Requests

Cli tool for creating merge requests for gitlab.

## Usage
When you are ready to create a new merge request, run `glmr new`
```
glmr new [OPTIONS]
```
This will create a file called `.glmr.json`
#### Options
 * `--title` -  the title of your MR. Defaults to the current commit message
 * `--source` - the source branch. Defaults to the current active branch
 * `--target` - the target branch. Defaults to the default branch of the project (usually master)
