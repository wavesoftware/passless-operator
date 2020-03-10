// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-ensure"
)

const (
	imageName           = "quay.io/wavesoftware/passless-operator"
	versionVariablePath = "github.com/wavesoftware/passless-operator/version.Version"
)

// Check will run all lints checks
func Check() {
	t := startMultilineTask("🔍", "Checking")
	mg.Deps(revive, staticcheck)
	t.end(nil)
}

func revive() error {
	mg.Deps(buildDeps)
	return sh.RunV("revive", "-config", "revive.toml", "-formatter", "stylish", "./...")
}

func staticcheck() error {
	mg.Deps(buildDeps)
	return sh.RunV("staticcheck", "-f", "stylish", "./...")
}

// Test will execute regular unit tests
func Test() {
	mg.Deps(Check, ensureBuildDir)
	t := startMultilineTask("✅", "Testing")
	cmd := "richgo"
	if color.NoColor {
		cmd = "go"
	}
	err := sh.RunV(cmd, "test", "-v",
		"-covermode=count",
		fmt.Sprintf("-coverprofile=%s/coverage.out", buildDir),
		fmt.Sprintf("-ldflags=-X %s=%s", versionVariablePath, gitVersion()),
		"./...",
	)
	t.end(err)
}

// Binary will build a binary executable file
func Binary() {
	mg.Deps(Test, ensureBuildDir)
	t := startTask("🔨", "Building")
	err := sh.RunV("go", "build",
		fmt.Sprintf("-ldflags=-X %s=%s", versionVariablePath, gitVersion()),
		"-o", managerBin, managerDir,
	)
	t.end(err)
}

// ImageSole will try to build an image only. It will not ensure that
// the project is tested beforehand
func ImageSole() {
	t := startMultilineTask("💿", "Creating an image")
	err := sh.RunV("operator-sdk", "build",
		"--go-build-args",
		fmt.Sprintf("-ldflags -X=%s=%s", versionVariablePath, gitVersion()),
		imageName,
	)
	t.end(err)
}

// Image will create a container image checking project first
func Image() {
	mg.SerialDeps(Test, ImageSole)
}

// A var named Default indicates which target is the default.  If there is no
// default, running mage will list the targets available.
var Default = Binary

// Clean will clean build directories
func Clean() {
	t := startTask("🚿", "Cleaning")
	err1 := os.RemoveAll(buildDir)
	imageHash, err2 := sh.Output("docker", "images", "-q", imageName)
	var err3 error
	if len(imageHash) > 0 {
		err3 = sh.Run("docker", "rmi", "-f", imageName)
	}
	t.end(err1, err2, err3)
}

// buildDeps install build dependencies
func buildDeps() error {
	for _, dep := range []string{
		"github.com/kyoh86/richgo",
		"github.com/mgechev/revive",
		"honnef.co/go/tools/cmd/staticcheck",
	} {
		err := sh.RunWith(map[string]string{"GO111MODULE": "off"}, "go", "get", dep)
		if err != nil {
			return err
		}
	}
	return nil
}

// ensureBuildDir creates a build directory
func ensureBuildDir() {
	d := path.Join(buildDir, "bin")
	ensure.NoError(os.MkdirAll(d, os.ModePerm))
}

var (
	repoDir             = currentFileDir()
	buildDir            = path.Join(repoDir, "build", "_output")
	managerDir          = path.Join(repoDir, "cmd", "manager")
	managerBin          = path.Join(buildDir, "bin", "passless-operator")
	git                 = sh.OutCmd("git")
	gitVerCache *string = nil
	mageTag             = color.New(color.FgCyan).Sprint("[MAGE]")
)

func gitVersion() string {
	if gitVerCache == nil {
		ver, err := git("describe", "--always", "--tags", "--dirty")
		ensure.NoError(err)
		gitVerCache = &ver
	}
	return *gitVerCache
}

func currentFileDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

type task struct {
	icon      string
	action    string
	multiline bool
}

func startTask(icon, action string) *task {
	t := &task{
		icon:      icon,
		action:    action,
		multiline: false,
	}
	t.start()
	return t
}

func startMultilineTask(icon, action string) *task {
	t := &task{
		icon:      icon,
		action:    action,
		multiline: true,
	}
	t.start()
	return t
}

func (t *task) start() {
	if t.multiline {
		fmt.Printf("%s %s %s\n", mageTag, t.icon, t.action)
	} else {
		fmt.Printf("%s %s %s... ", mageTag, t.icon, t.action)
	}
}

func (t *task) end(errs ...error) {
	var msg string
	merr := multierror.Append(nil, errs...)
	err := merr.ErrorOrNil()
	green := color.New(color.FgHiGreen).Add(color.Bold).SprintFunc()
	red := color.New(color.FgHiRed).Add(color.Bold).SprintFunc()
	if err != nil {
		if t.multiline {
			msg = mageTag + red(fmt.Sprintf(" %s have failed!\n", t.action))
		} else {
			msg = red(fmt.Sprintln("failed!"))
		}
	} else {
		if t.multiline {
			msg = mageTag + green(fmt.Sprintf(" %s was successful.\n", t.action))
		} else {
			msg = green(fmt.Sprintln("done."))
		}
	}

	fmt.Print(msg)
	ensure.NoError(err)
}
