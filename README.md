# server-monitor
服务器监控 nginx influxdb 

## 功能
1. 根据指定log文件(比图nginx的log)将log做筛选并存入influxdb，线程安全
2. 使用grafana读取influxdb数据

## 开发日程
1. go 并发接口
2. nginx docker运行
3. influxDB 操作
4. grafana 使用
5. 整合 
