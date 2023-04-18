### mitum-sto

*mitum-sto* is a security token contract model based on [mitum](https://github.com/ProtoconNet/mitum).

#### Installation

```sh
$ git clone https://github.com/ProtoconNet/mitum-sto

$ cd mitum-sto

$ go build -o ./mitum-sto ./main.go
```

#### Run

```sh
$ ./mitum-sto init --design=<config file> <genesis config file>

$ ./mitum-sto run --design=<config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.
[genesis-design.yml](genesis-design.yml) is a sample of `genesis config file`.