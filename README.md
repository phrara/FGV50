# FGV50
> jy, v50; hy, v50!  

## Architecture
![archi](assets/arch.png)

## Web API

+ "/"
  + Method: **GET**
  + Param: null
+ "/cmd"
  + Method: **POST**
  + ReqParam: 
  ```go
    type Cmd struct {
	    CmdType string `json:"cmd_type"`
      
	    Url string `json:"url"`
	    Ip string `json:"ip"`
	    Port int `json:"port"`
	    NetworkSegment string `json:"network_segment"`
    }
  ```
  + RespParam:
  ```go
    type Result struct {
      Host     string `json:"host"`
      Port     int    `json:"port"`
     Protocol string `json:"protocol"`
      TTL      int    `json:"ttl"`
      Buf      []byte `json:"buf"`
      Type     string `json:"type"`
      IdBool   bool   `json:"id_bool"`
      IdString string `json:"idstring"`
      BString  string `json:"banner"`
    }
  ```



