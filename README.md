# easyexcel

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
![https://opensource.org/licenses/MIT](https://img.shields.io/badge/license-MIT-green)
[![Go Reference](https://pkg.go.dev/badge/github.com/kexin8/easyexcel.svg)](https://pkg.go.dev/github.com/kexin8/easyexcel)

一个go语言开发的简单易用excel工具包

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Background
easyexcel旨在提供一个简单易用的excel工具包，让开发人员仅通过go结构体tags就能实现excel的导入与导出

## Install

```shell
go get github.com/kexin8/easyexcel
```

## Usage

### Import

```go
package main

import (
    "fmt"

    "github.com/kexin8/easyexcel"
)

type User struct {
	Name  string  `excel:"name:姓名"`
	Age   int     `excel:"name:年龄"`
	Sex   int     `excel:"name:性别;convertExp:0=男,1=女,2=未知"`
	Money float64 `excel:"name:金额"`
}

func main() {
    
    users, err := easyexcel.Import[User]("User.xlsx", easyexcel.NewOption())
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, user := users {
        fmt.Print("%+v",user)
    }
}

```

## API

## Maintainers

[@kexin8](https://github.com/https://github.com/kexin8)

## Contributing

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2023 zhengzi
