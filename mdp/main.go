package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
{{ .Body }}
	</body>
</html>
`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdownfile to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if os.Getenv("MDP_TEMPLATE") != "" {
		*tFname = os.Getenv("MDP_TEMPLATE")
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func preview(fname string) error {
	cName := ""
	cParams := []string{}
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	if err := exec.Command(cPath, cParams...).Run(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func run(filename string, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

	temp, err := os.CreateTemp("", "mdp-*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)
	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte, tFname string) (io.Reader, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var t *template.Template
	var err error

	if tFname != "" {
		if t, err = template.ParseFiles(tFname); err != nil {
			return nil, err
		}
	} else {
		if t, err = template.New("mdp").Parse(defaultTemplate); err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown preview tool",
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}
	return &buffer, nil
}

func saveHTML(outFname string, data io.Reader) error {
	bytes, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	return os.WriteFile(outFname, bytes, 0644)
}
