# telegram-bot-manager

[![view examples](https://img.shields.io/badge/learn%20by-examples-0C8EC5.svg?style=for-the-badge&logo=go)](https://github.com/pefish/telegram-bot-manager)

Read this in other languages: [English](README.md), [简体中文](README_zh-cn.md)

telegram-bot-manager is a robot manager for telegram.

## 安装

```
go get github.com/pefish/telegram-bot-manager/cmd/telegram-bot-manager
```

## 快速开始

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

## 文档

[doc](https://godoc.org/github.com/pefish/XXX)

## 贡献代码（非常欢迎）

1. Fork 仓库
2. 代码 Clone 到你本机
3. 创建feature分支 (`git checkout -b my-new-feature`)
4. 编写代码然后 Add 代码 (`git add .`)
5. Commin 代码 (`git commit -m 'Add some feature'`)
6. Push 代码 (`git push origin my-new-feature`)
7. 提交pull request

## 安全漏洞

如果你发现了一个安全漏洞，请发送邮件到[pefish@qq.com](mailto:pefish@qq.com)。

## 授权许可

[Apache License](LICENSE).

