# frida

## 准备工作

1. 安装安卓模拟器
2. 安装frida

   - 电脑端

   ```shell
    pip3 install frida-tools
    pip3 install frida
    npm install frida -g
   ```

   - 手机端
    下载对应版本的server，用adb推送到手机/data/system/frida目录下

    ```shell
      adb push ~/frida-server /data/system/frida/frida-server
    ```

3. 启动server

  ```shell
    adb shell
    cd /data/system/frida
    chmod +x frida-server
    ./frida-server
  ```

4. 对app进行过反编译

## okhttp3 抓包

使用`https://github.com/siyujie/OkHttpLogger-Frida`进行抓包
需要先启动松下智能家电app，再执行`frida -U -l okhttp_poker.js -f cn.pana.caapp --no-pause`
另外，本地抓包时，不能执行各种对应的命令，因此，直接修改了okhttp_poker.js，改为直接调用hold方法

## MyLog重载

okhttp3的请求中，未找到获取和设置状态的接口，并且发现app中很多地方都用MyLog打印日志，因此，重写MyLog(见my_log.js)，输出日志内容。发现使用的webview，打开的控制页面。

```shell
frida -U -l my_log.js -f cn.pana.caapp --no-pause
```

## webview重载

使用webview.js，获取访问页面地址

```shell
frida -U -l webview.js -f cn.pana.caapp --no-pause
```