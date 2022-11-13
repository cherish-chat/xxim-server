# http api gateway

## 1. 请求以及响应

### 1.1 请求

#### 1.1.1 请求头中的参数

| 参数名 | 说明                                                 | 是否必须 | 默认值                    |
| :--- |:---------------------------------------------------|:-----|:-----------------------|
| Content-Type | 请求的数据类型                                            | 是    | application/x-protobuf |
| Platform | 平台: ios,android,ipad,androidpad,windows,macos,wasm | 是    |                        |
| DeviceId | 设备id                                               | 是    |                        |
| UserId | 用户id                                               | 否    |                        |
| Token | 用户token                                            | 否    |                        |
| AppVersion | app版本号                                             | 否    |                        |
| OsVersion | 系统版本号                                              | 否    |                        |
| DeviceModel | 设备型号                                               | 否    |                        |
| Language | 语言 | 否    | zh_CN                  |

### 1.1.2 请求方法和路径

    请求方法固定为POST
    请求路径格式为: /{version}/{service}/[isNeedAuth]/{method}

#### 1.1.3 请求体

    请求体中的数据类型为`protobuf`，每一个接口定义的请求体都包含一个 requester 字段，用户不需要设置其中的值，只需要保证该字段的值不为null即可。

```protobuf
message Requester {
  string id = 1;
  string token = 2;
  string appVersion = 3;
  string ip = 4;

  // header
  string ua = 5;
  string osVersion = 6;
  string platform = 7;
  string deviceModel = 8;
  string deviceId = 9;
  string language = 10; // 系统语言：zh_CN en_US
}
```

### 1.2 响应

#### 1.2.1 http状态码

| 状态码 | 说明         |
| :--- |:-----------|
| 200  | 请求成功       |
| 400  | 请求参数错误     |
| 401  | 未授权         |
| 500 | 服务器内部错误 |

#### 1.2.2 响应体

    响应体中的数据类型为`protobuf`，以下是proto定义的消息体，其中`data`字段也是一个protobuf消息体，具体的消息体定义在各个接口的proto文件中。

```protobuf
message CommonResp {
  enum Code {
    Success = 0;

    InternalError = 2; // 内部错误
    RequestError = 3;  // 请求错误
    ToastError = 5;    // toast 错误 只有 message
    AlertError = 7;    // alert 错误 只有一个确定按钮
    RetryError = 8;    // alert 错误 有一个取消按钮 和一个重试按钮
  }
  Code code = 1;
  optional string msg = 2;
  bytes data = 3;
}
```

## 2. 通用proto消息体

### 2.1 通用

#### 2.1.1 CommonResp: 响应消息体

```protobuf
message CommonResp {
  enum Code {
    Success = 0;

    UnknownError = 1;  // 未知 error
    InternalError = 2; // 内部错误
    RequestError = 3;  // 请求错误
    AuthError = 4;     // 鉴权错误 // 应该退出登录
    ToastError = 5;    // toast 错误 只有 message
    AlertError = 7;    // alert 错误 只有一个确定按钮
    RetryError = 8;    // alert 错误 有一个取消按钮 和一个重试按钮
  }
  Code code = 1;
  optional string msg = 2;
  bytes data = 3;
}
```

#### 2.1.2 Requester: 请求消息体

```protobuf
message Requester {
  string id = 1;
  string token = 2;
  string appVersion = 3;
  string ip = 4;

  // header
  string ua = 5;
  string osVersion = 6;
  string platform = 7;
  string deviceModel = 8;
  string deviceId = 9;
  string language = 10; // 系统语言：zh-CN en-US
}
```

#### 2.1.3 IpRegion: ip地区

```protobuf
message IpRegion {
  string country = 1;
  string province = 2;
  string city = 3;
  string isp = 4;
}
```

#### 2.1.4 Page: 分页

```protobuf
message Page {
  int32 page = 1;
  int32 size = 2;
}
``` 

### 2.2 用户

#### 2.2.1 Xb: 用户性别

```protobuf
enum XB {
  UnknownXB = 0;
  Male = 1;
  Female = 2;
}
```

#### 2.2.2 Constellation: 星座

```protobuf
enum Constellation {
  UnknownConstellation = 0;
  Aries = 1;
  Taurus = 2;
  Gemini = 3;
  Cancer = 4;
  Leo = 5;
  Virgo = 6;
  Libra = 7;
  Scorpio = 8;
  Sagittarius = 9;
  Capricorn = 10;
  Aquarius = 11;
  Pisces = 12;
}
```

