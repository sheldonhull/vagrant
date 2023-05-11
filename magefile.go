//go:build mage
// +build mage

package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	// "os/exec"

	// "github.com/bitfield/script"
	// "github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

// var Aliases = map[string]interface{}{
// 	"up": Vagrant.Up,
// }

// Vagrant runs commands that are interactive and magical
func Vagrant() error {
	r, _ := pterm.DefaultInteractiveSelect.
		WithOptions([]string{"ubuntu-kinetic", "windows10", "windows11"}).
		WithDefaultText("Select a VM to bring up").
		Show()
	pterm.Info.Printfln("Opening: %s", r)

	action, _ := pterm.DefaultInteractiveSelect.
		WithOptions([]string{"up", "suspend", "destroy", "provision"}).
		WithDefaultText("🪄 What magic should I perform?").
		Show()
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	pterm.Info.Printfln("original directory: %s", wd)
	if err := os.Chdir(r); err != nil {
		return err
	}
	pterm.Info.Printfln("changed directory to: %s", r)
	defer func() {
		os.Chdir(wd) // flip back to project directory
		pterm.Info.Printfln("changed directory back to working directory: %s", wd)
	}()
	// cmd := exec.Command("vagrant", "up", "--install-provider")
	// cmd.Dir = r
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	if err := sh.RunV("vagrant", action); err != nil {
		if runtime.GOOS == "darwin" && strings.Contains(err.Error(), "available VMware adapters") {
			pterm.Warning.Printfln("Try running:\n\n" +
				"sudo launchctl stop com.vagrant.vagrant-vmware-utility && sudo launchctl start com.vagrant.vagrant-vmware-utility")
			return err
		}

		return err
	}

	return nil
}

// Upgrade upgrades the vagrant plugins.
func Upgrade() error {
	return sh.RunV("vagrant", "plugin", "update", "vmware_desktop")
}

// Init sets up as much tooling as possible, such as vagrant plugins.
func Init() error {
	switch runtime.GOOS {
	case "darwin":
		if err := sh.RunV("brew", "install", "hashicorp-vagrant"); err != nil {
			pterm.Error.Printfln("can't install required plugin: %v")
			pterm.Error.Println("visit https://www.vagrantup.com/docs/providers/vmware/installation")
			return err
		}
		if err := sh.RunV("vagrant", "plugin", "install", "vagrant-vmware-desktop"); err != nil {
			pterm.Error.Printfln("can't install required plugin: %v")
			pterm.Error.Println("visit https://www.vagrantup.com/docs/providers/vmware/installation")
			return err
		}
		if err := sh.RunV("brew", "install", "vagrant-vmware-utility"); err != nil {
			pterm.Error.Printfln("can't install required utility: %v")
			pterm.Error.Println("visit https://developer.hashicorp.com/vagrant/docs/providers/vmware/vagrant-vmware-utility")
			return err
		}
	case "linux":
		pterm.Warning.Println("not automated")
	case "windows":
		pterm.Warning.Println("not automated")
	}
	return nil
}

// Release using github cli (for now)
func Release() error {
	version, changelogFile, err := getVersion()
	if err != nil {
		pterm.Error.Printfln("failed to get version: %v", err)
		return err
	}
	return sh.Run("gh", "release", "create", version, "--title", version, "--notes-file", changelogFile, "--target", "main")
}

// getVersion returns the version and path for the changefile to use for the semver and release notes.
func getVersion() (releaseVersion, cleanPath string, err error) {
	releaseVersion, err = sh.Output("changie", "latest")
	if err != nil {
		pterm.Error.Printfln("changie pulling latest release note version failure: %v", err)
		return "", "", err
	}
	cleanVersion := strings.TrimSpace(releaseVersion)
	cleanPath = filepath.Join(".changes", cleanVersion+".md")
	if os.Getenv("GITHUB_WORKSPACE") != "" {
		cleanPath = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".changes", cleanVersion+".md")
	}
	return cleanVersion, cleanPath, nil
}
