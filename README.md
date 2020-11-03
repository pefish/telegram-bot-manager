# telegram-bot-manager

[![view examples](https://img.shields.io/badge/learn%20by-examples-0C8EC5.svg?style=for-the-badge&logo=go)](https://github.com/pefish/telegram-bot-manager)

Read this in other languages: [English](README.md), [简体中文](README_zh-cn.md)

telegram-bot-manager is a robot manager for telegram.

## Install

```
go get github.com/pefish/telegram-bot-manager/cmd/telegram-bot-manager
```

## Quick start

```shell
telegram-bot-manager --config="/path/to/sample.yaml"
```

Robot manager will reply all updates automatically according to the rules in `/path/to/sample.js`.

**/path/to/config.yaml**
```
token: "***"
commandsJsFile: "/path/to/sample.js"
```

**/path/to/sample.js**
```js
var commands = {
    "/test": function (args) {
        // console.log(args)
        return "test: " + JSON.stringify(args)
    },
    "/haha": function (args) {
        return "xixi"
    }
}
```


## Document

[doc](https://godoc.org/github.com/pefish/telegram-bot-manager)

## Contributing

1. Fork it
2. Download your fork to your PC
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Make changes and add them (`git add .`)
5. Commit your changes (`git commit -m 'Add some feature'`)
6. Push to the branch (`git push origin my-new-feature`)
7. Create new pull request

## Security Vulnerabilities

If you discover a security vulnerability, please send an e-mail to [pefish@qq.com](mailto:pefish@qq.com). All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the [Apache License](LICENSE).

