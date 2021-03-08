#### About

Alternate solution to action/cache for GHES and self-hosted runners on Docker.

https://github.com/actions/cache/issues/505

#### Run
```
go run main.go 3000
```

#### GHE workflow

https://github.com/marketplace/actions/action-cache-http

```yml
on: [push]
name: "Continuous Integration"

jobs:
  test:
    strategy:
      matrix:
        node-versions: ['10.16.3', '14.9']
    runs-on: self-hosted
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
      uses: kevincobain2000/action-cache-http@v1.0.3
      with:
        version: ${{ matrix.node-versions }}
        lock_file: yarn.lock
        install_command: yarn install
        destination_folder: node_modules
        cache_http_api: "https://yourdomain.com/path/to/installation/cache-http"
        http_proxy: ""
```

#### Deploy

```
go get github.com/lestrrat/go-server-starter/cmd/start_server
start_server --pid-file /path/to/pids/3000.pid -- main 3000 > /dev/null 2>&1 &
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