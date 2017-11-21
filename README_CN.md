# 修复了原来glide的两个问题

## windows环境下经常在export dependencies时无故报错

主要是由于glide在最终拷贝文件时使用的命令是`move`，glide上已经有pull request了，我的版本把这个merge进来了

[#889](https://github.com/Masterminds/glide/pull/889/commits/cc37dc711a3191c2b91b01b9593c685660eeb9af)

原因就是windows下诡异的文件权限导致的。

如果大家在glide报错的时候看到的是乱码，请先在cmd里面执行`chcp 65001`，然后再看到的就不再是乱码了。

## 解决使用mirror拉取子包时出现export至错误的目录路径问题

问题描述：

由于GFW的缘故，`golang.org/x/...`的包全部拉不下了，幸好github上做了镜像在`github.com/golang/...`

我们只需要使用`glide mirror set`命令来设置镜像配置：
`glide mirror set https://golang.org/x/sys https://github.com/golang/sys`
这样我们就能愉快的玩耍了。

但耍着耍着会发现，有些第三方依赖需要`go get golang.org/x/sys/unix`，比如`gin`web框架，这时候就操蛋了。

有两种设置镜像的方式:

1. `glide mirror set https://golang.org/x/sys/unix https://github.com/golang/sys`

这个方式虽然能成功拉取，但是会发现export的目录完全错了，会复制到`golang.org/x/sys/unix`下面，将会看到
 `golang.org/x/sys/unix/unix`, `golang.org/x/sys/unix/plan9`, `golang.org/x/sys/unix/windows`，这样就完全错了。

2. `glide mirror set https://golang.org/x/sys/unix https://github.com/golang/sys/unix`

这种方式就会报`Cannot detect VCS`错误

原有的`glide`有解决方案，就是在glide.yaml中将`golang.org/x/sys/unix`加入到`ignore`中，然后手动补刀，将该依赖加进去。

针对这个问题，我的解决方案是修改`glide`的mirror命令的相关代码，添加`--base`参数，即用户可以指定fetch的包export到哪个目录。
这样针对需要依赖子包时就不会出现上述问题了。操作方法如下：

1. glide mirror set https://golang.org/x/sys/unix https://github.com/golang/sys --base golang.org/x/sys
2. glide get golang.org/x/sys/unix

这时候在`mirrors.yaml`文件中会看到

```yaml
repos:
- original: https://golang.org/x/crypto
  repo: https://github.com/golang/crypto
- original: https://golang.org/x/crypto/acme/autocert
  repo: https://github.com/golang/crypto
  base: golang.org/x/crypto
- original: https://golang.org/x/sys/unix
  repo: https://github.com/golang/sys
  base: golang.org/x/sys
```

同样使用`glide mirror list`命令也能看到`base`参数的输出。
