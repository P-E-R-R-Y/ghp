
# Dev

retrieve executable using `go build`

example: `./gitperry --user P-E-R-R-Y`

# Release

install executable on $Home/go/bin: `go install github.com/P-E-R-R-Y/gitperry`

add go bin path to the env: `export PATH=$PATH:$(go env GOPATH)/bin`

example `~/ gitperry -user P-E-R-R-Y --filter "^i.*`

# Usage

gitperry

cmd:
- info get info about a spec org or user 
- list retrieve repo from a spec orga or user 

info
- user organsatiob or user