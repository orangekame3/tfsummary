/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfsummary",
	Short: "Summarize terraform plan output",
	Run:   runTfsummary,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

func runTfsummary(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)
	create := make([]string, 0)
	destroy := make([]string, 0)
	update := make([]string, 0)
	replace := make([]string, 0)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, "Enter a value:") {
			continue
		}
		fmt.Print(line)
		if strings.Contains(line, "Only 'yes'") {
			fmt.Print("Enter a value: ")
			continue
		}

		resource, action := extractResourceAction(line)
		fmt.Print(resource, action)
		if resource != "" {
			switch action {
			case "will be created":
				create = append(create, resource)
			case "will be destroyed":
				destroy = append(destroy, resource)
			case "will be updated in-place":
				update = append(update, resource)
			case "must be replaced":
				replace = append(replace, resource)
			}
		}
	}
	summary := Summary{
		Create:  create,
		Destroy: destroy,
		Update:  update,
		Replace: replace,
	}
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("## summary")
	if len(summary.Create) > 1 {
		fmt.Println("## create")
		for _, v := range summary.Create {
			fmt.Println("- ", v)
		}
	}
	if len(summary.Destroy) > 1 {
		fmt.Println("## destroy")
		for _, v := range summary.Destroy {
			fmt.Println("- ", v)
		}
	}
	if len(summary.Update) > 1 {
		fmt.Println("## update")
		for _, v := range summary.Update {
			fmt.Println("- ", v)
		}
	}
	if len(summary.Replace) > 1 {
		fmt.Println("## replace")
		for _, v := range summary.Replace {
			fmt.Println("- ", v)
		}
	}
	fmt.Println("------------------------------------------------------------------")
}

func extractResourceAction(line string) (string, string) {
	re := regexp.MustCompile(`#\s(.*?)\s(.*)`)
	match := re.FindSubmatch([]byte(line))

	if match == nil {
		return "", ""
	}

	resource := string(match[1])
	action := string(match[2])
	if !slices.Contains(defaultActionPattern, action) {
		return resource, action
	}

	return "", ""
}

var (
	defaultActionPattern = []string{"will be created", "will be destroyed", "will be updated in-place", "must be replaced"}
)

type Summary struct {
	Create  []string
	Destroy []string
	Update  []string
	Replace []string
}
