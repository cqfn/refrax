# Refrax

[![codecov](https://codecov.io/gh/cqfn/refrax/branch/master/graph/badge.svg)](https://codecov.io/gh/cqfn/refrax)

**Refrax** is an AI-powered refactoring agent for Java code, implemented in Go. It communicates using the [A2A protocol](https://google-a2a.github.io/A2A/latest/specification/).

> ⚠️ Early prototype — subject to rapid change.

## Installation

### Releases

Download the latest stable version from the [releases page](https://github.com/cqfn/refrax/releases). Pre-built binaries are available for MacOS, Windows, and Linux.

### From Sources

You need to have Go 1.24.1 or later installed on your system.

1. Clone the repository:

   ```bash
   git clone https://github.com/cqfn/refrax.git
   cd refrax
   ```

2. Build the binary:

   ```bash
   go build -o refrax
   ```

3. (Optional) Install the binary to your `$GOPATH/bin`:

   ```bash
   go install
   ```

## Usage

*[![asciicast](https://asciinema.org/a/IHrW8v68VS81vVNfw8ByioG4T.svg)](https://asciinema.org/a/IHrW8v68VS81vVNfw8ByioG4T)

- `refrax refactor [path]`: Refactor Java code in the specified directory (defaults to current directory).
- `refrax start [agent]`: Start the server for agents like fixer, critic, or facilitator.

## Configuration

- `--ai, -a`: Specify the AI provider (e.g., deepseek).
- `--token, -t`: Token for the AI provider.
- `--debug, -d`: Enable debug logging.

## License

Licensed under the [MIT](LICENSE.txt) License.

