## About

Alternate solution to action/cache for GHES and self-hosted runners on Docker.
[Github's "action/cache" plugin current doesn't support GHES and self-hosted runners.]
(https://github.com/actions/cache/issues/505)

## Github Enterprise / self-hosted runner workflow

The syntax is varies depending on whether your app that needs the cache is running 
in a Docker container or directly on the host.

Both cases the companion ["action-cache-http"](https://github.com/marketplace/actions/action-cache-http) Action.

When your app is built in a docker container, Github Actions will automatically will assign
the host name "cache-http" to this service container, so the configuration will
access by that host on port 80.

When your app is built directly on the host, you'll need to map port 80 inside the cache-http
container to a free port on your host. In the example below, we use port 3000.

In both cases we use a bind mount to allow the service container to persist the assets on the host.

Currently, there is no automated cache pruning but see [clear\_retention.sh(./clear\_retention.sh)
for an example that you could run on a cron job or system timer.

Note: There's no authentication on the default Docker image. It's assumed that access to
the service will be restricted to Dockerized builts or only to "localhost" requests from
your CI server.


```yml
on: [push]
name: "Continuous Integration"

jobs:
  test:
    strategy:
      matrix:
        node-versions: ['10.16.3', '14.9']
    runs-on: self-hosted
    # Cache service for GHE / self-hosted runners
    services:
      cache-http:
        image: docker.pkg.github.com/SOMEPATH/FIXME/cache-http:3
        # For now you still need to authenticate when accessing public packages
        credentials:
          username: $GITHUB_ACTOR
          password: ${{ secrets.SOME_SECRET }}
       # The left side is any path on your host you want to store the cache
       # Make sure the directory exists. Right side must be "/app/assets!"
        volumes:
          - /home/github/.cache/cache-http/assets:/app/assets
        # Only map ports for direct-on-host case. Not needed for Dockerized builds!
        ports:
           - 127.0.0.1:3000:80


    steps:
    - name: Setup Node.js ${{ matrix.node-versions }}
      uses: actions/setup-node@v1
      with:
        node-version: ${{ matrix.node-versions }}

    - name: Checkout
      uses: actions/checkout@v2

    # - name: Yarn Install (without cache)
    #   run: yarn install

    - name: Yarn Install (with cache)
      uses: kevincobain2000/action-cache-http@v3
      with:
        version: ${{ matrix.node-versions }}
        lock_file: yarn.lock
        install_command: yarn install
        operating_dir: ./
        destination_folder: node_modules
        # For Dockerized builds, this should match the service name.
        # cache_http_api: "http://cache-http"
        # For direct-on-host builds, use the custom port you mapped above
        cache_http_api: "http://127.0.0.1:3000"
        http_proxy: ""
```

If for some reason you want to use the Dockerized version of this service, below are details
on how to run it yourself. Also see the example `.sh` scripts distributed with the project.

**Without Docker**

```yml
    - name: Yarn Install (with cache)
      uses: kevincobain2000/action-cache-http@v3
      with:
        version: ${{ matrix.node-versions }}
        lock_file: yarn.lock
        install_command: yarn install
        operating_dir: ./
        destination_folder: node_modules
        # For Dockerized builds, this should match the service name.
        # cache_http_api: "http://cache-http"
        # For direct-on-host builds, use the custom port you mapped above
        cache_http_api: "http://127.0.0.1:3000"
        http_proxy: ""
```

#### Run
```
go run main.go 3000
```

#### Deploy

```
go build -o main
./main -host=localhost -port=3000 -pidDir=./ > /dev/null 2>&1 &
```

#### Development

```
go get github.com/cespare/reflex
reflex -r '\.go$' -s -- sh -c 'go run main.go 3000'
```

```
curl -X GET -u joe:secret localhost:3000/health
```

```
curl -X POST -u joe:secret --form file=@/path/to/file/sample.tar.gz localhost:3000/upload
curl -X GET -u joe:secret localhost:3000/assets/sample.tar.gz
```
