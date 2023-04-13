# 使用方法

以Linux为例

同级目录下需要放置配置文件`setting.ini`
格式如下

```ini
[log]
level = Debug
# level = Info
# level = Warn
# level = Error
[person]
key = <换成你自己的key>
```

```bash
chmod a+x ./uploadFileForLinux
./uploadFileForLinux <你要上传的文件名,为了避免歧义,最好也放在同级目录,可以使用相对路径>
cat media.id

```