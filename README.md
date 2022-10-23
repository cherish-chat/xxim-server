# xxim-server

<p align="center">
<img align="center" width="150px" src="https://raw.githubusercontent.com/cherish-chat/xx-doc/main/doc/images/xxim.jpg">
</p>

xxim-server是一个功能超多的开箱即用的IM服务器。它的诞生是为了让每个人都能使用自己的IM APP，不需要再使用令人不爽的第三方IM APP。xxim-server是一个开源项目，欢迎大家一起来完善它。

<div align=center>

[![Go](https://github.com/cherish-chat/xxim-server/workflows/Go/badge.svg?branch=master)](https://github.com/cherish-chat/xxim-server/actions)
[![codecov](https://codecov.io/gh/cherish-chat/xxim-server/branch/master/graph/badge.svg)](https://codecov.io/gh/cherish-chat/xxim-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/cherish-chat/xxim-server)](https://goreportcard.com/report/github.com/cherish-chat/xxim-server)
[![Release](https://img.shields.io/github/v/release/cherish-chat/xxim-server.svg?style=flat-square)](https://github.com/cherish-chat/xxim-server)
[![Go Reference](https://pkg.go.dev/badge/github.com/cherish-chat/xxim-server.svg)](https://pkg.go.dev/github.com/cherish-chat/xxim-server)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>


## 🤷‍ xxim 介绍
简体中文 | [English](README-EN.md)

xxim-server代码不复杂，im大多逻辑都在于客户端，所以xxim-server只是一个简单的im服务器，但它具备了一个IM应有的全部功能。

#### 包括但不限于：

* 发送消息（可定时的、可群发），包括：文本、图片、语音、视频、文件、位置、名片、撤回、转发、@、表情、对方正在输入、自定义消息等
* 按需拉取离线消息，且没有消息数量/天数限制
* 已读管理（对方是否已读、群内已读的成员）
* 音视频通话、IOS支持`Callkit`
* 当用户不在线时，通过厂商推送（极光、腾讯、Mob）将消息推送给用户
* 用户的每个会话都可以设置消息接收选项（接收、不接收、接收但不提醒）
* 不限人数的群聊

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/architecture-en.png" alt="Architecture" width="1500" />

## xxim的背景

2022年初，我们公司的社交产品需要一个IM，但是我们不想使用第三方IM，所以我们自己开发了一个IM，但是我们发现开发一个IM并不容易，所以我们决定开源出来，让更多的人能够使用自己的IM。

* 服务端使用 Go 语言开发
    * 高性能
    * 简单语法，易于维护代码
    * 部署简单
    * 服务器资源占用少
* 客户端使用 flutter 开发
    * 跨平台、一套代码多端运行
    * 支持原生系统调用，性能强大
    * 界面美观、交互流畅

## xxim的设计原则

通过im服务器，我们希望解决以下问题：

* 聊天受监控
* 消息漫游天数有限制
* 群聊人数有限制
* 消息占用磁盘空间过大

## xxim-server 架构

<img width="1067" alt="image" src="https://raw.githubusercontent.com/cherish-chat/xx-doc/main/doc/images/xxim.jpg">

## Benchmark

![benchmark](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/benchmark.png)

[Checkout the test code](https://github.com/smallnest/go-web-framework-benchmark)

## Documents

* [Documents](https://go-zero.dev/)
* [Rapid development of microservice systems](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)
* [Rapid development of microservice systems - multiple RPCs](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)
* [Examples](https://github.com/zeromicro/zero-examples)

##  Chat group

Join the chat via https://discord.gg/4JQvC5A4Fe

##  Cloud Native Landscape

<p float="left">
<img src="https://landscape.cncf.io/images/left-logo.svg" width="150"/>&nbsp;&nbsp;&nbsp;
<img src="https://landscape.cncf.io/images/right-logo.svg" width="200"/>
</p>

go-zero enlisted in the [CNCF Cloud Native Landscape](https://landscape.cncf.io/?selected=go-zero).

## Give a Star! ⭐

If you like or are using this project to learn or start your solution, please give it a star. Thanks!

[![Star History Chart](https://api.star-history.com/svg?repos=zeromicro/go-zero&type=Date)](#go-zero)

## Buy me a coffee

<a href="https://www.buymeacoffee.com/kevwan" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
