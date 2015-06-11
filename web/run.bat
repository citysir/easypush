set GOPATH=D:\my\gopath
cd %~dp0
go run api.auth.go api.messages.go api.message.go api.syncid.go api.uploaddid.go api.internal.push.go conf.go constant.go global.go main.go node.go redis.go protocol.go rpcclient.go rpcserver.go signalwatch.go stringutil.go usertoken.go

pause