#### Run
```
go run main.go 3000
```

#### GHE workflow

```yml
on: [push]
name: "Continuous Integration"

jobs:
  test:
    strategy:
      matrix:
        node-versions: ['10.16.3', '14.9']
    runs-on: self-hosted
    env:
      CACHE_HTTP_API: your-url.com
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
      run: |
        hit_tar=${{ runner.os }}-node-${{ matrix.node-versions }}-${{ hashFiles('yarn.lock') }}.tar.gz
        hit_url=https://${{env.CACHE_HTTP_API}}/cache-http/assets/$hit_tar

        install_command="yarn install"
        dst_folder=node_modules

        echo hit_url: $hit_url
        curl -X GET -u ${{ secrets.CACHE_HTTP_API_USER }}:${{ secrets.CACHE_HTTP_API_PASS }} \
          --noproxy ${{ env.CACHE_HTTP_API }} \
          -s $hit_url > $hit_tar
        file $hit_tar | grep -q 'gzip compressed data' && response=200 || response=404

        { test "$response" == 404; } && echo Cache Miss && \
          rm -rf $dst_folder && \
          $install_command && \
          tar zcf $hit_tar $dst_folder && \
          curl -X POST -u ${{ secrets.CACHE_HTTP_API_USER }}:${{ secrets.CACHE_HTTP_API_PASS }} \
            --noproxy ${{ env.CACHE_HTTP_API }} \
            --form file=@$hit_tar ${{env.CACHE_HTTP_API}}/cache-http/upload && \
          echo "Cache upload OK"

        tar xzf $hit_tar
        echo Dependency Install finished


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