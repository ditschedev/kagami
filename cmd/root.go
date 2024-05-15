package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var (
	versionFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "kagami",
	Short: "A tiny cli tool that help mirror git repositories",
	Long:  `Simply provide a git repository and kagami will mirror it for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "shows the version of the cli")
}

func Execute() {
	printBanner()

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute command")
		os.Exit(1)
	}
}

func printBanner() {
	color.Set(color.FgHiCyan)
	fmt.Println("   __                          ")
	fmt.Println("  / /_____ ____ ____ ___ _  (_)")
	fmt.Println(" /  '_/ _ `/ _ `/ _ `/  ' \\/ / ")
	fmt.Println("/_/\\_\\\\_,_/\\_, /\\_,_/_/_/_/_/  ")
	fmt.Println("          /___/                ")
	fmt.Println()
	color.Unset()
}
