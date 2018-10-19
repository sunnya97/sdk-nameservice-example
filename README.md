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

Start by installing Dep.

```
go get -v github.com/golang/dep/cmd/dep
```

Next, run

```
dep ensure
```

Finally run

```
make install
```

### Project directory structure

Through the course of this tutorial you will create the following files that make up your application:

```bash
./nameservice
├── Gopkg.toml
├── Makefile
├── app.go
├── cmd
│   ├── nameshakecli
│   │   └── main.go
│   └── nameshaked
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
