# Vagrant

## Install Mage

Use [aqua](https://aquaproj.github.io/docs/install) or if comfortable with Go development install manually `go install github.com/magefile/mage@latest` (requires `GOPATH/bin` to be set in your environment variables).

Run `mage` to see available tasks to spin up, provision, or tear down a VM.

## Snippets

For virtual box guest additions.

```shell
vagrant plugin install vagrant-vbguest
```

## Run Windows 10 Defaults

```shell
vagrant init gusztavvargadr/windows-10 && vagrant up
```

## Ubuntu Install Vagrant

```shell
 wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install vagrant
```
