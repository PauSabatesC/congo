<p align="center" >
  <img src="https://github.com/PauSabatesC/congo/blob/main/logo-img.png" alt="congo logo img" width="250"/>
</p>

<h1 align="center">
  Easy and unified way to connect from your terminal to AWS EC2 and ECS
</h1>

<p align="center" >
🚀 Connect easily to EC2 and ECS using the same tool using an interactive CLI
</p>

<p align="center" >
<img src="https://github.com/PauSabatesC/congo/blob/main/cli-example.png" alt="congo cli example" width="600"/>
</p>

<p align="center" >
🚀 Filter to find faster your resources
</p>

<p align="center" >
<img src="https://github.com/PauSabatesC/congo/blob/main/cli-filter-example.png" alt="congo cli filter" width="600"/>
</p>

<p align="center" >
🚀 Perfect for small or big AWS accounts
</p>

<p align="center" >
<img src="https://github.com/PauSabatesC/congo/blob/main/cli-filter-thousand-example.png" alt="congo cli thousand" width="600"/>
</p>

<br>

<p align="center" >
  <img alt="Go report card" src="https://goreportcard.com/badge/github.com/PauSabatesC/congo">
  <img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/PauSabatesC/congo">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/PauSabatesC/congo">
  <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/PauSabatesC/congo">
</p>


--- 

## Usage

Congo is really easy to use just run "ec2" or "ecs" and `the cli will continue guiding you interactively` from there:

```sh
congo ec2 [--id]
congo ecs
```

You can always run the command `congo help` to get a better understanding of each command.


## Installation

- Using brew:

```sh
brew tap PauSabatesC/congo https://github.com/PauSabatesC/congo
brew install congo
```

- Install using go:

```sh
go install github.com/PauSabatesC/congo
```

- Downloading the binaries:

You can install and run it just downloading your desired binary:

| platform     |
| ----------- | 
| [macOS ARM](https://github.com/PauSabatesC/congo/releases)
| [macOS 64 Bit](https://github.com/PauSabatesC/congo/releases)
| [Linux 32-Bit](https://github.com/PauSabatesC/congo/releases)
| [Linux ARM](https://github.com/PauSabatesC/congo/releases)
| [Linux 64 Bit](https://github.com/PauSabatesC/congo/releases)
| [Windows ARM](https://github.com/PauSabatesC/congo/releases)
| [Windows 32 Bit](https://github.com/PauSabatesC/congo/releases)
| [Windows 64 Bit](https://github.com/PauSabatesC/congo/releases)

Then you can add the executable binary file downloaded into your PATH

- Install the golang package to use it on your project:

```sh
go get github.com/PauSabatesC/congo/congo
```

---

## Prerequisites

- Congo uses AWS SSM to start a session in a EC2 instance in order to connect to it. So make sure your EC2 and IAM role passes the requisites to allow this secure connection.
- Also congo uses aws CLI to exec to a container.
- Congo detects automatically your AWS credentials from environment variables. So even if you exported the AWS key id,secret and token, or you exported you AWS_PROFILE it will work. Make sure though to export the desired AWS_REGION too.

---

## Contributing

All kinds of Pull Requests are welcomed!

## Credits

- [Cobra](https://github.com/spf13/cobra)
- [Charm-Bubbletea](https://github.com/charmbracelet/bubbletea)
- [tedsmitt/ecsgo](https://github.com/tedsmitt/ecsgo)
- [gjbae1212/gossm](https://github.com/gjbae1212/gossm)


## License

Congo is released under the Apache 2.0 license. See [`LICENSE.txt`](https://github.com/PauSabatesC/congo/blob/master/LICENSE.txt) for more details.

