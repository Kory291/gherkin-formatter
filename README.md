# Gherkin Formatter

This project still is somewhat work in progress - depends on my motivation.
This will for now be focused on the project structure that is guided by the behave package since it is the domain I want to use this in.

## How to use:

To just get an idea what will be done: 
```
gherkin-formatter format
```


If your sure what you are doing you can also run:
```
gherkin-formatter format --write
```
This will write to the `.feature` files that were found.


You can also initialize a configuration file with:
```
gherkin-formatter configuration init
```
This will create a file `gherkinFormatter.toml`

## Installation

```
go install github.com/Kory291/gherkin-formatter
```

## ToDo

### Scenario discovery

- [x] be able to scan for all `.feature` files in project structure
Note: Assumed project structure is that all `.feature` files are in a directory `features/` or in sub-directories within that `feature/` directory.
- [x] read feature files that where found and save them

### Formatting options

- [x] set default value for intendation
This will be 2 spaces for now
- [ ] set allignment

### Configuration

- [x] Have input file to define configuration
- [ ] Configuration options are supported in `.pyproject.toml`
- [x] Have `configuration init` command available to help configuration

### Formatting

- [x] Add command `format` to apply configuration
- [x] Add support to extend tags
- [x] Add handling of multiple empty lines

### Code quality

- [x] Add unit tests
- [x] Add linter, unit-tests in CI

### Misc

- [ ] Have logging

### Building

- [x] Have a build and release pipeline available
