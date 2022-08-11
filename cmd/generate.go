/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"html/template"
	"log"
	"os"
	"path"
	"sync"

	"github.com/mmcdole/gofeed"
	"github.com/nguyendhst/gonews/pkg/fetch"
	"github.com/nguyendhst/gonews/pkg/render"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	configPath   string
	outputDir    string
	templatePath string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate static HTML files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read config file
		configFile, err := os.Open(configPath)
		if err != nil {
			log.Println("Error opening config file: ;%w", err)
			os.Exit(1)
		}
		defer configFile.Close()

		// Parse config file
		var newsSources = make(map[string]fetch.NewsSource)

		err = yaml.NewDecoder(configFile).Decode(&newsSources)
		if err != nil {
			log.Println("Error decoding config file: ;%w", err)
			os.Exit(1)
		}

		tmpl, err := template.
			New(path.Base(templatePath)).
			Funcs(template.FuncMap{
				"time":         render.Time(),
				"trim":         render.Trim(),
				"simplify":     render.Simplify(),
				"unescapeHTML": render.UnescapeHTML(),
			}).ParseFiles(templatePath)

		if err != nil {
			log.Printf("[Error] failed parsing template: %s; %v\n", templatePath, err)
		}

		// Fetch and parse XML to gofeed.Feed
		// responded := make(chan *fetch.NewsFeed, len(newsSources))
		// defer close(responded)

		wg := sync.WaitGroup{}

		// Concurrent fetch with timeout
		parsedFeed := make(map[string]*gofeed.Feed)
		for _, source := range newsSources {
			wg.Add(1)
			go func(source fetch.NewsSource) {
				// Starts fetching, some sources may respond with 403 (Forbidden) due to strict bot policy.
				var fetchErr error
				feed, fetchErr := source.Fetch()
				if fetchErr != nil {
					log.Printf("[Error] failed fetching source: %s; %v\n", source.URL, err)
				}

				parsedFeed[feed.Source.Name] = feed.Feed
				wg.Done()
			}(source)
		}

		// wait for all goroutines to finish
		wg.Wait()

		if len(parsedFeed) == 0 {
			log.Println("[Error] no feed parsed")
			os.Exit(1)
		}

		// Generate HTML files
		renderer := render.Renderer{
			Tmpl:    tmpl,
			OutPath: outputDir,
			Feeds:   &parsedFeed,
		}
		err = renderer.RenderPages()
		if err != nil {
			log.Printf("[Error] failed rendering: %v\n", err)
			os.Exit(1)
		}

		// index main page
		err = renderer.RenderIndex()
		if err != nil {
			log.Printf("[Error] failed rendering index: %v\n", err)
			os.Exit(1)
		}

		log.Println("Successfully generated HTML files")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "./static", "Output directory")
	generateCmd.Flags().StringVarP(&configPath, "config", "c", "config.yml", "Config path")
	generateCmd.Flags().StringVarP(&templatePath, "template", "t", "./resource/template.html", "Template path")
}
