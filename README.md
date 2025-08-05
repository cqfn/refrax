# Refrax

[![codecov](https://codecov.io/gh/cqfn/refrax/branch/master/graph/badge.svg)](https://codecov.io/gh/cqfn/refrax)

**Refrax** is an AI-powered refactoring agent for Java code, implemented in Go. It communicates using the [A2A protocol](https://google-a2a.github.io/A2A/latest/specification/).

> ⚠️ Early prototype — subject to rapid change.

## Installation

### Releases
Download the latest stable version from the [releases page](https://github.com/cqfn/refrax/releases). Pre-built binaries are available for macOS, Windows, and Linux.

### Using Go

If you have Go 1.24.1 or later installed, you can run:

```bash
go install github.com/cqfn/refrax@latest
```

To install a specific version, use:

```bash
go install github.com/cqfn/refrax@v0.0.1
```

[Releases page](https://github.com/cqfn/refrax/releases).

### From Source

Ensure that Go 1.24.1 or later is installed on your system.

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

### Example

You can try refactoring the testing project located in this repository. To do so, you will need to clone the repository:

```
git clone https://github.com/cqfn/refrax.git
```

Then, run the following command:

```sh
refrax refactor --output="./out" --ai=deepseek refrax/test/test_data/java/person
```

Or, if you are already in the `refrax` folder, simply run:

```
refrax refactor --output="./out" --ai=deepseek test/test_data/java/person
```

## Configuration

- `--ai, -a`: Specify the AI provider (e.g., deepseek, openai).
- `--token, -t`: Token for the AI provider.
- `--debug, -d`: Enable debug logging.

## Authentication

Some operations in Refrax require AI authentication using an API token. You can provide the token using one of the following methods:

### Command-Line Flag

```sh
refrax refactor . --token your-token-here
```

### AI Providers

Supported AI providers are:
* `deepseek`
* `openai`

### Environment Variable

✅ The `DEEPSEEK_TOKEN` variable is the recommended option for `deepseep` AI provider
✅ The `OPENAI_TOKEN` variable is the recommended option for `openai` AI provider
⚠️ The `TOKEN` variable is still supported for any AI provider but deprecated.


Set the environment variable:

```sh
export DEEPSEEK_TOKEN=your-token-here
refrax start facilitator
```

### `.env` File

If a `.env` file is present in the working directory, Refrax will attempt to read the token from it:

```
# .env
DEEPSEEK_TOKEN=your-token-here
```

### Priority Order

If multiple sources are provided, the following priority order is applied (highest priority first):

1. `--token` command-line flag
2. `DEEPSEEK_TOKEN` environment variable
3. `TOKEN` environment variable (deprecated)
4. `.env` file (`DEEPSEEK_TOKEN` > `TOKEN`)

## Statistics

To gather interaction statistics, you can use the following command:

```sh
refrax refactor . --ai=deepseek --stats --stats-format=csv --stats-output=stats.csv
```

This command generates a `stats.csv` file containing the interaction statistics.
The `--stats-output` and `--stats-format` parameters are optional.
If you omit them, `refrax` will output the statistics directly to the console.

## License

Licensed under the [MIT](LICENSE.txt) License.
