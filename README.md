# HellSpawner

![Logo](hsassets/images/d2logo.png)

[Join us on Discord!](https://discord.gg/pRy8tdc)\
[Development Live stream](https://www.twitch.tv/essial/)\
[Support us on Patreon](https://www.patreon.com/bePatron?u=37261055)

## About this project

HellSpawner is the toolset used for developing games on the [OpenDiablo2](https://github.com/OpenDiablo2/OpenDiablo2) engine.

## Getting the source

To download code use a following command:

`git clone --recurse-submodules https://github.com/OpenDiablo2/hellspawner`

Then, you need to install [Go](https://golang.org/doc/install) 1.16.2 or higher.

Then you can use the Makefile (try `make help` for more information:

```shell
make run # (works for Windows, MacOS, Ubuntu, Fedora)
```

*Note*: If you're using any other Unix-based OS, you can build project with included building script: run `./build.sh`

## Contributing

If you find something you'd like to fix that's obviously broken, create a branch, commit your code, and submit a pull request. If it's a new or missing feature you'd like to see, add an issue, and be descriptive!

If you'd like to help out and are not quite sure how, you can look through any open issues and tasks, or ask
for tasks on our discord server.

### Lint Errors

We use `golangci-lint` to catch lint errors.
To install it use: `make setup` and to run it use: `make lint`.

### Testing

We encourage you to write tests for your code.
To run the tests, use `make test` and `make race` (this works only for Linux systems)

*Note*: Use `make cover` to check the testing coverage of your code.

## VS Code Extensions

The following extensions are recommended for working with this project:

*   ms-vscode.go
*   defaltd.go-coverage-viewer

When you open the workspace for the first time, Visual Studio Code will automatically suggest these extensions for installation.

Alternatively you can get to it by going to settings <kbd>Ctrl+,</kbd>, expanding `Extensions` and selecting `Go configuration`,
then clicking on `Edit in settings.json`. Just paste that section where appropriate.

## Status

For now (start of march 2021) we can use HellSpawner to:

*   create projects
*   browse MPQ archives
*   view following file formats:
    *   DC6 - and DCC animations
    *   WAV - sound files
    *   TXT - data tables
    *   PL2 - palette transforms
    *   DS1 - map data
*   edit:
    *   COF - animation data
    *   TBL - font tables
    *   TBL - string tables
    *   TXT - text files
    *   DAT - palettes
    *   DT1 - map tiles
    *   D2 - animation data

Much work has been made in the background, but a lot of work still has to be done for the project to be complete.

Feel free to contribute!

## Screenshots

![Screenshot](docs/overview.png)
