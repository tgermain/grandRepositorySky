# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/tgermain/grandRepositorySky
ADD ./web/client /static/ 



# Build the outyet command inside the container.
# fetching dependencies
RUN go get github.com/nu7hatch/gouuid
RUN go get github.com/op/go-logging
RUN go get github.com/gorilla/mux
RUN go get github.com/spf13/cobra

# install the main project
RUN go install github.com/tgermain/grandRepositorySky


# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/grandRepositorySky"]


# default command
CMD ["--help"]