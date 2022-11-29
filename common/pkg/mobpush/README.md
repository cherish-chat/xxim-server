# [MobPush API for Golang](http://wiki.mob.com/mobpush-rest-api-接口文档/)

![image](https://github.com/MOBX/MOB-SMS-WEBAPI/blob/master/doc/images/logo.png)

**[MobPush API for Golang](http://wiki.mob.com/mobpush-rest-api-接口文档/)** 
为了帮助开发者更方便接入MobPush免费推送SDK，提供完整的API接口的Golang实现，包含设备操作相关接口、推送操作相关接口以及公共接口。

了解更多 [MobPush 免费推送SDK.](http://mobpush.mob.com)


## 优势

**免费使用**、**自定义UI**、**稳定服务**、**流程体验**、**数据同步**、**专业技术团队服务**

# 接口
* 推送接口:
	* 广播推送 pushAll
    * 别名推送 pushByAlias
    * 用户标签推送 pushByTags
    * Registration ID推送 pushByRids
    * 复杂地理位置推送 pushByAreas
    * 用户自定义配置推送 pushTaskV3         
* 推送任务详情查询接口
	* 查询推送任务详情(根据batchId) getPushByBatchId
	* 查询推送任务详情(根据workno) getPushByWorkno
* 推送任务的处理接口
    * 取消推送任务(根据workId) cancelPushTask
    * 替换推送任务(根据workId) replacePushTask
    * 撤回推送任务(根据workId) recallPushTask
* 查询推送统计接口
    * 根据推送任务id查询统计 getStatsByWorkId
    * 根据推送任务id批量查询统计 getStatsByWorkIds
    * 根据用户id查询统计 getStatsByWorkno
    * 按小时查询统计 getStatsByHour
    * 按日期查询统计 getStatsByDay
    * 根据id查询任务下发给设备的详情统计 getStatsByDevice
       


# 使用方式

 go get github.com/MobClub/mobpush-websdkv3-go   
 从开发者平台获取 appkey 和 appSecert

## 使用示例

#### 发送推送示例片段代码 
    
   
    workNo := "a7b9f1918c596ea"
    appkey,appSecert := "appkey","appSecert"
	client := NewPushClient(appkey,appSecert)
	res, err := client.PushByRids(workNo, "测试推送标题", "测试推送内容", []string{"rids1", "rids1"})

#### 统计查询示例片段代码

 
    workNo := "a7b9f1918c596ea"
    appkey,appSecert := "appkey","appSecert"
    client := NewPushClient(appkey,appSecert)
    res,err := client.GetStatsByWorkIds([]string{workNo})
    if err != nil{
    	t.Error("error")
    }
    
2.0.0版本新增设备管理和批量推送接口，以及定速推送功能
