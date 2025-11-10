# Gherkin Formatter

This project still is very much work in progress. But I want to try to have a formatter for Gherkin.
This will for now be focused on the project structure that is guided by the behave package since it is the domain I want to use this in.

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
- [ ] Have `configure` command available to help configuration

### Formatting

- [x] Add command `format` to apply configuration
- [ ] Add support to extend tags
- [ ] Add handling of multiple empty lines

### Code quality

- [ ] Add unit tests
- [ ] Add linter, unit-tests in CI

### Building

- [x] Have a build and release pipeline available
