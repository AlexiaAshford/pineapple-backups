package cfg

import (
	"fmt"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

var (
	Book_id      string
	Account      string
	Password     string
	App_type     string
	Search_key   string
	Max_thread   int
	show_info    bool
	up_date      bool
	command_line []string
)

func Console() ([]string, bool) {
	spaceRe, _ := regexp.Compile(`\s+`)
	inputs := spaceRe.Split(strings.TrimSpace(Input(">")), -1)
	if len(inputs) > 0 && inputs[0] != "" {
		return inputs, true
	} else {
		return nil, false
	}
}

var ruleCmd = &cobra.Command{
	Use:   "https://github.com/VeronicaAlexia/pineapple-backups",
	Short: "you can use this command tools to backup your data",
	Long:  "[warning] you login required to use this command tools",
}

func CommandInit() []string {
	ruleCmd.Flags().StringVarP(&Book_id, "download", "d", "", "")
	ruleCmd.Flags().StringVar(&Account, "account", "", "input account")
	ruleCmd.Flags().StringVar(&Password, "password", "", "input password")
	ruleCmd.Flags().StringVarP(&App_type, "app", "a", "sfacg", "input app type")
	ruleCmd.Flags().StringVarP(&Search_key, "search", "s", "", "input search keyword")
	ruleCmd.Flags().IntVarP(&Max_thread, "max", "m", 32, "change max thread number")
	ruleCmd.Flags().BoolVar(&show_info, "show", false, "show config")
	ruleCmd.Flags().BoolVar(&up_date, "update", false, "update config")
	if err := ruleCmd.Execute(); err != nil {
		fmt.Println("ruleCmd.Execute:", err)
	}
	if show_info {
		command_line = []string{"show", "config"}
	}
	if Book_id != "" {
		command_line = []string{"download", Book_id}
	}
	if Search_key != "" {
		command_line = []string{"search", Search_key}
	}
	return command_line
}
