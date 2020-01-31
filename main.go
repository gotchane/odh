package main

import (
	"fmt"
	"flag"
	"strings"
    "github.com/c-bata/go-prompt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/opsworks"
)

var (
	profile string
	region string
	suggests []prompt.Suggest
	stackId string
	appId string
	command string
)

type Completer struct {
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}

	return c.argumentsCompleter(args)
}

var commands = []prompt.Suggest{
	{Text: "deploy", Description: "deploy app"},
	{Text: "undeploy", Description: "undeploy app"},
}

func (c *Completer) argumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}

	command := args[0]
	switch command {
	case "deploy", "undeploy":
		second := args[1]
		if len(args) == 2 {
			return prompt.FilterContains(suggests, second, true)
		}

		third := args[2]
		if len(args) == 3 {
			var stackId string
			for _, v := range suggests {
				if v.Text == second {
					stackId = v.Description
				}
			}
			return prompt.FilterHasPrefix(fetchStackApps(stackId), third, true)
		}
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}

func fetchStackApps(stackId string) []prompt.Suggest {
	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile:profile}))
	svc := opsworks.New(sess, aws.NewConfig().WithRegion(region))
	result, err := svc.DescribeApps(&opsworks.DescribeAppsInput{
		StackId: &stackId,
	})
	if err != nil {
		panic(err)
	}
	var apps []prompt.Suggest
	for _, v := range result.Apps {
		apps = append(apps, prompt.Suggest{
			Text: aws.StringValue(v.AppId),
			Description: aws.StringValue(v.Name),
		})
	}
	return apps
}

func fetchSuggestStacks(profile, region string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile:profile}))
	svc := opsworks.New(sess, aws.NewConfig().WithRegion(region))
	result, err := svc.DescribeStacks(nil)
	if err != nil {
		panic(err)
	}
	for _, b := range result.Stacks {
		suggests = append(suggests, prompt.Suggest{
			Text: aws.StringValue(b.Name),
			Description: aws.StringValue(b.StackId),
		})
	}
}

func main() {
	flag.StringVar(&profile, "p", "", "Aws profile")
	flag.StringVar(&region, "r", "ap-northeast-1", "Aws region")
	flag.Parse()

	if (profile == "") {
		fmt.Println("Not existing profile:", profile)
		flag.Usage()
		return
	}

	fetchSuggestStacks(profile, region)

    // in := prompt.Input(
	// 	">>> ",
	// 	completer,
	// 	prompt.OptionTitle("opsworks-helper"),
	// 	prompt.OptionHistory(nil),
	// 	prompt.OptionPrefixTextColor(prompt.Yellow),
	// 	prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
	// 	prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
	// 	prompt.OptionSuggestionTextColor(prompt.Black),
	// 	prompt.OptionSuggestionBGColor(prompt.Blue),
	// )
    // fmt.Println("Your input: " + in)
	c := Completer{}
    in := prompt.Input(
		">>> ",
		c.Complete,
		prompt.OptionTitle("opsworks-helper"),
		prompt.OptionHistory(nil),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionSuggestionBGColor(prompt.Blue),
	)
    fmt.Println("Your input: " + in)
}
