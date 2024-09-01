CREATE TABLE user
(
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deletedAt DATETIME DEFAULT NULL COMMENT '删除时间',
    userId    BIGINT(20) PRIMARY KEY COMMENT '用户id',
    userName  VARCHAR(15) UNIQUE NOT NULL COMMENT '用户名',
    gender    CHAR(1)  DEFAULT '男' COMMENT '性别',
    birth     DATE DEFAULT NULL COMMENT '生日',
    sign      VARCHAR(50) COMMENT '个性签名',
    grade     SMALLINT DEFAULT 0 COMMENT '等级',-- 检查等级在0~100的范围，用代码实现CHECK (userGrade BETWEEN 0 AND 100 )
    exp       SMALLINT DEFAULT 0 COMMENT '经验',--  检查等级在0~499的范围，用代码实现CHECK(userExp BETWEEN 0 AND 499)
    img       VARCHAR(255) /* DEFAULT 默认头像*/ COMMENT '用户头像'
) COMMENT '用户表';

CREATE TABLE userLogin
(
    createdAt DATETIME                                 DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME                                 DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deletedAt DATETIME                                 DEFAULT NULL COMMENT '删除时间',
    userId    BIGINT(20) PRIMARY KEY COMMENT '用户id',
    userName  VARCHAR(15) UNIQUE NOT NULL COMMENT '用户名',
    email     varchar(30) UNIQUE                       DEFAULT NULL COMMENT '邮箱',
    phone     char(11) UNIQUE CHECK (phone LIKE '1%' ) DEFAULT NULL COMMENT '手机号',
    password  varchar(85)        NOT NULL COMMENT '密码'
) COMMENT '登录授权表';

CREATE TABLE userFavor
(
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deletedAt DATETIME DEFAULT NULL COMMENT '删除时间',
    userId    BIGINT(20) NOT NULL COMMENT '用户id',
    beFavorId BIGINT(20) NOT NULL COMMENT '关注用户id',--  不能关注自己，用代码实现CHECK ( userId!=beFavorId )
    isFriend  BOOL     DEFAULT 0 COMMENT '是否互关',
    PRIMARY KEY (userId, beFavorId)                    -- userId 和 beFavorId 共同为主键
) COMMENT '用户关注表';

CREATE TABLE userCollect
(
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deletedAt DATETIME DEFAULT NULL COMMENT '删除时间',
    userId    BIGINT(20) NOT NULL COMMENT '用户id',
    postId    BIGINT(20) NOT NULL COMMENT '收藏帖子id',
    PRIMARY KEY (userId, postId) -- userId 和 postId 共同为主键
) COMMENT '用户收藏表';

CREATE TABLE post
(
    createdAt   DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    deletedAt   DATETIME DEFAULT NULL COMMENT '删除时间',
    postId      BIGINT(20) PRIMARY KEY COMMENT '帖子id',
    userId      BIGINT(20) UNIQUE NOT NULL COMMENT '用户id',
    collection  BIGINT   DEFAULT 0 COMMENT '收藏数',
    star        BIGINT   DEFAULT 0 COMMENT '点赞数',
    comment     BIGINT   DEFAULT 0 COMMENT '评论数',
    content     VARCHAR(2047)     NOT NULL COMMENT '帖子内容',
    title       VARCHAR(20)       NOT NULL COMMENT '标题',
    isScan      BOOL     DEFAULT 1 COMMENT '可见性',
    communityId BIGINT(20)        NOT NULL COMMENT '社区id'
) COMMENT '帖子表';

CREATE TABLE postComment
(
    createdAt   DATETIME   DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    deletedAt   DATETIME   DEFAULT NULL COMMENT '删除时间',
    commentId   BIGINT(20) PRIMARY KEY COMMENT '评论id',
    postId      BIGINT(20)  NOT NULL COMMENT '帖子id',
    userId      BIGINT(20)        NOT NULL COMMENT '用户id',
    content     VARCHAR(511)      NOT NULL COMMENT '评论内容',
    star        BIGINT   DEFAULT 0 COMMENT '点赞数',
    reply       BIGINT   DEFAULT 0 COMMENT '回复数',
    beCommentId BIGINT(20) DEFAULT NULL COMMENT '关联评论id'
) COMMENT '评论表';

CREATE TABLE community
(
    createdAt     DATETIME    DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt     DATETIME    DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deletedAt     DATETIME    DEFAULT NULL COMMENT '删除时间',
    communityId   BIGINT(20) UNIQUE  NOT NULL COMMENT '社区id',
    description   VARCHAR(50) DEFAULT NULL COMMENT '简介',
    communityName VARCHAR(15) UNIQUE NOT NULL COMMENT '社区名',
    member        BIGINT      DEFAULT 0 COMMENT '成员数',
    leaderId      BIGINT(20)         NOT NULL COMMENT '社区主持',
    manageId      BIGINT(20)         NOT NULL COMMENT '管理员',
    img           VARCHAR(255) /* DEFAULT 默认头像*/ COMMENT '社区头像',
    PRIMARY KEY (communityId, leaderId)                     -- communityId 和 leaderId 共同为主键
) COMMENT '社区表';

CREATE TABLE `private_messages` (
    `deletedAt` datetime DEFAULT NULL COMMENT '删除时间',
    `chatId` bigint NOT NULL COMMENT '消息id',
    `type`   varchar(20)  NOT NULL  COMMENT  '类型', -- "like", "reply", "system", "mention", "whisper"
    `sendTime`  datetime   NOT NULL  COMMENT  '发送时间',
    `sdUserId` bigint NOT NULL COMMENT '发送方',
    `acUserId` bigint NOT NULL COMMENT '接受方',
    `content` varchar(255) NOT NULL COMMENT '内容',
    `isRead`  boolean    DEFAULT FALSE  COMMENT '是否已读',
     PRIMARY KEY (`chatId`),
     INDEX (acUserId)
)  COMMENT='消息表';

CREATE TABLE `group_messages` (
    `chatId`     bigint(20) PRIMARY KEY ,
    `sendTime` datetime NOT NULL COMMENT '消息发送时间',
    `content` varchar(255) NOT NULL COMMENT '内容',
    `sdUserId` bigint NOT NULL COMMENT '发送方id',
    `communityId` bigint NOT NULL COMMENT '社区id'
)  COMMENT='群聊表';