# go-isso

![go-isso_s.png](https://i.loli.net/2018/10/16/5bc556ea1ae9a.png)

a commenting server similar to Disqus, while keeping completely API compatible with [isso](https://posativ.org/isso/)

see more doc in <https://go-isso.lowentropy.me/>

## why another isso

`isso` is so good, but it is hard to install and customize. what' more, `isso`'s frontend part use something no longer updated which means it is hard to develop.

go-isso is distributed as a single binary, which means it can be installed and used easily.

## Getting Started

### Prerequisites

go-isso is commenting server written in Go language.

Make sure you have [go installed](https://golang.org/doc/install).

### Installing

Download the code: `go get -u github.com/jinxiapu/go-isso`

then `cd /cmd/go-isso/`

run `go build`

and play with it!

## Work in progress

**This is still a work in progress** so there's still bugs to iron out and as this
is my first project in Go the code could no doubt use an increase in quality,
but I'll be improving on it whenever I find the time. If you have any feedback
feel free to [raise an issue](https://github.com/jinxiapu/go-isso/issues)/[submit a PR](https://github.com/jinxiapu/go-isso/pulls).

## Built With

* [Go Standard library](https://golang.org/pkg/)

* [logrus](https://github.com/sirupsen/logrus)

## Contributing

I know nothing about javascript. I need someone to 

## Authors

* **Ray Wong** - *Initial work* - [LowEntropy](https://lowentropy.me)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details