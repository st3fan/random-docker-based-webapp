# Technology

## Application and development

Build the web app with Go modules and the standard library’s `net/http` package. Format Go source with `gofmt`, run `go test ./...`, and verify that the Docker image builds.

## Packaging and local use

Package the app as a single container image using a Dockerfile. Developers build the image from the repository and run it locally with `docker run`.

## Intent previews

Use the self-hosted Preview Controller in single-container mode. For each SDLC Intent, the controller builds the selected Git commit from the Dockerfile. The app must serve the configured HTTP health path. The controller switches the Intent’s stable URL to the new container after its health check passes. [Preview Controller documentation](https://github.com/st3fan/preview-controller)

Local preview commands read the controller connection settings from the repository-root `.env`. Keep that file out of Git and the application image.

## Open decisions

Persistence, image publishing, production hosting, and production deployment remain undecided.
