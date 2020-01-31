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
	suggests []prompt.Suggest
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
	{Text: "start", Description: "Display one or many resources"},
}

func (c *Completer) argumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(suggests, args[0], true)
	}

	second := args[1]
	if len(args) == 2 {
		subcommands := []prompt.Suggest{
			{Text: "get"},
		}
		return prompt.FilterHasPrefix(subcommands, second, true)
	}

	third := args[2]
	if len(args) == 3 {
		switch second {
		case "g", "get", "gets":
			return prompt.FilterContains(getSubcommandsuggestions(), third, true)
		}
	}
	return []prompt.Suggest{}
}

func getSubcommandsuggestions() []prompt.Suggest {
	s := make([]prompt.Suggest, 2)
	s[0] = prompt.Suggest{Text: "hoge"}
	s[1] = prompt.Suggest{Text: "fuga"}
	return s
}

func fetchSuggestStacks(profile, region string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile:profile}))
	svc := opsworks.New(sess, aws.NewConfig().WithRegion(region))
	result, err := svc.DescribeStacks(nil)
	if err != nil {
		panic(err)
	}
	for _, b := range result.Stacks {
		suggests = append(suggests, prompt.Suggest{Text: aws.StringValue(b.Name)})
	}
}

func completer(in prompt.Document) []prompt.Suggest {
    return prompt.FilterHasPrefix(suggests, in.GetWordBeforeCursor(), true)
}

func main() {
	var (
		profile string
		region string
	)

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
