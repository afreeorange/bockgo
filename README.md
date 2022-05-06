## Development Notes

```bash
# Initialize this project
go mod init afreeorange/bock

# Remove unused mods
go mod tidy

# Remove a package
go get package@none

# Build. Version info is injected at build...
go build -o ~/Downloads/bock -v -ldflags="-X main.version=1.0.0" .
```

```go
r, _ := git.PlainOpen(*articleRoot)

commits, _ := repository.Log(&git.LogOptions{FileName: &fileName})
commits.ForEach(func(c *object.Commit) error {
  f, err := c.Files()

  if err != nil {
    fmt.Println("Could not get files for commit: ", c.Hash)
  } else {
    f.ForEach(func(f *object.File) error {
      if f.Name == fileName {
        fileContents, _ := f.Contents()
        render([]byte(fileContents), buffer)

        fmt.Println("---", c.Hash.String())
        os.MkdirAll(outputFolder+"/"+title+"/"+c.Hash.String()[0:8], os.ModePerm)
        os.WriteFile(outputFolder+"/"+title+"/"+c.Hash.String()[0:8]+"/index.html", buffer.Bytes(), os.ModePerm)

        buffer.Reset()
      }
      return nil
    })
  }

  return nil
})
```

## Libraries

* https://github.com/vbauerster/mpb

## References

* https://maelvls.dev/go111module-everywhere/
* https://github.com/flosch/pongo2/issues/68
* [Colors in `fmt`](https://golangbyexample.com/print-output-text-color-console/)
* [Versioning](https://stackoverflow.com/questions/11354518/application-auto-build-versioning)
* [Strings](https://dhdersch.github.io/golang/2016/01/23/golang-when-to-use-string-pointers.html)
* [getopts](https://pkg.go.dev/github.com/pborman/getopt)
