# Command directory

This directory houses custom [cobra](https://github.com/spf13/cobra) commands which can be added as subcommands to the root `klyraprotocold` command defined in `cmd/klyraprotocold/main.go`.

Conventionally, each package in this directory should define a public `Command()` method which returns a `*cobra.Command`.

These commands can be added as a subcommand to the `klyraprotocold` root command defined in `main.go` in the following way:

```go
rootCmd, _ := NewRootCmd(...)

rootCmd.AddCommand(mycommandpkg.Command())
```

The above will surface your command as `klyraprotocold mycommand`.

If instead you wish to define your command as the subcommand of a subcommand (i.e. You wish define something like `klyraprotocold tendermint mycommmand`), you can first search for the subcommand (i.e. `tendermint`) and subsequently add your command to it like so:

```go
// Fetch Tendermint subcommand.
tmCmd, _, err := rootCmd.Find([]string{"tendermint"})
if err != nil {
  os.Exit(1)
}

// Add "mycommand" command to Tendermint subcommand.
tmCmd.AddCommand(mycommandpkg.Command())
```
