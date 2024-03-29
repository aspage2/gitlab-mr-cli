/*
Copyright © 2020

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/mintel/personal-dev/apage/glmr/internal/repository"
	"gitlab.com/mintel/personal-dev/apage/glmr/internal/userinput"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new MR in this repository.",
	Long: `Create a new merge request in this repository. Gathers MR parameters from the User and the current environment:
 * Title - MR title (editor input)
 * Description - MR description (editor input)
 * Source Branch - Current active branch
 * Target Branch - master
 * Assignee - choosing "yes" sets the assignee to the API token holder.

All fields may be overridden by the relevant CLI flag (see below).
`,
	Run: func(cmd *cobra.Command, args []string) {
		actualCreateCommand()
	},
}

var (
	DryRun      bool
	Title       string
	Description string
	MergeSource string
	MergeTarget string
)

func init() {
	rootCmd.AddCommand(createCmd)

	flagSet := createCmd.Flags()

	flagSet.BoolVar(&DryRun, "dry-run", false, "Log HTTP request without performing it.")
	flagSet.StringVar(&Title, "title", "", "Title for the MR.")
	flagSet.StringVar(&MergeSource, "source", "", "Branch to use as merge source. Defaults to current branch.")
	flagSet.StringVar(&MergeTarget, "target", "", "Branch to use as merge target.")
	flagSet.StringVar(&Description, "description", "", "Description for the MR. Prompts with a template if not provided.")

	flagSet.Bool("delete-source", false, "Delete source branch on merge. Overrides MROptions.DeleteSourceBranch")
	flagSet.Bool("squash-commits", false, "Squash commits on merge. Overrides MROptions.SquashCommits")

	viper.BindPFlag("MROptions.DeleteSourceBranch", flagSet.Lookup("delete-source"))
	viper.BindPFlag("MROptions.SquashCommits", flagSet.Lookup("squash-commits"))
}

// Panic on error, noop otherwise
func chk(err error) {
	if err != nil {
		log.Fatalf("%s: %s", red("FATAL"), err.Error())
	}
}

// ANSI-bold-green string
func green(s string) string {
	return fmt.Sprintf("\u001b[32;1m%s\u001b[0m", s)
}

// ANSI-bold-red string
func red(s string) string {
	return fmt.Sprintf("\u001b[31;1m%s\u001b[0m", s)
}

// Pretty-prints the json payload for a Create MR call
func printRequest(opt *gitlab.CreateMergeRequestOptions) {
	data, err := json.MarshalIndent(opt, "", "  ")
	chk(err)
	fmt.Println(string(data))
}

// Submit the Create MR request; display the MR url on success
func submitAndReport(client *gitlab.Client, opt *gitlab.CreateMergeRequestOptions, slug string) {
	// Submit the request
	mr, _, err := client.MergeRequests.CreateMergeRequest(slug, opt)
	chk(err)
	fmt.Printf("View online at %s\n", mr.WebURL)
}

// Entrypoint for `glmr create`
func actualCreateCommand() {
	client, err := gitlab.NewClient(viper.GetString("APIToken"))
	chk(err)
	inputMethod, err := userinput.UseEditor(viper.GetString("Editor"))
	chk(err)

	var opt gitlab.CreateMergeRequestOptions

	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	chk(err)

	slug := getRepoSlug(repo)
	fmt.Printf("Using %s as project slug\n", green(slug))

	opt.SourceBranch = gitlab.String(getSourceBranch(repo))
	fmt.Printf("Using %s as source branch\n", green(*opt.SourceBranch))

	opt.TargetBranch = gitlab.String(getTargetBranch(client, slug))
	fmt.Printf("Using %s as target branch\n\n", green(*opt.TargetBranch))

	if *opt.SourceBranch == *opt.TargetBranch {
		chk(errors.New("cannot merge a branch into itself: " + *opt.SourceBranch))
	}

	opt.Title = gitlab.String(getMRTitle())
	if *opt.Title == "" {
		chk(errors.New("title cannot be empty"))
	}

	opt.Description = gitlab.String(getMRDescription(inputMethod))

	opt.AssigneeID = gitlab.Int(getAssignee(client))

	opt.RemoveSourceBranch = gitlab.Bool(viper.GetBool("MROptions.DeleteSourceBranch"))
	opt.Squash = gitlab.Bool(viper.GetBool("MROptions.SquashCommits"))

	// If DryRun, Log the request & quit
	if DryRun {
		fmt.Println(red("-- DRY RUN --"))
		printRequest(&opt)
	} else {
		submitAndReport(client, &opt, slug)
	}

}

// --- getX routines for sourcing MR options ---

func getRepoSlug(repo *git.Repository) string {
	remote, err := repo.Remote("origin")
	chk(err)
	slug, err := repository.GetRepoSlug(remote.Config().URLs[0])
	chk(err)
	return slug
}

func getSourceBranch(repo *git.Repository) string {
	if MergeSource == "" {
		head, err := repo.Head()
		chk(err)
		if !head.Name().IsBranch() {
			chk(errors.New("can't determine merge source: HEAD is not a branch"))
		}
		items := strings.Split(head.Name().String(), "/")
		return items[len(items)-1]
	}
	return MergeSource
}

func getTargetBranch(gitlab *gitlab.Client, slug string) string {
	if MergeTarget == "" {
		proj, _, err := gitlab.Projects.GetProject(slug, nil)
		chk(err)
		return proj.DefaultBranch
	} else {
		return MergeTarget
	}
}

func getMRTitle() string {
	if Title == "" {
		return userinput.StdinPrompt("Please provide an MR title")
	} else {
		return Title
	}
}

func getMRDescription(inputMethod userinput.LargeInputStrategy) string {
	if Description == "" {
		var template string
		files, err := repository.GetRepoTemplates()
		if err != nil || len(files) == 0 {
			template = viper.GetString("DescriptionTemplate")
		} else if len(files) == 1 {
			fmt.Printf("\nFound MR template: %s\n", files[0])
			data, err := ioutil.ReadFile(files[0])
			chk(err)
			template = string(data)
		} else {
			fmt.Println("\nFound multiple MR templates in .gitlab/merge_request_templates:")
			choice := userinput.MultiChoicePrompt(files, "Choose an MR Template")
			data, err := ioutil.ReadFile(choice)
			chk(err)
			template = string(data)
		}
		description, err := userinput.LargeInput(template, inputMethod)
		chk(err)
		return description
	} else {
		return Description
	}
}

func getAssignee(client *gitlab.Client) int {
	user, _, err := client.Users.CurrentUser()
	chk(err)
	assignSelf, err := userinput.YesOrNo(fmt.Sprintf("Assign MR to %s (token holder)?", green(user.Username)), true)
	chk(err)
	if assignSelf {
		return user.ID
	}
	return 0
}