#### 2.2.3 BirthdayInfo: 生日信息

```protobuf
message BirthdayInfo {
  int32 year = 1;
  int32 month = 2;
  int32 day = 3;
  int32 age = 4;
  Constellation constellation = 5;
}
```

#### 2.2.4 LevelInfo: 等级信息

```protobuf
message LevelInfo {
  int32 level = 1;
  int32 exp = 2;
  // 下一级所需经验
  int32 nextLevelExp = 3;
}
```

#### 2.2.5 UserBaseInfo: 用户基本信息

```protobuf
message UserBaseInfo {
  string id = 1;
  string nickname = 2;
  string avatar = 3;
  XB xb = 4;
  // 生日信息
  BirthdayInfo birthday = 5;
  // 最后一次连接 ip所在地
  IpRegion ipRegion = 6;
}
```

### 2.3 im

#### 2.3.1 MsgNotifyOpt: 消息通知选项

```protobuf
//MsgNotifyOpt 消息通知选项
message MsgNotifyOpt {
  bool preview = 1; // 是否预览
  bool sound = 2; // 是否声音
  string soundName = 3; // 声音名称
  bool vibrate = 4; // 是否震动
}
```

### 2.3 消息

#### 2.3.1 ConvType: 会话类型

```protobuf
enum ConvType {
  SINGLE = 0; // 单聊
  GROUP = 1; // 群聊
}
```

#### 2.3.2 ContentType: 消息类型

```protobuf

enum ContentType {
  UNKNOWN = 0;
  TYPING = 1; // 正在输入
  READ = 2; // 已读
  REVOKE = 3; // 撤回

  TEXT = 11; // 文本
  IMAGE = 12; // 图片
  AUDIO = 13; // 语音
  VIDEO = 14; // 视频
  FILE = 15; // 文件
  LOCATION = 16; // 位置
  CARD = 17; // 名片
  MERGE = 18; // 合并
  EMOJI = 19; // 表情
  COMMAND = 20; // 命令

  CUSTOM = 100; // 自定义消息
}
```

#### 2.3.3 MsgData: 消息体

```protobuf

message MsgData {
  message OfflinePush {
    string title = 1;
    string content = 2;
    string payload = 3;
  }
  message Options {
    // 是否需要离线推送
    bool offlinePush = 1;
    // 服务端是否需要保存消息
    bool storageForServer = 2;
    // 客户端是否需要保存消息
    bool storageForClient = 3;
    // 消息是否需要计入未读数
    bool unreadCount = 4;
    // 是否需要解密 （端对端加密技术，服务端无法解密）
    bool needDecrypt = 5;
    // 是否需要重新渲染会话
    bool updateConv = 6;
  }
  message Receiver {
    optional string userId = 1; // 单聊时为对方的userId
    optional string groupId = 2; // 群聊时为群组id
  }
  string clientMsgId = 1;
  string serverMsgId = 2;
  string clientTime = 3;
  string serverTime = 4;

  string sender = 11; // 发送者id
  string senderInfo = 12; // 发送者信息
  string senderConvInfo = 13; // 发送者在会话中的信息

  Receiver receiver = 21; // 接收者id (单聊时为对方id, 群聊时为群id)
  string convId = 22; // 会话id
  repeated string atUsers = 23;   // 强提醒用户id列表 用户不在线时，会收到离线推送，除非用户屏蔽了该会话 如果需要提醒所有人，可以传入"all"

  ContentType contentType = 31; // 消息内容类型
  bytes content = 32; // 消息内容
  string seq = 33; // 消息序号 会话内唯一且递增

  Options options = 41; // 消息选项
  OfflinePush offlinePush = 42; // 离线推送

  bytes ext = 100;
}
```

#### 2.3.4 MsgDataList: 消息列表

```protobuf
message MsgDataList {
  repeated MsgData msgDataList = 1;
}
```

## 3. 接口

### 3.1 用户模块

#### 3.1.1 Login: 登录

- 请求地址：`/v1/user/white/login`
- 请求体：

```protobuf
message LoginReq {
  Requester requester = 1;
  string id = 2; // 用户id 只能是英文和数字_，长度为6-20
  string password = 3; // 密码 // md5 数据库中会存入该值加盐后的值
}
```

- 响应体：

