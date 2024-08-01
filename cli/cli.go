package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/strutil"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"

	"github.com/essentialkaos/detheme/converter"
	"github.com/essentialkaos/detheme/theme"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic utility info
const (
	APP  = "detheme"
	VER  = "0.0.4"
	DESC = "SublimeText color theme downgrader (sublime-color-scheme → tmTheme converter)"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_OUTPUT   = "o:output"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_OUTPUT:   {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	err := convert(args)

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	fmtc.DisableColors = true

	if fmtc.IsColorsSupported() {
		fmtc.DisableColors = false
	}

	if os.Getenv("NO_COLOR") != "" {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// convert starts converting process
func convert(args options.Arguments) error {
	themeFile := args.Get(0).Clean().String()
	stTheme, err := theme.Read(themeFile)

	if err != nil {
		return fmt.Errorf("Can't load theme: %v", err)
	}

	printThemeInfo(stTheme)

	data, err := converter.Convert(stTheme)

	if err != nil {
		return fmt.Errorf("Can't convert theme: %v", err)
	}

	outputFile := strutil.Q(
		options.GetS(OPT_OUTPUT),
		strings.ReplaceAll(themeFile, ".sublime-color-scheme", ".tmTheme"),
	)

	err = os.WriteFile(outputFile, data, 0644)

	if err != nil {
		return fmt.Errorf("Can't save theme: %v", err)
	}

	fmtc.Printf("{g}Theme successfully converted and saved as {*}%s{!*}\n\n", outputFile)

	return nil
}

// printThemeInfo prints info about theme
func printThemeInfo(t *theme.Theme) {
	fmtutil.Separator(false, strutil.Q(t.Name, "Theme"))

	fmtc.Printf("  {*}Author:{!}    %s\n", strutil.Q(t.Author, "—"))
	fmtc.Printf("  {*}Variables:{!} %s\n", fmtutil.PrettyNum(len(t.Variables)))
	fmtc.Printf("  {*}Globals:{!}   %s\n", fmtutil.PrettyNum(len(t.Globals.Data)))
	fmtc.Printf("  {*}Rules:{!}     %s\n", fmtutil.PrettyNum(len(t.Rules)))

	fmtutil.Separator(false)
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	if len(a) == 0 {
		fmtc.Fprintln(os.Stderr, "{r}"+f+"{!}")
	} else {
		fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, "detheme"))
	case "fish":
		fmt.Print(fish.Generate(info, "detheme"))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, "detheme"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "theme-file")

	info.AddOption(OPT_OUTPUT, "Path to output file", "path")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"my-theme.sublime-color-scheme",
		"Convert custom theme to thTheme format",
	)

	info.AddExample(
		"-o theme1.thTheme my-theme.sublime-color-scheme",
		"Convert custom theme to thTheme format and save as theme1.thTheme",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
		about.UpdateChecker = usage.UpdateChecker{
			"essentialkaos/detheme",
			update.GitHubChecker,
		}
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
