# Multi user Manager for [shadowsocks](https://github.com/shadowsocks/shadowsocks-libev)
## How to Use
### BUild
```
git clone https://github.com/linexjlin/ss-web-manager.git
cd ss-web-manager
go build
```

### Prepare redis DB and run redis

```
cp -r redisDB.example redisDB 
cd redisDB
nohup redis-server redis.conf
```

### Run
```
cd ss-seb-manager
./ss-web-manager
###
```

### Create the first user the Admininstrator

Open the browser, open link http://127.0.0.1:8033/new_user 