```protobuf
message LoginResp {
  CommonResp commonResp = 1;
  // 是否是新用户
  bool isNewUser = 2;
  // token
  string token = 3; // 如果是新用户，token为空
}
```

#### 3.1.2 ConfirmRegister: 确认注册

- 请求地址：`/v1/user/white/confirmRegister`
- 请求体：

```protobuf
message ConfirmRegisterReq {
  Requester requester = 1;
  string id = 2; // 用户id 只能是英文和数字_，长度为6-20
  string password = 3; // 密码 // md5 数据库中会存入该值加盐后的值
}
```

- 响应体：

```protobuf
message ConfirmRegisterResp {
  CommonResp commonResp = 1;
  string token = 2;
}
```

#### 3.1.3 SearchUsersByKeyword: 使用关键词搜索用户

- 请求地址：`/v1/user/searchUsersByKeyword`
- 请求体：

```protobuf
message SearchUsersByKeywordReq {
  Requester requester = 1;
  string keyword = 2;
}
```

- 响应体：

```protobuf
message SearchUsersByKeywordResp {
  CommonResp commonResp = 1;
  repeated UserBaseInfo users = 2;
}
```

#### 3.1.4 GetUserHome: 获取用户主页信息

- 请求地址：`/v1/user/getUserHome`
- 请求体：

```protobuf
message GetUserHomeReq {
  Requester requester = 1;
  string id = 2;
}
```

- 响应体：

```protobuf
message GetUserHomeResp {
  CommonResp commonResp = 1;
  string id = 2;
  string nickname = 3;
  string avatar = 4;
  XB xb = 5;
  BirthdayInfo birthday = 6;
  IpRegion ipRegion = 7;
  // 个性签名
  string signature = 8;
  // 等级信息
  LevelInfo levelInfo = 9;
}
```

#### 3.1.5 GetUserSettings: 获取用户设置

- 请求地址：`/v1/user/getUserSettings`
- 请求体：

```protobuf
enum UserSettingKey {
  HowToAddFriend = 0; // 如何添加好友
  HowToAddFriend_NeedAnswerQuestionCorrectly_Question = 1; // 如何添加好友 需要回答的问题
  HowToAddFriend_NeedAnswerQuestionCorrectly_Answer = 2; // 如何添加好友 需要回答的问题的答案
}
message GetUserSettingsReq {
  Requester requester = 1;
  repeated UserSettingKey keys = 2;
}
```

- 响应体：

```protobuf
message UserSetting {
  string userId = 1;
  UserSettingKey key = 2;
  string value = 3;
}

message GetUserSettingsResp {
  CommonResp commonResp = 1;
  map<int32, UserSetting> settings = 2;
}
```

#### 3.1.6 SetUserSettings: 更新用户设置

- 请求地址：`/v1/user/setUserSettings`
- 请求体：

```protobuf
message SetUserSettingsReq {
  Requester requester = 1;
  repeated UserSetting settings = 2;
}
```

- 响应体：

```protobuf
message SetUserSettingsResp {
  CommonResp commonResp = 1;
}
```

### 3.2 关系模块

#### 3.2.1 RequestAddFriend: 请求添加好友

- 请求地址：`/v1/relation/requestAddFriend`
- 请求体：

```protobuf
message RequestAddFriendReq {
  Requester requester = 1;
  string to = 2;
  // 附加消息
  string message = 3;
}
```

- 响应体：

```protobuf
message RequestAddFriendResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.2 AcceptAddFriend: 接受添加好友请求

- 请求地址：`/v1/relation/acceptAddFriend`
- 请求体：

```protobuf
message AcceptAddFriendReq {
  Requester requester = 1;
  string applyUserId = 2; // 申请人id
  optional string requestId = 3; // 申请id
}
```

- 响应体：

```protobuf
message AcceptAddFriendResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.3 RejectAddFriend: 拒绝添加好友请求

- 请求地址：`/v1/relation/rejectAddFriend`
- 请求体：

```protobuf
message RejectAddFriendReq {
  Requester requester = 1;
  string applyUserId = 2; // 申请人id
  string requestId = 3; // 申请id
  bool block = 4; // 是否拉黑
}
```

- 响应体：

```protobuf
message RejectAddFriendResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.4 BlockUser: 拉黑用户

- 请求地址：`/v1/relation/blockUser`
- 请求体：

```protobuf
message BlockUserReq {
  Requester requester = 1;
  string userId = 2;
}
```

- 响应体：

```protobuf
message BlockUserResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.5 DeleteBlockUser: 取消拉黑用户

