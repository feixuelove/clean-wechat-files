# clean-wechat-files

自动删除 PC 端微信自动下载的大量文件、视频、图片等数据内容，解放一年几十 G 的空间占用。

>

### 该工具不会删除文字的聊天记录，请放心使用。请给个 Star 吧，非常感谢！

全面支持 windows mac linux系统中的所有微信版本。

配置简单
只需要配置 config.yaml 文件中的微信目录即可

```
path: "C:\\Users\\[用户名]\\Documents\\WeChat Files\\[微信ID]\\FileStorage\\MsgAttach"
days: 60
interval: 1h
log_file: "run.log"
```

path 为微信的文件目录
days 为删除多少天之前的数据
interval 每多少小时进行检查一次
log_file 运行日志

设置完毕后 windows 可以设置为系统服务
mac 或linux可以设置为开机服务

