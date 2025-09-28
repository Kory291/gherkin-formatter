# Gherkin Formatter

This project still is very much work in progress. But I want to try to have a formatter for Gherkin.
This will for now be focused on the project structure that is guided by the behave package since it is the domain I want to use this in.

## ToDo

### Scenario discovery

- [ ] be able to scan for all `.feature` files in project structure
Note: Assumed project structure is that all `.feature` files are in a directory `features/` or in sub-directories within that `feature/` directory.

### Formatting options

- [ ] set default value for intendation
This will be 2 spaces for now
- [ ] set allignment

### Configuration

- [ ] Have input file to define configuration
- [ ] Configuration options are supported in `.pyproject.toml`
