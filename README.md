# GTI-OBD-Scanner

📹️ 基于 Zap 封装的日志组件

> 项目名的灵感，来自汽车中的故障诊断仪，该部件可以监测车辆的各种性能数据，并在出现故障时生成相应的故障代码

## 输出预览
![1a4e2837e9a653a449db014bb1921737](https://github.com/user-attachments/assets/0c3570d8-a9aa-4b92-882f-fa57d359186c)

## 依赖
| 技术         | 说明   | 仓库地址                                           |
|------------|------|------------------------------------------------|
| zap        | 日志框架 | https://go.uber.org/zap                        |
| rotatelogs | 日志切割 | https://github.com/lestrrat-go/file-rotatelogs |

## 功能
- 分级存储
- 调用文件全路径
- 堆栈信息携带行号
- 终端日志级别颜色
- 自定义日志留存时间
- 自定义日志最大尺寸
- 支持文本或 JSON 输出格式
- 日志切割（根据年、月、日、时、分和秒）

## 下一期功能规划
- 压缩日志文件
