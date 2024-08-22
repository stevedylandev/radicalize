# Radicalize

![radicle banner](https://dweb.mypinata.cloud/ipfs/QmUFwBiweWHtGBxftQ7xNpiS5xSBHJyZJgsHXXGRy2qyLH)

A script that helps clone existing git repos to [Radicle](https://radicle.xyz).

## Pre-Installation

In order to use this script you will already need Radicle installed on your machine. Follow the instructions [here](https://radicle.xyz/#get-started) before moving forward.

After installing be sure to initialize Radicle with the following command:

```
rad auth
```

Then spin up the node to make sure that is working as well.

```
rad node start
```

## Installation

This script is written in Go so there are a few install options.

### Install with Homebrew
```
brew install stevedylandev/radicalize/radicalize
```

### Install directly with Go
Have Go installed and run this command.
```
go install github.com/stevedylandev/radicalize@latest
```

### Clone and Build with Go
For this method simply clone the repo, build, and install with Go
```
git clone https://github.com/stevedylandev/radicalize
cd radicalize
go build
go install .
```

## Usage

This CLI (`radicalize`) has two commands:

### `local`

To start backing up your local repos simply run the command `radicalize local` in the parent directory of all your projects. This will creep through all directories for .git repos on a surface level, so if you have sub directories you will want to navigate to those separately. After finding all git repos you can select which ones you would like to init to your Radicle node. Once you have confirmed the selection it will work through each repo and initialize it.

![radicalize local gif](https://dweb.mypinata.cloud/ipfs/QmaPFKdTuS7aauJ5RpiZssJj82RWuLqypAow7MfMZ5Nzkp)

**Private Repos**

By default the program will make repos public, however you can pass in the `--private` flag so all selected repos will be private instead

```
radicalize local --private
```

### `remote`

The CLI can also pull the latest 100 repos from any given user or organization on GitHub, clone them, and initalize them locally as well. Run `radicalize remote` which will then prompt you for the username or org name to search. Once provided it will put the repos into a selection list where you can start typing the name as well as scroll through it. Once all are selected it will begin the process of cloning and initalizing.

![radicalize remote gif](https://dweb.mypinata.cloud/ipfs/QmdErhQsJAshuTqaVPjPZKaxBDdwmzuB4JsVUw7zpRVbQi)
