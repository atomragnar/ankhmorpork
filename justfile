# print available targets
default:
    @just --list --justfile {{justfile()}}

# evaluate and print all just variables
evaluate:
    @just --evaluate

# print system information such as OS and architecture
system-info:
  @echo "architecture: {{arch()}}"
  @echo "os: {{os()}}"
  @echo "os family: {{os_family()}}"


golangci_lint_exe := "./bin/golangci-lint"
go_app_path := "."

# Recipes
@bin-directory:
    # Check if the bin directory exists, if not, create it
    if [ ! -d "./bin" ]; then \
        echo "Creating bin directory..."; \
        mkdir -p ./bin; \
    fi

@golangci-lint: bin-directory
    # Check if golangci-lint is installed, if not, install it
    if [ ! -f {{golangci_lint_exe}} ]; then \
        echo "Installing golangci-lint..."; \
        curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.55.2; \
    fi


lint: golangci-lint
    # Running linting
    echo "Running golangci-lint..."
    {{golangci_lint_exe}} run ./filetree/*.go --verbose


installed-go-tools:
    go list -m -u -json all

format:
    @echo "Formatting source code ..."
    gofmt -l -s -w .