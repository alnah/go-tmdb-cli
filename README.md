# go-tmdb-cli

A simple command-line interface (CLI) to fetch and display data from [The Movie Database (TMDB)](https://www.themoviedb.org/). The idea is taken from [roadmap.sh](https://roadmap.sh/projects/tmdb-cli).

![Demo](demo.gif)

## Install

Configure the TMDB API key:

- The CLI looks for a YAML file in your **home directory**: `~/.go-tmdb-cli/config.yaml`.
- The file must include your TMDB API key in the following format: `api_key: YOUR_API_KEY`.
- [Get an API Key](https://developer.themoviedb.org/docs/getting-started).
- By default, `config.yaml` is expected, you can pass a different file to `newRootCmd("filename.yaml")` in `main.go`.

Setup the CLI:

- Option 1: Install the CLI globally, making it accessible from anywhere in your terminal by running the command:

  ```
  go install github.com/alnah/go-tmdb-cli@latest
  ```

- Option 2: Clone the repo and use the Makefile to compile it inside `./bin/go-tmdb-cli`:

  ```
  git clone https://github.com/alnah/go-tmdb-cli && cd go-tmdb-cli && make
  ```

- Option 3: Download the appropriate [artifact](https://github.com/alnah/go-tmdb-cli/releases):
  - For Windows, download the `.zip` file.
  - For Linux and macOS, download the `.tar.gz` file.
  - Extract the downloaded file and install the CLI.

## Usage

Fetch curated lists like **now playing**, **popular**, **top rated**, and **upcoming** movies directly from TMDB:

```
go-tmdb-cli list -n
go-tmdb-cli list -p
go-tmdb-cli list -t
go-tmdb-cli list -u
```

Specify filters such as **language**, **year**, **average rating**, **genres**, etc., to discover movies:

```
go-tmdb-cli discover -l=en -y=2000,2005 -g=comedy,action -a=6.5,10 -v=100,50000 -m=100 -s=average,desc
go-tmdb-cli discover -l=fr -y=1960,gte -g=history -a=7,gte -v100,gte -m=50 -s=title,asc
go-tmdb-cli discover -l=pt -y=1960,lte -w=comedy -a=9.0,lte -v=2000,lte -m=10 -s=votes,asc
```

Fore more details:

```
go-tmdb-cli [command] --help
```

Run all tests and benchmarking:

```
make test && make benchmark
```

## Dependencies

- [golangci-lint](https://github.com/golangci/golangci-lint) for static code analysis and linting.
- [tablewriter](https://github.com/olekukonko/tablewriter) for tabular data.
- [cobra](https://github.com/spf13/cobra) for CLI structure.
- [viper](https://github.com/spf13/viper) for configuration.
- [backoff](https://github.com/cenkalti/backoff) for retrying on failed requests.

## License

This project is distributed under the Apache License. See the [license](LICENCE) file for more details.
