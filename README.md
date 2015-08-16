# Torokko

A build proxy server using Docker container for Golang apps.

## Quick start

1. Start a torokko server.
2. Access a build endpoint.

```sh
$ curl -X POST http://torokko-server/{remote}/{owner}/{repo}/{GOOS}/{GOARCH}/{version}
# e.g. http://torokko-server/github.com/monochromegane/torokko/linux/amd64/v0.0.1
```

And access a download endpoint after a few minutes :beers:

```sh
$ curl -OJL http://torokko-server/{remote}/{owner}/{repo}/{GOOS}/{GOARCH}/{version}
```

[Here](http://torokko.monochromegane.com) is a demo API server. Try it now :)

## Overview

![torokko_overview](https://cloud.githubusercontent.com/assets/1845486/9293092/eba8190c-4457-11e5-9176-19d9f7ac3363.jpg)

## Endpoints

### Build

Build a go binary.

`POST /{remote}/{owner}/{repo}/{GOOS}/{GOARCH}/{version}`

**Example request:**

```sh
$ curl -X POST http://torokko-server/github.com/monochromegane/torokko/linux/amd64/v0.0.1
```

- **remote** - Remote repository. (e.g. github.com)
- **owner** - Repository owner.
- **repo** - Repository name.
- **GOOS** - Cross compilation environment. (e.g. linux, darwin, windows)
- **GOARCH** - Cross compilation environment. (e.g. amd64, 386)
- **version** - Repository version tag name. (e.g. v0.0.1)

See also **Custom build** section.

**Example response:**

```json
{
  "build_id": "1a7452d077faed659af0e85731664938"
}
```

- **build_id** - build id (See also **Log** endpoint)

**Status code**

- **202** - no error (build offer is accepted)
- **409** - already exists (The binary is already built, try download)
- **500** - server error

### Log

Get build logs.

`GET /builds/{build_id}/logs`

**Example request:**

```sh
$ curl http://torokko-server/builds/1a7452d077faed659af0e85731664938/logs
```

**Example response:**

```json
{
  "build_id": "1a7452d077faed659af0e85731664938",
  "level": "info",
  "msg": "Your build started.",
  "time": "2015-08-13T11:46:26+09:00",
  "workspace": "795098082"
}
{
  "build_id": "1a7452d077faed659af0e85731664938",
  "level": "info",
  "msg": "checking binary...",
  "time": "2015-08-13T11:46:26+09:00"
}
{
  "build_id": "1a7452d077faed659af0e85731664938",
  "level": "info",
  "msg": "cloning repository...",
  "time": "2015-08-13T11:46:26+09:00"
}
```

**Status code**

- **200** - no error
- **500** - server error

### Download

Download a go binary.

`GET /{remote}/{owner}/{repo}/{GOOS}/{GOARCH}/{version}`

- **remote** - Remote repository. (e.g. github.com)
- **owner** - Repository owner.
- **repo** - Repository name.
- **GOOS** - Cross compilation environment. (e.g. linux, darwin, windows)
- **GOARCH** - Cross compilation environment. (e.g. amd64, 386)
- **version** - Repository version tag name. (e.g. v0.0.1)

**Example request:**

```sh
$ curl -OJL http://torokko-server/github.com/monochromegane/torokko/linux/amd64/v0.0.1
```

- Specify `-OJL` option, because torokko server redirect and add `Content-Disposion` header.
- If you use `wget`, try `--content-disposition` option.


**Status code**

- **200** - no error
- **404** - not found
- **500** - server error

## Custom build

Torokko use `make` command for building your app.
If your repository don't have `Makefile`, Torokko use default Makefile:

```make
build:
	go get -d ./...
	go build
```

If your app have to customize build step, put a Makefile on your repository.

## Private repository

If your repository is private, you can add `Authorization: token <TOKEN>` header to build and download requests.

```sh
$ curl -H "Authorization: token <TOKEN>" -X POST http://torokko-server/github.com/monochromegane/torokko/linux/amd64/v0.0.1
```

**Caution**

[Demo API server](http://torokko.monochromegane.com) is **non-SSL** site.
I don't recommend you use for private repository.

## Installation

```sh
$ go get github.com/monochromegane/torokko
```

## Requirement

Torokko server require the following.

- Docker
- Docker Remote API
- Golang build images (Now torokko use `golang:1.4.2-cross`)
- git
- lsof
- workspace, log, storage directories.

## TODO

- Support build trigger from GitHub Webhook.
- Dockernize.
- Separate build queue process.
- Support GitHub releases backend.
- Support S3 backend.
- Add more tests.

## Contribution

1. Fork it
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## License

[MIT](https://github.com/monochromegane/torokko/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
