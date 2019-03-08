**This is still a work in progress**

New develop take place in branch `construction`.

If you want to develop `go-isso`, fell free to [contact](https://lowentropy.me/about/#about-me) me.

# go-isso

![go-isso_s.png](https://i.loli.net/2018/10/16/5bc556ea1ae9a.png)

a commenting server similar to Disqus, while keeping completely API compatible with [isso](https://posativ.org/isso/)

see more doc in <https://lowentropy.me/go-isso/>

## why another isso

`isso` is good, but it's hard to be installed or customized.
What' more, the frontend part of `isso` use library that is no longer updated.

`go-isso` is different from `isso`:

* Written in Go (Golang)
* Works with Sqlite3 but easy to add other database support.
* Doesn't use any ORM
* Doesn't use any complicated framework
* Use only modern vanilla Javascript (ES6 and Fetch API)
* Single binary compiled statically without dependency

### Why choose Golang as a programming language?

Go is probably the best choice for self-hosted software:

* Go is a simple programming language.
* Running code concurrently is part of the language.
* It’s faster than a scripting language like PHP or Python.
* The final application is a binary compiled statically without any dependency.
* You just need to drop the executable on your server to deploy the application.
* You don’t have to worry about what version of PHP/Python is installed on your machine.
* Packaging the software using RPM/Debian/Docker is straightforward.

## Roadmap

1. rewrite isso backend part <https://github.com/RayHY/go-isso/blob/master/roadmap.md>
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

## Thanks

[isso](https://posativ.org/isso/) & [isso's contributors](https://github.com/posativ/isso/graphs/contributors).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details