# 参考
  http://www.blogjava.net/yongboy/archive/2014/03/05/410636.html
  http://www.w3cschool.cc/redis/redis-sorted-sets.html

# 接口
- auth(username, did, password)
  token, node

- messages(token, topic, synckey)
  synckeys, max_synckey

- sync(token, topic, synckey)
  success or failure