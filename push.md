参考微信消息推送设计

推送消息
DeviceId
Payload

Android
websocket 长连接
30秒ping一次

IOS推送
应用在前台：websocket 长连接
应用在后台：APNS

DeviceToken 管理
获取
存储
上报
更新
注销


上报
DeviceId 上报
消息到达上报


参考：
http://blog.csdn.net/ryantang03/article/details/8482259
http://blog.csdn.net/ryantang03/article/details/8540362
http://www.tuicool.com/articles/zq2Ez2
http://guafei.iteye.com/blog/1808445
