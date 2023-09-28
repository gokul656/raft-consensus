## main concepts

- [ ] leader
- [ ] follow
- [ ] client

## admin

- [ ] add peer
- [ ] remove peer
- [ ] make peer to vote
- [ ] change leadership

## todo

- [ ] create two cli for admin and peers
- [ ] manage Log replication
- [ ] add gRPC implementation for peers
- [ ] add health check functionalities

## learnings

* go work init
* always keep recover function in top
* generic for function ref : `FromByte`

* setting up gRPC server
    * setup struct for grpc server
    * implementing methods in proto file is mandatory, method should have 
    `connectionServer` as function param and it should return `error`
    * implement a `Connection` service and listen to it infinitely using `for {}`
    