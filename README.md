# Example SDK Module Tutorial

In this tutorial, you will build a functional [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) module, and in the process, learn the basic concepts and structures in the SDK. Then you can get started building your own modules and include them in decentralized applications. By the end of this tutorial you will have a functional `nameservice`, a mapping of strings to other strings (`map[string]string`). This is similar to [Namecoin](https://namecoin.org/), [ENS](https://ens.domains/), or [Handshake](https://handshake.org/), which all model the traditional DNS systems (`map[domain]zonefile`). Users will be able to buy unused names, or sell/trade their name.

All of the final source code for this tutorial project is in this directory (and compiles), however, it is best to follow along manually and try building the project yourself!

### Requirements:

- `golang` >1.11 installed
- A working `$GOPATH`
- Desire to create your own blockchain!

### Tutorial Parts:

1. Start by building your [`Keeper`](./tutorial/keeper.md)
2. Define interactions with your chain through [`Msgs` and `Handlers`](./tutorial/msgs-handlers.md)
	* [`SetName`](./tutorial/set-name.md)
	* [`BuyName`](./tutorial/buy-name.md)
3. Make views on your state machine with [`Queriers`](./tutorial/queriers.md)
4. Register your types in the encoding format using [`sdk.Codec`](./tutorial/codec.md)
5. Create [CLI interactions for your module](./tutorial/cli.md)
6. Put it all together in [`./app.go`](./tutorial/app.md)!
7. Create the [`nameserviced` and `nameservicecli` entry points](./tutorial/entrypoint.md)
8. Setup [dependency management using `dep`](./tutorial/dep.md)


### Build the `nameservice` application!

If you want to build the `nameservice` application in this repo to see the functionality, first you need to install `dep`. Below there is a command for using a shell script from `dep`'s site to preform this install. If you are uncomfortable `|`ing `curl` output to `sh` (you should be) then check out [your platform specific installation instructions](https://golang.github.io/dep/docs/installation.html).

```bash
# Install dep
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Initialize dep and install dependencies
dep init
dep ensure -v -upgrade

# Install the app into your $GOBIN
go install -v ./cmd/...

# Now you should be able to run the following commands:
nameserviced help
nameservicecli help
```

### Using `nameservicecli`

TODO: Write a list of the interactions that are enabled with the `nameservice` module and commands that are enabled by them

### Project directory structure

Through the course of this tutorial you will create the following files that make up your application:

```bash
./nameservice
├── Gopkg.toml
├── app.go
├── cmd
│   ├── nameservicecli
│   │   └── main.go
│   └── nameserviced
│       └── main.go
└── x
    └── nameservice
        ├── client
        │   └── cli
        │       ├── query.go
        │       └── tx.go
        ├── codec.go
        ├── handler.go
        ├── keeper.go
        ├── msgs.go
        └── querier.go
```
