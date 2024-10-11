CREATE TABLE user
(
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deletedAt DATETIME DEFAULT NULL COMMENT '删除时间',
    userId    BIGINT(20) PRIMARY KEY COMMENT '用户id',
    userName  VARCHAR(15) UNIQUE NOT NULL COMMENT '用户名',
    gender    CHAR(1)  DEFAULT '男' COMMENT '性别',
    birth     DATE     DEFAULT NULL COMMENT '生日',
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
    postId      BIGINT(20)   NOT NULL COMMENT '帖子id',
    userId      BIGINT(20)   NOT NULL COMMENT '用户id',
    content     VARCHAR(511) NOT NULL COMMENT '评论内容',
    star        BIGINT     DEFAULT 0 COMMENT '点赞数',
    reply       BIGINT     DEFAULT 0 COMMENT '回复数',
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
    PRIMARY KEY (communityId, leaderId) -- communityId 和 leaderId 共同为主键
) COMMENT '社区表';

CREATE TABLE `group_messages`
(
    `chatId`      bigint(20) PRIMARY KEY,
    `sendTime`    datetime     NOT NULL COMMENT '消息发送时间',
    `content`     varchar(255) NOT NULL COMMENT '内容',
    `sdUserId`    bigint       NOT NULL COMMENT '发送方id',
    `communityId` bigint       NOT NULL COMMENT '社区id'
) COMMENT ='群聊表';


CREATE TABLE `private_messages`
(
    `private_message_id` bigint       NOT NULL COMMENT '私信id',
    `sender_id`          bigint       NOT NULL COMMENT '发送方',
    `recipient_id`       bigint       NOT NULL COMMENT '接受方',
    `content`            varchar(255) NOT NULL COMMENT '内容',
    `state`              boolean DEFAULT FALSE COMMENT '是否已读',
    `send_time`          datetime     NOT NULL COMMENT '发送时间',
    PRIMARY KEY (`private_message_id`),
    INDEX (`recipient_id`)
) comment '私信表';



create table `private_chat`
(
    private_chat_id      bigint comment '会话id',
    user1_id             bigint(20) comment '用户1的id',
    user2_id             bigint(20) comment '用户2的id',
    last_message_content varchar(255) comment '最后一条消息内容',
    last_message_time    datetime comment '最后一条消息时间',
    primary key (private_chat_id)
) comment '会话表';

create table `manager_system_notice`
(
    system_notice_id bigint,
    tile             varchar(30) comment '标题',
    content          varchar(255) comment '内容',
    type             varchar(10) comment '类型',        -- single 单用户  all 全体用户
    state            boolean comment '状态',            -- 是否被拉取过
    notice_user_id   bigint(20) comment '通知的用户id', -- 如果为全体用户则为0
    manager_id       bigint(20) comment '管理员用户id',
    publish_time     datetime comment '发布时间',
    primary key (system_notice_id)
) comment '管理员系统通知表';


create table `user_system_notice`
(
    user_notice_id   bigint,
    state            boolean comment '是否已读',
    system_notice_id bigint comment '系统通知id',
    notice_user_id   bigint comment '接受通知的用户id',
    pull_time        datetime comment '拉取时间',
    primary key (user_notice_id),
    index (system_notice_id)
) comment '用户系统通知表';


create table `like_remind`
(
    like_remind_id bigint,
    like_id        bigint comment '点赞源id',-- 评论id 帖子id
    like_type      varchar(10) comment '点赞源类型', -- "comment","post"
    like_content   varchar(255) comment '点赞源内容',
    url            varchar(255) comment '点赞源链接',
    state          boolean comment '是否已读',
    sender_id      boolean comment '点赞人id',
    recipient_id   boolean comment '接受通知的人的id',
    remind_time    datetime comment '提醒时间',
    primary key (like_remind_id)

) comment '点赞提醒表';



create table `mention_remind`
(
    mention_remind_id bigint,
    mention_id        bigint comment '@源id',-- 评论id 帖子id
    mention_type      varchar(10) comment '@源类型', -- "comment","post"
    mention_content   varchar(255) comment '@源内容',
    url               varchar(255) comment '@源链接',
    state             boolean comment '是否已读',
    sender_id         boolean comment '@人id',
    recipient_id      boolean comment '接受通知的人的id',
    remind_time       datetime comment '提醒时间',
    primary key (mention_remind_id)
) comment '@提醒表';


create table `reply_remind`
(
    reply_remind_id bigint,
    reply_id        bigint comment '回复源id',-- 评论id 帖子id
    reply_type      varchar(10) comment '回复源类型', -- "comment","post"
    reply_content   varchar(255) comment '回复源内容',
    url             varchar(255) comment '回复源链接',
    state           boolean comment '是否已读',
    sender_id       boolean comment '回复人id',
    recipient_id    boolean comment '接受通知的人的id',
    remind_time     datetime comment '提醒时间',
    primary key (reply_remind_id)
) comment '回复提醒表';
