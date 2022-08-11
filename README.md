# Gonews -- a RSS/Atom feed static site generator
## Note
This is purely a hobby project and it is meant to be used for reference only. It is not meant to be used in production.
## Built with
- [Cobra](https://github.com/spf13/cobra)
### Progress
- Unstable, not guaranteed to work with any particular source.

## Usage
Either ```git clone``` and build or use:
```bash
go install github.com/nguyendhst/gonews@latest
```
and then run:
```bash
gonews generate
```

## Existing problems
- Duplicate .html file name in the same directory are not handled yet.
- Not yet tested on Windows.

## References
- [RSS 2.0 Specification](https://cyber.harvard.edu/rss/rss.html)
- [Atom 1.0 Specification](https://tools.ietf.org/html/rfc4287)
- [gofeed](https://github.com/mmcdole/gofeed)
- [goread](https://github.com/bake/goread)
- [Escaping HTML](https://h1z3y3.me/posts/go-html-template-script-unescape/)