- 请求地址：`/v1/relation/deleteBlockUser`
- 请求体：

```protobuf
message DeleteBlockUserReq {
  Requester requester = 1;
  string userId = 2;
}
```

- 响应体：

```protobuf
message DeleteBlockUserResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.6 DeleteFriend: 删除好友

- 请求地址：`/v1/relation/deleteFriend`
- 请求体：

```protobuf
message DeleteFriendReq {
  Requester requester = 1;
  string userId = 2;
  bool block = 3; // 是否拉黑
}
```

- 响应体：

```protobuf
message DeleteFriendResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.7 SetSingleChatSetting: 设置单聊设置

- 请求地址：`/v1/relation/setSingleChatSetting`
- 请求体：

```protobuf
message SetSingleChatSettingReq {
  Requester requester = 1;
  string userId = 2;
  // 置顶选项
  optional bool top = 3;
  // 单聊的备注
  optional string remark = 4;
  // 免打扰选项
  optional bool disturb = 11;
}
```

- 响应体：

```protobuf
message SetSingleChatSettingResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.8 GetSingleChatSetting: 获取单聊设置

- 请求地址：`/v1/relation/getSingleChatSetting`
- 请求体：

```protobuf
message GetSingleChatSettingReq {
  Requester requester = 1;
  string userId = 2;
}
```

- 响应体：

```protobuf
message GetSingleChatSettingResp {
  CommonResp commonResp = 1;
  // 置顶选项
  bool top = 2;
  // 单聊的备注
  string remark = 3;
  // 免打扰选项
  bool disturb = 4;
}
```

#### 3.2.9 SetSingleMsgNotifyOpt: 设置单聊消息通知选项

- 请求地址：`/v1/relation/setSingleMsgNotifyOpt`
- 请求体：

```protobuf
message SetSingleMsgNotifyOptReq {
  Requester requester = 1;
  // 群ID
  string userId = 2;
  // 消息通知类型
  MsgNotifyOpt opt = 3;
}
```

- 响应体：

```protobuf
message SetSingleMsgNotifyOptResp {
  CommonResp commonResp = 1;
}
```

#### 3.2.10 GetSingleMsgNotifyOpt: 获取单聊消息通知选项

- 请求地址：`/v1/relation/getSingleMsgNotifyOpt`
- 请求体：

```protobuf
message GetSingleMsgNotifyOptReq {
  Requester requester = 1;
  // 群ID
  string userId = 2;
}
```

- 响应体：

```protobuf
message GetSingleMsgNotifyOptResp {
  CommonResp commonResp = 1;
  // 消息通知类型
  MsgNotifyOpt opt = 2;
}
```

#### 3.2.11 GetFriendList: 获取好友列表

- 请求地址：`/v1/relation/getFriendList`
- 请求体：

```protobuf
message GetFriendListReq {
  Requester requester = 1;
  // 分页
  Page page = 2;
  enum Opt {
    WithBaseInfo = 0; // 带用户的基本信息
    OnlyId = 1; // 只有用户id
  }
  Opt opt = 10;
}
```

- 响应体：

```protobuf
message GetFriendListResp {
  CommonResp commonResp = 1;
  repeated string ids = 2;
  map<string, UserBaseInfo> userMap = 3;
}
```

### 3.3 群组

#### 3.3.1 CreateGroup: 创建群组

- 请求地址：`/v1/group/createGroup`
- 请求体：

```protobuf
message CreateGroupReq {
  Requester requester = 1;
  // 拉人进群
  repeated string members = 2;
  // 群名称(可选参数)
  optional string name = 3;
  // 群头像(可选参数)
  optional string avatar = 4;
}
```

- 响应体：

```protobuf
message CreateGroupResp {
  CommonResp commonResp = 1;
  // 群ID
  optional string groupId = 2;
}
```

### 3.4 消息

#### 3.4.1 SendMsg: 发送消息

- 请求地址：`/v1/msg/sendMsg`
- 请求体：

```protobuf
message SendMsgListReq {
  Requester requester = 1;
  repeated MsgData msgDataList = 2;
  // options
  // 1. 延迟时间（秒） 不得大于 864000秒 也就是10天 只有开启了Pulsar的延迟消息功能才有效
  optional int32 deliverAfter = 11;
}
```

- 响应体：

```protobuf
message SendMsgListResp {
  CommonResp commonResp = 1;
}
```

