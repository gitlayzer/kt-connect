Ktctl Completion
---

为`ktctl`工具开启命令和参数的`Tab键`自动补全功能。包含3个子命令：

- bash        生成Bash的自动补全配置
- zsh         生成Zsh的自动补全配置
- fish        生成Fish的自动补全配置

详细使用方法如下：

<!-- tabs:start -->

#### ** bash **

MacOS系统：

```bash
ktctl completion bash > /usr/local/etc/bash_completion.d/ktctl
```

Linux系统：

```bash
ktctl completion bash > /etc/bash_completion.d/ktctl
```

对命令执行后所有新打开的窗口生效。

#### ** zsh **

MacOS系统：

```bash
ktctl completion zsh > /usr/local/share/zsh/site-functions/_ktctl
```

Linux系统：

```bash
ktctl completion zsh > "${fpath[1]}/_ktctl"
```

对命令执行后所有新打开的窗口生效。

#### ** fish **

仅对当前Shell窗口生效：

```bash
ktctl completion fish | source
```

对所有新的Shell窗口生效：

```bash
ktctl completion fish > ~/.config/fish/completions/ktctl.fish
```

<!-- tabs:end -->

`ktctl`的自动补全功能包括"命令补全"和"参数补全"，举例如下（其中`<tab>`为按下键盘上的TAB键）：
 
- 命令补全：输入`ktctl ex<tab>`，将自动补全为`ktctl exchange`
- 参数补全：输入`ktctl connect --m<tab>`，将自动补全为`ktctl connect --mode`

当存在多种匹配的补全结果时，可通过连续按Tab键，在多种结果之间切换。
