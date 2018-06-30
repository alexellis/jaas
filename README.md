# Jobs as a Service (JaaS)

Version 2.0

See propsal for more info: https://github.com/alexellis/jaas/issues/31

## Building:

```
mkdir -p ~/go/src/github.com/alexellis
cd  ~/go/src/github.com/alexellis

git clone https://github.com/alexellis/jaas
cd jaas

git checkout add_server

go build
```

## Server

Prepare credentials (stored in ~/.jaas/)

```
echo -n test | ./jaas server login -u=alex --password-stdin
```

Start server (uses stored credentials and basic auth):

```
./jaas server start
```

## Client

The client will send basic auth credentials

Run `hello-world` image:

```
./jaas server run --image hello-world
Running: hello-world 
Status: 200
```

List JaaS tasks:

```
./jaas server list

Listnaughty_tesla#vivved25azr6btmrr90ku9n6g 1 complete
hardcore_ardinghelli#z6qf7zfurhe5a7bptckhpcdnb 1 preparing
naughty_sammet#69hcmeembr5mdsog3r40e4ak9 1 complete
friendly_poitras#eb4275wcbkof4ee0ybz9p4xj6 1 complete
quirky_ritchie#qcs9xibtgo0udxv1hmxop1hy0 1 complete
distracted_brahmagupta#w87jqdloxe8vtd1xvjfskgwy8 1 complete
fervent_jackson#n5d5ipxhtpdcpcmwvl6txl4pw 1 complete
blissful_lamport#nh2q3oxaryvxsfw1drsbpw9da 1 complete
nervous_mayer#hbrb8a29zj9zpcm7mxntazm72 1 complete
modest_jones#t30cu5a79dy3h4d5vkvcizhs2 1 complete
```

