# addlicense
This project adds the given license to all source file in a given directory.
The tool will try to avoid source files that already has a header. This
project can be used as a CLI as well as a library imported to a go script.

## Inspiration
The initial inspriation is from this
[repo](https://github.com/google/addlicense) from Google. While the google
addlicense project is good, but there are some additional requirements I
needed.

1. I want the code to be a importable library in addition to just a tool. In
this way, I can import this to write my own customized scripts. 
2. Template is an overkill for most people. Most of the time, I just want to
attach a license to all my files. If I am making changes to the license, such
as updating the year, I will just update a single license file. So I
simplified the codebase.

## Build
```bash
make
```

## Test
```bash
go test -v ./...
```

## Usage
This project can be consumed either as a CLI or as a lib.

### CLI

To install:
```bash
go get -u github.com/unkies/addlicense/cmd/addlicense
```
Note: You may need to add you `$GOPATH/bin` to your `$PATH`.

To uninstall:
```bash
go clean -i github.com/unkies/addlicense/cmd/addlicense
```

To use, first create a file with the license header you want to add. See
`./LICENSE` as an example of the MIT license.
```bash
addlicense --license <license file> --ignore <if any ignore patterns> <root of source code>
```

For usage, consult `addlicense -h` and `make dogfood` as an example of adding
`./LICENSE` to this repo.

### Lib

Follow the reference [here](https://pkg.go.dev/mod/github.com/unkies/addlicense)

```golang
import "github.com/unkies/addlicense/libaddlicense"

// Read the license file
license, _ := ioutil.ReadFile("./LICENSE")

// Add licenses to a project root, but ignore certain files
srcRoot := "."
ignore := []string{"testdata"}
_ := libaddlicense.AddLicense(srcRoot, license, ignore)
// Or no need to ignore
_ := libaddlicense.AddLicense(srcRoot, license, []string{})

// To add a license to a file
_ := libaddlicense.AddLicense(path, license, []string{})

```

## License

MIT license.