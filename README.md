**This is still a work in progress**
If you want to develop `go-isso`, fell free to [contact](https://lowentropy.me/about/#about-me) me.

# go-isso

![go-isso_s.png](https://i.loli.net/2018/10/16/5bc556ea1ae9a.png)

a commenting server similar to Disqus, while keeping completely API compatible with [isso](https://posativ.org/isso/)

see more doc in <https://go-isso.lowentropy.me/>

## why another isso

`isso` is good, but it's hard to be installed or customized.
What' more, the frontend part of `isso` use library that is no longer updated.

go-isso is distributed as a single binary, which means it can be installed and used easily.

## Roadmap

1. rewrite isso backend part <https://github.com/RayHY/go-isso/projects/1>
2. Pray that someone will help me rewrite the front part of isso.

## Getting Started

### Prerequisites

go-isso is commenting server written in Go language.

Make sure you have [go installed](https://golang.org/doc/install).

### Developing

Download the code: `go get -u github.com/RayHY/go-isso`

then `cd $GOPATH/src/github.com/RayHY/go-isso/cmd/go-isso/`

run `go build`

and play with it!

## Work in progress

**This is still a work in progress** so there's still bugs to iron out and as this
is my first project in Go the code could no doubt use an increase in quality,
but I'll be improving on it whenever I find the time. If you have any feedback
feel free to [raise an issue](https://github.com/jinxiapu/go-isso/issues)/[submit a PR](https://github.com/jinxiapu/go-isso/pulls).


## Contributing

I know NOTHING about javascript. I need someone to HELP ME!!!

If you want to develop `go-isso`, fell free to [contact](https://lowentropy.me/about/#about-me) me.

## Authors

* **Ray Wong** - *Initial work* - [LowEntropy](https://lowentropy.me)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details