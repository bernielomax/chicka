# chicka
A pluggable Go monitoring system

### Building

To build Chicka you need a working Go environment.

##### OSx:

``` 
make osx
```

##### Linux:

``` 
make linux
```

Make will place a working `chicka` executable inside the Chicka source folder.

### Configuration

The configuration for Chicka can belong in the following locations:

* $HOME/.chicka/config.yaml
* /etc/chicka/config.yaml
* ./config.yaml

The following is an example of the config.yaml

```
git:
  url: git@github.com:bernielomax/chicka-plugins.git # git plugin repo URL
  pull: true # pulls latest on startup
plugins:
  path: /Users/justin.pye/Documents/chicka-plugins/ # local path where plugins will reside
logging:
  path: /tmp/chicka.log # local path where the logs will reside
cache:
  ttl: 120 # amount of time to store tests in cache
http:
  api_addr: :9090 # the listen address for the API HTTP server
  frontend_addr: :8080 # the listen address for the HTTP web console
tests: # a list of tests (plugins) to run and their intervals
- command: docker-images.py --threshold 10
  interval: 5
- command: load-average.rb --threshold 2 --expect true
  interval: 30
- command: stringmatch.sh /tmp/test.txt dog 3
  interval: 60 
- command: process-running.js --name foobar --expect false
  interval: 30 
- command: process-running.js --name ssh-agent --expect true
  interval: 30

```

### Writing custom plugins

As long as your plugin outputs conforms to the following JSON standard your plugin will work:

```
{
    "description": "the total number of docker images must not exceed 5",  // describes the test
    "result":      false,                                                  // did the test pass?
    "data":        28,                                                     // the result data that the result logic was based on.
    "expect":      true                                                    // what result were you expecting?
}
```