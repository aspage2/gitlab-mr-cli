/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"strings"

	git "github.com/libgit2/git2go/v27"

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

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&DryRun, "dry-run", false, "Log HTTP request without performing it.")
	createCmd.Flags().StringVar(&Title, "title", "", "Title for the MR.")
	createCmd.Flags().StringVar(&MergeSource, "source", "", "Branch to use as merge source. Defaults to current branch.")
	createCmd.Flags().StringVar(&MergeTarget, "target", "master", "Branch to use as merge target.")
	createCmd.Flags().StringVar(&Description, "description", "", "Description for the MR. Prompts with a template if not provided.")
}

func chk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	DryRun      bool
	Title       string
	Description string
	MergeSource string
	MergeTarget string
)

func actualCreateCommand() {
	client, err := gitlab.NewClient(viper.GetString("APIToken"))
	chk(err)
	inputMethod, err := userinput.UseEditor(viper.GetString("Editor"))
	chk(err)

	var opt gitlab.CreateMergeRequestOptions

	repo, err := git.OpenRepository(".")
	chk(err)

	slug := getRepoSlug(repo)
	fmt.Printf("Using \u001b[32;1m%s\u001b[0m as project slug\n", slug)

	opt.SourceBranch = gitlab.String(getSourceBranch(repo))
	fmt.Printf("Using \u001b[32;1m%s\u001b[0m as source branch\n", *opt.SourceBranch)

	opt.TargetBranch = gitlab.String(getTargetBranch())
	fmt.Printf("Using \u001b[32;1m%s\u001b[0m as target branch\n", *opt.TargetBranch)

	opt.Title = gitlab.String(getMRTitle(inputMethod))

	opt.Description = gitlab.String(getMRDescription(inputMethod))

	opt.AssigneeID = gitlab.Int(getAssignee(client))

	// If DryRun, Log the request & quit
	if DryRun {
		fmt.Println("\n\u001b[31;1m -- DRY RUN -- \u001b[0m")
		printRequest(&opt)
	} else {
		submitAndReport(client, &opt, slug)
	}

}

func getRepoSlug(repo *git.Repository) string {
	remote, err := repo.Remotes.Lookup("origin")
	chk(err)
	slug, err := repository.GetRepoSlug(remote.Url())
	chk(err)
	return slug
}

func getSourceBranch(repo *git.Repository) string {
	if MergeSource == "" {
		head, err := repo.Head()
		chk(err)
		if !head.IsBranch() {
			chk(errors.New("can't determine merge source: HEAD is not a branch"))
		}
		items := strings.Split(head.Name(), "/")
		return items[len(items)-1]
	}
	return MergeSource
}

func getTargetBranch() string {
	if MergeTarget == "" {
		return "master"
	} else {
		return MergeTarget
	}
}

func getMRTitle(inputMethod userinput.LargeInputStrategy) string {
	if Title == "" {
		title, err := userinput.LargeInput("Please provide an MR title", inputMethod)
		chk(err)
		return strings.TrimSpace(title)
	} else {
		return Title
	}
}

func getMRDescription(inputMethod userinput.LargeInputStrategy) string {
	if Description == "" {
		template := viper.GetString("DescriptionTemplate")
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
	assignSelf, err := userinput.YesOrNo(fmt.Sprintf("Assign %s (token holder)?", user.Username), true)
	chk(err)
	if assignSelf {
		return user.ID
	}
	return 0
}

func printRequest(opt *gitlab.CreateMergeRequestOptions) {
	data, err := json.MarshalIndent(opt, "", "  ")
	chk(err)
	fmt.Println(string(data))
}

func submitAndReport(client *gitlab.Client, opt *gitlab.CreateMergeRequestOptions, slug string) {
	// Submit the request
	mr, _, err := client.MergeRequests.CreateMergeRequest(slug, opt)
	chk(err)
	fmt.Printf("View online at %s\n", mr.WebURL)
}
