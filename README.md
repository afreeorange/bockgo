## Development Notes

```bash
# Initialize this project
go mod init afreeorange/bock

# Remove unused mods
go mod tidy

# Remove a package
go get package@none
```

I wanted to use Pongo2 but it causes this nonsense when I try to use `NewFSLoader`:

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x140d4e2]

goroutine 6 [running]:
github.com/flosch/pongo2/v4.(*Template).newBufferAndExecute(0x0, 0x1532e60?)
	/Users/nikhilanand/go/pkg/mod/github.com/flosch/pongo2/v4@v4.0.2/template.go:187 +0x22
github.com/flosch/pongo2/v4.(*Template).Execute(0x1577460?, 0xc0000c3ad0?)
	/Users/nikhilanand/go/pkg/mod/github.com/flosch/pongo2/v4@v4.0.2/template.go:231 +0x1d
main.render({0xc000014500, 0x24c5, 0x24c6}, {{0xc0002ca450, 0x24}, {0x0, 0x0}, {0xc000197c20, 0x0, 0x0}, ...})
	/Users/nikhilanand/personal/go-bock/render.go:28 +0x1c5
main.processArticle({0xc0002dc230, 0x4f}, {0x7ffeefbff6ed, 0x33}, {0x7ffeefbff724, 0x1f}, {0x16edbe0, 0xc0002da4e0}, 0x0?, 0xc000033500)
	/Users/nikhilanand/personal/go-bock/main.go:73 +0x570
created by main.process.func1
	/Users/nikhilanand/personal/go-bock/helpers.go:35 +0x238
exit status 2
```

## Libraries

* https://github.com/vbauerster/mpb

## References

* https://maelvls.dev/go111module-everywhere/
