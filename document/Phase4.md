Phase4 数据持久化与好友关系

一.架构
1.使用MongoDB来存储用户的聊天记录
2.使用MySQL来存储用户的好友关系

二.数据库设计
1.用户表(已有)users:
id created_at updated_at deleted_at username password email avatar

2.好友关系表(relation):
owner_id target_id relation_type note(备注)
1代表好友,2代表黑名单

3.两表关系
(1).用户表和好友关系表是通过owner_id和target_id来关联的
(2).好友关系表中relation_type字段表示好友关系的类型,如好友、黑名单等
(3).好友关系表中Desc字段表示好友关系的备注,如"张三"、"李四的好友"等

4.消息表(message):
id session_id sender_id receiver_id content type timestamp

三.API
1.请求添加好友
POST /relation/add
{
    "target_id": ,
    "relation_type": ,
    "Desc": 
}

2.请求删除好友
POST /relation/delete
{
    "target_id": ,
}

3.请求获取好友列表
GET /relation/list

4. 请求获取聊天记录
GET /chat/history
{
    "owner_id": ,
    "target_id": ,
}
