# Mocks

We use [mockery](https://github.com/vektra/mockery) for generating mocks from Go interfaces for unit tests. 

## Adding a new Mock

To add a new mock, append a line to the `Makefile` in this directory in the following form:

```sh
mockery --name=InterfaceName --dir=path/to/package --recursive --output=./mocks
```

Note that if the mock being generated is for an external package, you'll need to use the $(GOPATH) variable to reference the package, otherwise the `--dir` argument should be relative to the root of this repository.

After adding your Mock to the Makefile, run `make mock-gen` from the repository root. Mocks are checked in to source control along with your tests.

_Be aware that updating any of the interfaces used by mocks will require you to rerun `make mock-gen`, otherwise any tests using those mocks will fail to compile (as the interfaces will no longer match)._
