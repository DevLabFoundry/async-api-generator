# Usage

## SourceCode Annotation

The program works by identifying markers in the source code, and extracing them and sorting them based on ancestral precedence - i.e. who is a parent/child/grandchild.

### Markers and Annotation

Markers have to be used in the exact form, beginning `//+gendoc annotationKey=annotationVal` and end `//-gendoc`.

### Tips

You can do multiline annotations by adding `\` before a line break e.g. 

```csharp
public class Bar() {
    /*
    //+gendoc category=message \
    type=example 
    */
    public class Foo() {

    }
    //-gendoc
}
```

## CLI

Download the published binary from [here](TODO).

> replace `darwin` in `.../gendoc-darwin/overview` for your platform = `linux|windows|darwin`

```sh
chmod +x ./gendoc-$PLATFORM
mv ./gendoc-$PLATFORM /usr/local/bin/gendoc
```

After moving to a location which is on your $PATH, you should be able to run the following commands.

The CLI comes with the following commands, for furher info and help use the `--help|-h` flag at any point.

```bash
gendoc -h
gendoc --version
```

### Commands

The CLI has 2 main commands that are run either against a single repo that generates the interim output and the command that reads in the interim output and generates the AsyncAPI compliant document.

- `--input`|`--output` options currently support 2 types of `"storage implementation"`
    - `local://` => pointing to a local filesystem
    - `azblob://` => pointing to an Azure storageaccount/blob in this format `azblob://STORAGE_ACCOUNT_NAME/CONTAINER_NAME`. The utility handles the virtual path and object creation.
    - additional `storageClients` can be added easily by providing a new implementation on the storageAdapter

For ease of use, you can enable shell completion for your shell.

`gendoc completion -h` to see the options, e.g. for powershell `gendoc completion powershell`

>Not tested on Windows, users may need to suffix the binary with `.exe`.

#### SingleContext

This command is run against a single directory which holds source code of any type and generates an interim output.

It can be run in validate only mode by setting  `--dry-run` flag, it will ensure that any annotations that have been added are set correctly and there is no syntax errors.

```sh
gendoc single-context --input local:///path/to/src/domain.sample --output local:///path/to/out/interim --is-service --verbose --dry-run
```

If any errors occur they are returned to the terminal. In dry-run mode no files are emitted to the output location.

```sh
gendoc single-context --input local:///path/to/src/domain.sample --output local:///path/to/out/interim --is-service --business-ctx s2s --business-domain domain
```

`--business-ctx` and `--business-domain` are purely for tagging/description/name generation purposes

> Currently `--input` for single-context can only be a `local://` i.e. stored on the local filesystem

##### EnvVariable expansion

The content can include environment variable like text to avoid repetition, however it will fail if the variable is not set.

Example:

```text
//+gendoc category=message type=description id=foo
this is some description with $foo
//-gendoc
```

Ensure that the environment variable is present otherwise it will fail with either an `unset` or `set but empty` error.

See tests for [more examples](../src/go/async-api-gen-doc/internal/parser/parser_test.go)

#### GlobalContext

This command is run against a directory containing zero or more interim output files from generated from across many repos or (single-context sources). 

```sh
gendoc global-context --input local:///path/to/src/domain.sample --output local:///path/to/out/interim
```

### Local Example

Point it to an input directory of any repo - e.g. `domain.Foo.BarActionQux`.

This will generate the interim code that the 

```sh
gendoc single-context --input local://$FULL_PATH_TO/domain.Foo.BarActionQux --is-service --bounded-ctx Foo --business-domain domain \
--repo "https://github.com/repo" \
--output local://$HOME/.gendoc/poc
```

The output will be populated with a directory called `current` which will include the interim output(s) from the single-context runs.

This is then used as in input for the global-context and it will output a full AsyncAPI document in the output directory in this case `local://$HOME/.gendoc/poc/processed` here.

```sh
gendoc global-context --input local://$HOME/.gendoc/poc/current --output local://$HOME/.gendoc/poc/processed
```

The files are emitted with the `AsyncAPI.ID` as the name in the `asyncapi` directory, e.g.: `asyncapi/urn:domain:Foo:domain.Foo.BarActionQux.yml`.
