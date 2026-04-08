# go-acl-proxy
a simple content-filtering proxy written in go.

## how to use
provide the URL to your access control list(ACL) in the config.json
then run the codebase using 
    $ go run *.go

## NOTE
- considering the length of your ACL, your ACL is loaded concurrently in a goroutine, to minimize server startup time
  > therefore, access control rules later in you list might take a while to load
- you might have to install dependencies found in the go.mod or go.sum file using go get
- work still in progress

### ACLs: 
- https://easylist.to/easylist/easylist.txt
- https://big.oisd.nl/

### TODO / FEATURES TO BE IMPLEMENTED:
- adding additional rules on runtime
- imrove performance >> hardcode security header middleware
- dockerfile for convenient setup
- healthcheck
- ...
