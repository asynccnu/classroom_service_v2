## classroom_service_v2

匣子空闲教室 Golang 版

### API Doc

[文档](./api.md)

### Env

```shell
export CCNUBOX_CLASSROOM_DB_NAME=''
export CCNUBOX_CLASSROOM_DB_URL=''
```

### Run

```
make
./main
```

### 数据导入

每学期根据选课手册（`*.xlsx`）导入教室数据，数据导入到一个新的数据库，即仍保存上学期的教室数据，不覆盖。

DB 命名格式：`classroom_<学年>_<学期>`。如 `classroom_20_21_1`，20-21学年第一学期

执行
```shell
make

# -p 后面跟文件路径
./main -p $FILE
```
