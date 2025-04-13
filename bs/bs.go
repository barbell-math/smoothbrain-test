package main

import (
	"bytes"
	"context"
	"strings"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

// Generate a stage that runs `git diff` and returns an error if there are any
// differences. This is mainly used by ci targets.
func gitDiffStage(errMessage string, targetToRun string) sbbs.StageFunc {
	return sbbs.Stage(
		"Run Diff",
		func(ctxt context.Context, cmdLineArgs ...string) error {
			var buf bytes.Buffer
			if err := sbbs.Run(ctxt, &buf, "git", "diff"); err != nil {
				return err
			}
			if buf.Len() > 0 {
				sbbs.LogErr(errMessage)
				sbbs.LogQuietInfo(buf.String())
				sbbs.LogErr(
					"Run build system with %s and push any changes",
					targetToRun,
				)
				return sbbs.StopErr
			}
			return nil
		},
	)
}

func main() {
	// Register a target that updates all dependences. Dependencies that are in
	// the `barbell-math` repo will always be pinned at latest and all other
	// dependencies will be updated to the latest version.
	sbbs.RegisterTarget(
		context.Background(),
		"updateDeps",
		sbbs.Stage(
			"barbell math package cmds",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				var packages bytes.Buffer
				if err := sbbs.Run(
					ctxt, &packages, "go", "list", "-m", "-u", "all",
				); err != nil {
					return err
				}

				lines := strings.Split(packages.String(), "\n")
				// First line is the current package, skip it
				for i := 1; i < len(lines); i++ {
					iterPackage := strings.SplitN(lines[i], " ", 2)
					if !strings.Contains(iterPackage[0], "barbell-math") {
						continue
					}

					if err := sbbs.RunStdout(
						ctxt, "go", "get", iterPackage[0]+"@latest",
					); err != nil {
						return err
					}
				}
				return nil
			},
		),
		sbbs.Stage(
			"Non barbell math package cmds",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := sbbs.RunStdout(ctxt, "go", "get", "-u", "./..."); err != nil {
					return err
				}
				if err := sbbs.RunStdout(ctxt, "go", "mod", "tidy"); err != nil {
					return err
				}

				return nil
			},
		),
	)
	sbbs.RegisterTarget(
		context.Background(),
		"updateReadme",
		sbbs.Stage(
			"Run gomarkdoc",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				err := sbbs.RunStdout(
					ctxt, "gomarkdoc", "--embed", "--output", "README.md", ".",
				)
				if err != nil {
					sbbs.LogQuietInfo("Consider running build system with installGoMarkDoc target if gomarkdoc is not installed")
				}
				return err
			},
		),
	)
	sbbs.RegisterTarget(
		context.Background(),
		"installGoMarkDoc",
		sbbs.Stage(
			"Install gomarkdoc",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(
					ctxt, "go",
					"install", "github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest",
				)
			},
		),
	)
	sbbs.RegisterTarget(
		context.Background(),
		"fmt",
		sbbs.Stage(
			"Run go fmt",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "go", "fmt", "./...")
			},
		),
	)
	sbbs.RegisterTarget(
		context.Background(),
		"unitTests",
		sbbs.Stage(
			"Run go test",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "go", "test", "-v", "./...")
			},
		),
	)
	sbbs.RegisterTarget(
		context.Background(),
		"buildBs",
		sbbs.Stage(
			"Run go build",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "go", "build", "-o", "./bs/bs", "./bs")
			},
		),
	)

	// Registers a target that will update all deps and run a diff to make sure
	// that the commited code is using all of the correct dependencies.
	sbbs.RegisterTarget(
		context.Background(),
		"ciCheckDeps",
		sbbs.TargetAsStage("updateDeps"),
		gitDiffStage("Out of date packages were detected", "updateDeps"),
	)
	// Registers a target that will install gomarkdoc, update the readme, and
	// run a diff to make sure that the commited readme is up to date.
	sbbs.RegisterTarget(
		context.Background(),
		"ciCheckReadme",
		sbbs.TargetAsStage("installGoMarkDoc"),
		sbbs.TargetAsStage("updateReadme"),
		gitDiffStage("Readme is out of date", "updateReadme"),
	)
	// Registers a target that will run go fmt and then run a diff to make sure
	// that the commited code is properly formated.
	sbbs.RegisterTarget(
		context.Background(),
		"ciCheckFmt",
		sbbs.TargetAsStage("fmt"),
		gitDiffStage("Fix formatting to get a passing run!", "fmt"),
	)
	// Registers a target that will run all mergegate checks. This includes:
	//	- checking that the code is formatted
	//	- checking that the readme is up to date
	//	- checking that the dependencies are up to date
	//	- checking that all unit tests pass
	sbbs.RegisterTarget(
		context.Background(),
		"mergegate",
		sbbs.TargetAsStage("ciCheckFmt"),
		sbbs.TargetAsStage("ciCheckReadme"),
		sbbs.TargetAsStage("ciCheckDeps"),
		sbbs.TargetAsStage("unitTests"),
	)

	sbbs.Main("build")
}
