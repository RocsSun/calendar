# calendar

#### 介绍
生成法定节假日和股市交易日历。调休请假有效时长计算。

#### 软件架构
编写了cli代码，可以直接cli使用，也可以使用源码。


#### 编译

1. `git clone https://github.com/RocsSun/calendar.git`
2. `cd calendar/cmd/calendar`
3. `go build .`
4. `./calendar -h`

#### 引入代码

1. `go get github.com/RocsSun/calendar`
2. `_ "github.com/RocsSun/calendar/cache"` # 添加本地缓存文件。
3. 核心逻辑在calendar下，导入对应的包使用。

#### 使用说明

1. 命令行提供work（工作日历），share（股票交易日历），effectTime（计算有效的请假或者调休时长）三个命令输入。
2. work和share 默认是当前年份的，两个支持参数-y xxxx输入年份。前提是国务院官网能查档查询年份的放假通知。2007年之前不支持。参数-j 文件名。将结果输出到指定的json文件。
3. effectTime计算有效的请假或者调休时长，需要输入--start开始时间。--end 结束时间。--amStart早上上班时间，--amEnd早上下班时间。--pmStart 下午上班时间。--pmEnd下午下班时间。其中start和end的时间格式是`YYYY-MM-DD hh:mm`，其余的为`hh:mm`。
