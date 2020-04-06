package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"wssocket/impl"
)

var upgrade=websocket.Upgrader{
		//允许跨域
	CheckOrigin: func(r *http.Request)bool{
			return true
	},
}

func wsHandle(w http.ResponseWriter, r *http.Request){
	var(
		wsConn *websocket.Conn
		err error
		data []byte
		conn *impl.Connection
	)
	if wsConn,err= upgrade.Upgrade(w,r,nil);err!=nil{
		return
	}
	if conn,err=impl.InitConnection(wsConn);err!=nil{
		goto ERR
	}

	go func(){
		var(
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heatbeat")); err != nil {
				return
			}
			time.Sleep(time.Second)
		}
	}()

	for{
		if data,err=conn.Readmessage();err!=nil{
			goto ERR
		}
		if err=conn.WriteMessage(data);err!=nil{
			goto ERR
		}
	}
	ERR:
		conn.Close()
}

func main(){
	http.HandleFunc("/ws",wsHandle)
	http.ListenAndServe("0.0.0.0:8888",nil)
}