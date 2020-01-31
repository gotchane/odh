package main

import(
	"fmt"
	"flag"
)

func main() {
	var profile string
	flag.StringVar(&profile, "p", "", "Aws profile")
	flag.Parse()

	if (profile != "staging") {
		fmt.Println("Not existing profile:", profile)
		flag.Usage()
		return
	}
	fmt.Println(profile)
}
