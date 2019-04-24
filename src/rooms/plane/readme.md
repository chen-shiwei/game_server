# 猜飞机

## 1.匹配操作
	CMD_PLAN_MATCH uint8 = 50
    typ  PACKET_JSON = 4
    
    + 请求数据:
        {
            "plane1":{
                "rotate":1,
                "centerid":1,
                "head":1,
                "body":1,
            },
            "plane2":{
                "rotate":1,
                "centerid":1,
                "head":1,
                "body":1,
            },
            "plane3":{
                "rotate":1,
                "centerid":1,
                "head":1,
                "body":1,
            }
        }
    + 返回数据
        {
            "code":1001,
            "data":{
                "matchStatus":1
            }                        
        }
    注意：status 1 进入匹配队列，等待匹配 2.匹配成功，进入战斗 
## 2.战斗
    CMD_PLAN_HIT   uint8 = 51
    typ PACKET_JSON = 4
    请求数据：
    {   
        "hitId":1 //int类型
    }
    返回数据：
    {
        "code":1001,
        "data":{
            "hitResult":1,  //1 击中头部 2 击中身体 3 没有击中
            "destroyPlaneIndex":1 //0没有打中飞机 1代表plane1 以此类推 
        }        
    }


## 3.notice: 
    1.匹配成功 进入房间通知
    CMD_PLAN_MATCH_NOTICE uint8 = 52
    {   
        "code":1001,
        "data":{
            "player":{
                "userId":1,
                "uname":"111",
                "head":"111"
            },
            "playerData":{
                "plane1":{
                    "rotate":1,
                    "centerId":1,
                    "headId":1,
                    "bodyPos":1,
                },
                "plane2":{
                    "rotate":1,
                    "centerId":1,
                    "headId":1,
                    "bodyPos":1,
                },
                "plane3":{
                    "rotate":1,
                    "centerId":1,
                    "headId":1,
                    "bodyPos":1,
                }
            },
            "playerHistoryData":{
                "winNum":11,
                "allNum":11
            }
        }
    }
    2.打击通知
	CMD_PLAN_HIT_NOTICE   uint8 = 54
    {
        "code":1001,
        "data":{
            "hitId":1,
            "hitResult":1,  //1 击中头部 2 击中身体 3 没有击中
            "destroyPlaneIndex":1 //0没有打中飞机 1代表plane1 以此类推 
        }
    }    
    
## 4.取消匹配
	CMD_PLAN_CANCEL_MATCH uint8 = 53    
    请求

    响应
    {
        "code":1001,
        "data":null
    }
## 5.结算推送
	CMD_PLAN_OVER         uint8 = 55
    响应
    {
        "code":1001,
        "data":{
            "players"[
                {
                    "user":{
                        "userId":1,
                        "uname":"222",
                        "head":""
                    },
                    "iswin":1,
                    "getscore":1
                },
            ]
        }
     
    }    
## 错误返回
    {
        code:601
    }
    codes:
        + 600 server处理错误
        + 601 请求参数错误
        + 1301 用户不在匹配中
        + 1302 禁止操作



    
