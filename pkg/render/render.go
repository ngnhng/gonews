// This package includes functions used in rendering HTML
package render

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	TIMEZONE string = "Asia/Saigon"
)

type Renderer struct {
	Tmpl    *template.Template
	OutPath string
	Feeds   *map[string]*gofeed.Feed
}

type page struct {
	Sources     []string
	Title       string
	Items       []*gofeed.Item
	Description string
}

func Time() func() time.Time {
	currentTimeZone, _ := time.LoadLocation(TIMEZONE)
	return func() time.Time {
		return time.Now().In(currentTimeZone)
	}
}

func Trim() func(string) string {
	return strings.TrimSpace
}

func Simplify() func(string) string {
	return func(s string) string {
		return strings.Map(func(r rune) rune {
			if r >= 'a' && r <= 'z' {
				return r
			}
			return -1
		}, strings.ToLower(s))
	}
}

func UnescapeHTML() func(string) template.HTML {
	return func(s string) template.HTML {
		return template.HTML(s)
	}
}

func FormatLocalTime() func(string) string {
	// format: Mon, 02 Jan 2006 15:04:05 +0000
	return func(t string) string {
		tm, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", t)
		if err != nil {
			return t
		}
		currentTimeZone, _ := time.LoadLocation(TIMEZONE)
		return tm.In(currentTimeZone).Format("Mon, 02 Jan 2006 15:04")
	}
}

// RenderPages renders HTML pages for each source
func (rder *Renderer) RenderPages() error {
	var sources []string
	for name := range *rder.Feeds {
		sources = append(sources, name)
	} // bad design :v
	for name, feed := range *rder.Feeds {
		p := page{
			Sources:     sources,
			Title:       feed.Title,
			Items:       feed.Items,
			Description: feed.Description,
		}
		// check if directory exists
		if _, err := os.Stat(rder.OutPath); os.IsNotExist(err) {
			if e := os.Mkdir(rder.OutPath, 0755); e != nil {
				return fmt.Errorf("render: Error creating directory; %w", e)
			}
		}
		// create file
		outPath := path.Join(rder.OutPath, Simplify()(name)+".html")
		file, err := os.Create(outPath)

		if err != nil {
			return fmt.Errorf("render: Error creating file; %w", err)
		}
		defer file.Close()

		err = rder.Tmpl.Execute(file, p)
		if err != nil {
			return fmt.Errorf("render: Error executing template; %w", err)
		}
	}
	return nil
}

func (rder *Renderer) RenderIndex() error {
	var sources []string
	for name := range *rder.Feeds {
		sources = append(sources, name)
	} // bad design :v
	p := page{
		Sources:     sources,
		Title:       "Main page",
		Description: "Select a source to read",
	}
	// check if directory exists
	if _, err := os.Stat(rder.OutPath); os.IsNotExist(err) {
		if e := os.Mkdir(rder.OutPath, 0755); e != nil {
			return fmt.Errorf("render: Error creating directory; %w", e)
		}
	}
	// create file
	outPath := path.Join(rder.OutPath, "index.html")
	file, err := os.Create(outPath)

	if err != nil {
		return fmt.Errorf("render: Error creating file; %w", err)
	}
	defer file.Close()

	err = rder.Tmpl.Execute(file, p)
	if err != nil {
		return fmt.Errorf("render: Error executing template; %w", err)
	}
	return nil
}
