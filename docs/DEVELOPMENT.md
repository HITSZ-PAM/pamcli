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

### PAMClient
需要完成的功能
- 使用`client_id`和`client_secret`从中心端获取`token`
- 解析`pam://`并从中心端签出凭据
- 异常处理
  
**struct PAMConfig**
含有必要的pam连接和认证信息

**NewPAMClient** public 方法，用于从PAMConfig生成一个可用的PAMClient
如果为了提高扩展性，Client的创建可参考`secrethub-cli`使用工厂模式，但是目前单一的用途工厂模式不必要。参考`run.go`中的`NewRunCommand`函数。此函数需要执行PAMClient.init()

**PAMCLient.init()**
通过OAuth module获取Token，并存入PAMClient.Token中，配置好一个http.Client等待使用。

**PAMClient.resolve(url String)**
将url转换为真正的pam接口，并发起请求、解析返回的body。如果出现错误则处理错误。

### runCmd
参考`secrethub-cli`


### Cli Tools
使用[cobra](https://github.com/spf13/cobra)作为命令行`subcommand`和`flags`的管理库。

### 配置管理
使用[viper](https://github.com/spf13/viper)读取配置文件、命令行参数或者环境变量，主要用于pamcli自身的配置读取，对环境变量的替换依然使用`os.Environ`的方式。

## 测试方法

Mock 关键api 并且使用黑盒测试。