# MCPing-cli
The goal of this software is to provide a UNIX ping like tool for Minecraft server administrators.

This software is an implementation of [my MCPing library](https://github.com/ammario/mcping)
## Install
```bash
git clone https://github.com/ammario/mcping-cli
cd mcping-cli
export GOPATH=`pwd`
go get ./...
go build
sudo cp mcping-cli /usr/bin/mcping
```


## Basic Usage
```bash
mcping -h vapormc.co:25565
(0) vapormc.co:25565; latency=76ms players=(78/80)
(1) vapormc.co:25565; latency=76ms players=(78/80)
(2) vapormc.co:25565; latency=76ms players=(77/80)
```
All responses are colored. If your count increases, the output will be green; if it drops, the output will be yellow; if it doesn't change, the output will be white.
