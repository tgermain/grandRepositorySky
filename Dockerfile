# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/tgermain/grandRepositorySky
ADD ./web/client /static/ 



# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/nu7hatch/gouuid
RUN go get github.com/op/go-logging
RUN go get github.com/gorilla/mux
RUN go get github.com/spf13/cobra
RUN go install github.com/tgermain/grandRepositorySky

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/grandRepositorySky -s /static/

# Document that the service listens on port 4321.
EXPOSE 4321