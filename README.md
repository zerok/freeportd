# FreeportD

**Warning:** This is right now just a toy-project. Use at your own risk!

This project aims at providing a simple service that returns the next
available TCP port on the local machine. The idea for this came when I
ran into a situation where multiple test-runners on a Jenkins server
attempted to bind the same port. Now they can reserve a port for their
exclusive use using a simple HTTP GET request.

Additionally, the server keeps an internal cache in order to reserve
ports for a given grace-period (`--grace-period`). This was necessary as
the requesting CI jobs might take a while to actually bind their ports.

## Usage

```
Usage of ./freeportd:
--cache-size int          Maximum number of ports in the cache (default 1024)
--grace-period duration   Graceperiod for how long ports are kept in the internal store (default 5m0s)
--http-addr string        Address to listen on for HTTP requests (default "localhost:8888")
```


## Inspiration

This project was heavily inspired by https://github.com/phayes/freeport. Big
thanks to [Patrick D Hayes](https://github.com/phayes) for this project!
