package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/magaluCloud/mgccli/cmd"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
)

var RawVersion string

var version string = func() string {
	if RawVersion == "" {
		return getVCSInfo("v0.0.0")
	}

	return strings.Trim(RawVersion, " \t\n\r")
}()

func getVCSInfo(version string) string {
	if info, ok := debug.ReadBuildInfo(); ok {
		var vcs, rev, status string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs":
				vcs = setting.Value
			case "vcs.revision":
				rev = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					status = " (modified)"
				}
			}
		}

		if vcs != "" {
			return fmt.Sprintf("%s %s%s", version, rev, status)
		}
	}
	return "v0.0.0 dev"
}

func getLang() string {
	lang, isPresent := os.LookupEnv("CLI_LANG")
	if isPresent {
		return lang
	}
	return ""
}

func main() {
	args := cmdutils.NewArgsParser()

	panicOff := os.Getenv("CLI_PANIC_OFF")
	if panicOff == "" {
		defer panicRecover()
	}
	ctx := context.Background()

	lang := getLang()
	if lang == "" {
		langFlagValue, _, _ := args.GetValue(cmd.LangFlag)
		lang = langFlagValue
	}

	manager := i18n.Init18n(lang)
	version := fmt.Sprintf("%s (%s)", version, manager.GetLanguage())

	rootCmd := cmd.RootCmd(ctx, version, args)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func panicRecover() {
	err := recover()
	if err != nil {

		Url := "https://github.com/magaluCloud/mgccli/issues/new"
		args := strings.Join(os.Args, " ")

		query := url.Values{}
		query.Add("title", fmt.Sprintf("Error report at '%s'", args))
		query.Add("body", fmt.Sprintf("Version: %s\nSO: %s / %s\nArgs: %s\nError: %s\n",
			version,
			runtime.GOOS,
			runtime.GOARCH,
			args,
			err))
		Url = Url + "?" + query.Encode()

		manager := i18n.GetInstance()

		fmt.Fprintf(os.Stderr, `
😔 %s
     %s: %s
     %s: %s / %s  
     %s: %s 
     %s: %s

%s
	%s

%s
`,
			manager.T("cli.panic_message"),
			manager.T("cli.version"),
			version,
			manager.T("cli.os"),
			runtime.GOOS,
			runtime.GOARCH,
			manager.T("cli.args"),
			args,
			manager.T("cli.error"),
			err,
			manager.T("cli.panic_help"),
			Url,
			manager.T("cli.panic_thanks"))
		os.Exit(1)
	}
}
