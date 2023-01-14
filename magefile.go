//go:build mage
// +build mage

package main

import (
	"os"
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
		WithOptions([]string{"ubuntu-kinetic", "windows10"}).
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
	return sh.RunV("vagrant", "plugin", "install", "vagrant-vmware-desktop")
}
