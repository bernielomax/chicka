# Chicka
A pluggable Go monitoring system

### Building

To build Chicka you need a working Go environment.

#### OSx:

``` 
make osx
```

#### Linux:

``` 
make linux
```

#### Docker

``` 
make docker
```

This will place a working `chicka` executable inside the Chicka source folder.

To run the docker container run the following:

``` 
docker run --name chicka -d -p "8080:8080" -p "9090:9090" -v $(pwd)/test.txt:/tmp/test.txt -v /var/run/docker.sock:/var/run/docker.sock chicka
```

### Configuration

#### Paths:

The configuration for Chicka can belong in the following locations:

* `$HOME/.chicka/config.yaml`
* `/etc/chicka/config.yaml`
* `./config.yaml`

#### Example:

The following is an example of the config.yaml

```
git:
  url: git@github.com:bernielomax/chicka-plugins.git # git plugin repo URL
plugins:
  path: /Users/justin.pye/Documents/chicka-plugins/ # local path where plugins will reside
logging:
  path: /tmp/chicka.log # local path where the logs will reside
cache:
  ttl: 120 # amount of time to store tests in cache
http:
  api: :9090 # the listen address for the API HTTP server
  www: :8080 # the listen address for the HTTP web console
tests: # a list of tests (plugins) to run along with their interval and expected result.
- command: docker-images.py --threshold 10
  interval: 5
  expect: true
- command: load-average.rb --threshold 2
  interval: 30
  expect: true
- command: stringmatch.sh /tmp/test.txt dog 3
  interval: 60 
  expect: true
- command: process-running.js --name foobar --inverse true
  interval: 30 
  expect: true
- command: process-running.js --name ssh-agent
  interval: 30
  expect: true

```

### Get plugins via git

Chicka has a builtin Git integration. The `git.url` configuraiton must be specified in the `config.yaml`. Plugins can then be downloaded to the `plugins.path` location by running the following command:

```
chicka get
```

#### Git SSH authentication

```
chicka get --help
       get latest chicka plugins repo
       
       Usage:
         chicka get [flags]
       
       Flags:
         -h, --help             help for get
         -k, --ssh-key string   The SSH key to use for Git authentication (default "/Users/justin.pye/.ssh/id_rsa")
```

### Plugins

[Chicka plugins can be found here](https://github.com/bernielomax/chicka-plugins)

### Writing custom plugins

As long as your plugin outputs conforms to the following JSON standard your plugin will work:

```
{
    "description": "the total number of docker images must not exceed 5",  // describes the test
    "result":      false,                                                  // did the test pass?
    "data":        28                                                      // the data that the result was based on.
}
```

### API


You can access the API using the following URL (depending on your `http.api` config):

``` 
http://localhost:9090
```

### Web Console


You can access the web console using the following URL (depending on your `http.www` config):

``` 
http://localhost:8080
```
