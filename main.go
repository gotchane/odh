package main

import (
	"fmt"
	"flag"
    "github.com/c-bata/go-prompt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/opsworks"
)

var (
	suggests []prompt.Suggest
)

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

    in := prompt.Input(
		">>> ",
		completer,
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
