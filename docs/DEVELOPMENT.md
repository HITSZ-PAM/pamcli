# 基本要求

## 功能

本项目实现从PAM平台获取敏感信息，注入环境变量后供目标程序使用，避免了传统凭据管理硬编码或使用配置文件的凭据泄露风险。

## 开发平台

Golang 的包管理较高地依赖国外网络，国内使用建议配置GOPROXY使用国内镜像，例如

```
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

## 开发流程

Checkout - Pull Request模式

### Commit Message规范
使用AngularJS规范的 commit message
，格式如下
```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```
#### 允许的 `<type>`

- feat (feature)
- fix (bug fix)
- docs (documentation)
- style (formatting, missing semi colons, …)
- refactor (重构)
- test (when adding missing tests)
- chore (maintain)

#### 允许的 `<scope>`

Scope could be anything specifying place of the commit change. For example $location, $browser, $compile, $rootScope, ngHref, ngClick, ngView, etc...

#### `<subject>` 主要更改描述

use imperative, present tense: “change” not “changed” nor “changes”
don't capitalize first letter
no dot (.) at the end

#### Message body

just as in use imperative, present tense: “change” not “changed” nor “changes”
includes motivation for the change and contrasts with previous behavior

#### Message footer
- Breaking Change时说明具体的改变
- Close Issue
  
#### 样例

```
fix($compile): couple of unit tests for IE9

Older IEs serialize html uppercased, but IE9 does not...
Would be better to expect case insensitive, unfortunately jasmine does
not allow to user regexps for throw expectations.

Closes #392
Breaks foo.bar api, foo.baz should be used instead
```



# 项目开发

## 概要设计

## 详细设计

## 测试方法

Mock 关键api 并且使用黑盒测试。