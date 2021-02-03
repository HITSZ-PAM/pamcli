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

Checkout 一个新的分支，完成对代码的修改以后创建一个Pull Request

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
程序的实现分为几个部分
- 初始化：创建需要的各个对象。创建 httpClient 并做好认证。
- 读取全部环境变量并选出含有`pam://`的进行接下来的操作
- 将这些环境变量逐一交由 Client 解析
- 将新的环境变量与不需要修改的环境变量注入目标程序的运行环境
- 创建子进程
- 设定好stdin pass through, signal pass through
- 运行子进程
- 等待子进程退出后退出

## 详细设计
### 初始化
- 读取环境变量，获取认证用的`client_id`和`client_secret`
- 使用`client_id`和`client_secret`创建Client，例如命名为pamClient

### pamClient
需要完成的功能
- 使用`client_id`和`client_secret`从中心端获取`token`
- 解析`pam://`并从中心端签出凭据
- 异常处理
  
为了扩展性，Client的创建可参考`secrethub-cli`使用工厂模式，但是目前单一的用途工厂模式不必要。建议参考`run.go`中的`NewRunCommand`函数。

### 其他部分
基本与`secrethub-cli`一致，按实际需求修改。

## 测试方法

Mock 关键api 并且使用黑盒测试。