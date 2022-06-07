[![CircleCI](https://circleci.com/gh/afreeorange/bockgo/tree/master.svg?style=svg)](https://circleci.com/gh/afreeorange/bockgo/tree/master)

## Development Notes

```bash
# Initialize this project
go mod init afreeorange/bock

# Remove unused mods
go mod tidy

# Remove a package
go get package@none

# Build
CGO_ENABLED=1 go build --tags "fts5" -o "dist/bock-$(uname)-$(uname -m)" .
```

### TODO

* [ ] Table of Contents
* [x] Configurable Syntax Highlight
* [x] Fix builds on cimg/go:1.18
* [ ] [Markdown highlight in Raw view](https://www.zupzup.org/go-markdown-syntax-highlight-chroma/)

### Versioning

Did this before I specified the version in `constants.go`. It has its advantages.

```golang
//go:embed VERSION
var b []byte
var version = string(b)
```

---

```bash
rm -rf ~/Desktop/temp/*; time go run --tags "fts5" . -a /Users/nikhilanand/personal/wiki.nikhil.io.articles -o /Users/nikhilanand/Desktop/temp
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

* A possible [progress bar](https://github.com/vbauerster/mpb).
* Hugo uses [afero](https://github.com/spf13/afero) as its filesystem abstraction layer. I have not needed it. Yet.
* [This is Commander](https://github.com/spf13/cobra) but for golang <3 Maybe not necessary here since the `flag` library in STDLIB has everything I need. But longopts are nice!

## References

* https://maelvls.dev/go111module-everywhere/
* https://github.com/flosch/pongo2/issues/68
* [Colors in `fmt`](https://golangbyexample.com/print-output-text-color-console/)
* [Versioning](https://stackoverflow.com/questions/11354518/application-auto-build-versioning)
* [Strings](https://dhdersch.github.io/golang/2016/01/23/golang-when-to-use-string-pointers.html)
* [getopts](https://pkg.go.dev/github.com/pborman/getopt)
* [Concurrency and Parallelism "Crash Course"](https://levelup.gitconnected.com/a-crash-course-on-concurrency-parallelism-in-go-8ea935c9b0f8)
* [Templates and Embed](https://philipptanlak.com/mastering-html-templates-in-go-the-fundamentals/#parsing-templates)
* [Recursive copying](https://blog.depa.do/post/copy-files-and-directories-in-go). I love that you have to implement quite a few things by hand in Golang!
* [Chroma/Pygment style reference](https://xyproto.github.io/splash/docs/all.html)
* [Enabling FTS5 with `go-sqlite`](https://github.com/mattn/go-sqlite3/issues/340)
* [Build Tags in Golang](https://www.digitalocean.com/community/tutorials/customizing-go-binaries-with-build-tags)
* [Go Routines Under the Hood](https://osmh.dev/posts/goroutines-under-the-hood)

### Books

* [Learning Go](https://miek.nl/go/learninggo.html)
* [Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE)
