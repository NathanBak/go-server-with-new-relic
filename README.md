# go-server-with-new-relic
A sample go server that can talk to New Relic

## Introduction
This repo was original copied from my [go-server](https://github.com/NathanBak/go-server) and then I added in bits an pieces so that the server can send data to New Relic.

## Using this code
This repo is published under an [MIT License](LICENSE) which tends to be flexible for the user while still protecting the creator.  Please feel free to use this code and also to submit Pull Requests with any fixes/improvements/suggestions/etc.  Also, if you use code from the repo I request (but do not require) that you star the repo in GitHub.

## Running the code

### Setup development environment
The server uses properties stored in a `.env` file.  To generate said file, from the project root run `scripts/create_env.sh`.

### Running the Server
From the project root run `go run .`

### Building the Docker image
From the project root run `docker build .`

If you would prefer to use a pre-build image, see [https://hub.docker.com/repository/docker/bakchoy/go-server-with-new-relic/general](https://hub.docker.com/repository/docker/bakchoy/go-server-with-new-relic/general).

### Prerequisites
- Go (see the [go.mod](https://github.com/NathanBak/go-server/blob/755e067fd4b192641c8478422a49549e316e137c/go.mod#L3) file for the correct version) 
- Git
- New Relic account (free tier works fine)