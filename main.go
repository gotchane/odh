package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"github.com/c-bata/go-prompt"
	"strings"
)

var (
	profile  string
	region   string
	dryRun   bool
	suggests []prompt.Suggest
	stackID  string
	appID    string
	command  string
)

// complete executes go-prompt suggestions
func complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}
	return argumentsCompleter(args)
}

func argumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(suggests, args[0], true)
	}

	first := args[0]
	second := args[1]
	if len(args) == 2 {
		for _, v := range suggests {
			if v.Text == first {
				stackID = v.Description
			}
		}
		return prompt.FilterContains(fetchStackApps(stackID), second, true)
	}

	third := args[2]
	if len(args) == 3 {
		appID = second
		subcommands := []prompt.Suggest{
			{Text: "deploy"},
		}
		return prompt.FilterHasPrefix(subcommands, third, true)
	}
	return []prompt.Suggest{}
}

func fetchStackApps(stackID string) []prompt.Suggest {
	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profile}))
	svc := opsworks.New(sess, aws.NewConfig().WithRegion(region))
	result, err := svc.DescribeApps(&opsworks.DescribeAppsInput{
		StackId: &stackID,
	})
	if err != nil {
		panic(err)
	}
	var apps []prompt.Suggest
	for _, v := range result.Apps {
		apps = append(apps, prompt.Suggest{
			Text:        aws.StringValue(v.AppId),
		})
	}
	return apps
}

func fetchSuggestStacks(profile, region string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profile}))
	svc := opsworks.New(sess, aws.NewConfig().WithRegion(region))
	result, err := svc.DescribeStacks(nil)
	if err != nil {
		panic(err)
	}
	for _, b := range result.Stacks {
		suggests = append(suggests, prompt.Suggest{
			Text:        aws.StringValue(b.Name),
			Description: aws.StringValue(b.StackId),
		})
	}
}

func executeDeploy() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profile}))
	svc := opsworks.New(sess, aws.NewConfig().WithRegion(region))
	str := "deploy"
	result, err := svc.CreateDeployment(&opsworks.CreateDeploymentInput{
		StackId: &stackID,
		AppId:   &appID,
		Command: &opsworks.DeploymentCommand{
			Name: &str,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func main() {
	flag.StringVar(&profile, "p", "", "Aws profile")
	flag.StringVar(&region, "r", "ap-northeast-1", "Aws region")
	flag.BoolVar(&dryRun, "d", false, "Dry run")
	flag.Parse()

	if profile == "" {
		fmt.Println("Not existing profile:", profile)
		flag.Usage()
		return
	}

	fetchSuggestStacks(profile, region)

	in := prompt.Input(
		">>> ",
		complete,
		prompt.OptionTitle("opsworks-helper"),
		prompt.OptionHistory(nil),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionSuggestionBGColor(prompt.Blue),
	)
	if (dryRun == true) {
		fmt.Println(in)
	} else {
		executeDeploy()
	}
}
