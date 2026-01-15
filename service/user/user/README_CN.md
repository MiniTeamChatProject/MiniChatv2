用户注册微服务 - MiniChatGroup 项目

作者：Jonas Peng；联系方式：jonaspong@outlook.com；Github：jonas-peng

请检查 README_CN.md 以获取中文版！！

    您可以获取文本示例，请查看 Test_Examples.txt，您可以在 Linux 终端上复制命令来测试我的微服务是否正常工作。

    如果您需要开发另一个微服务或前端，请查看 For_Other_Microservices_and_FrontEnd_Developement.txt 了解详细信息。

1，概述

这是为 MiniChatGroup 项目的微服务编写的 readme 文件，作为我们团队首次合作的小尝试。

它最初在 Linux 上使用 Kratos 构建，并附带 Dockerfile

该微服务负责将用户信息存储到数据库中，这对我们作为聊天平台（非去中心化）的项目至关重要。我负责这个微服务。我编写了一些 API 供其他微服务调用，以便您可以注册、登录、删除和更新您的用户帐户。
2，开发环境设置：

如果您还没有安装它们，请通过谷歌搜索来设置它们！！！
2.1，系统要求

    操作系统：Linux（推荐 Ubuntu 24+）

    Go：版本 1.22 或更高

    数据库：PostgreSQL 16+

2.2，构建工具（代码生成所需）

您需要安装以下 CLI 工具来生成 Proto 代码和依赖注入：
工具	用途	安装命令
Kratos CLI	项目管理	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
Wire	依赖注入	go install github.com/google/wire/cmd/wire@latest
Protoc	Protobuf 编译器	sudo apt install -y protobuf-compiler
Protoc Go 插件	Go 代码生成	（参见下面的设置脚本）
2.3，核心库（Go 模块）

这些依赖通过 go.mod 管理。运行 go mod tidy 来安装。

    框架：github.com/go-kratos/kratos/v2

    ORM：gorm.io/gorm 和 gorm.io/driver/postgres

    认证：github.com/golang-jwt/jwt/v5 和 golang.org/x/crypto

    配置：google.golang.org/protobuf

2.4，技术栈

    语言：Go 1.22+

    框架：Kratos v2（微服务框架）

    数据库：PostgreSQL 16

    ORM：GORM（支持 AutoMigrate）

    认证：JWT（JSON Web Token）+ Bcrypt（密码哈希）

    依赖注入：Google Wire

    协议：HTTP / gRPC

2.5，Kratos 结构

    api/：Protobuf 定义（API 合约）。

    cmd/：应用程序入口点（main.go、wire_gen.go）。

    internal/biz/：领域层。核心业务逻辑和领域对象（DO）。

    internal/data/：基础设施层。数据库实现和持久化对象（PO）。

    internal/service/：接口层。处理 HTTP/gRPC 请求和 DTO 转换。

3，快速开始
3.1，先决条件

确保您已安装以下内容：

    Go 环境

    PostgreSQL（创建数据库：benutzer_db）

3.2，运行服务

您可以使用 Kratos CLI 或标准的 Go 命令启动服务：
text

# 安装依赖
go mod tidy

# 运行服务器
kratos run
# 或者
cd cmd/user && go run .

结语

您也可以使用一键 bash 脚本来设置所有依赖项，在 Linux 上，您需要打开终端并将以下 bash 脚本复制到终端中。
bash

chmod +x ./ persetup.sh

./presetup.sh

祝您好运！

请通过电子邮件联系我或在 Github 上留言。

我们的团队由武汉理工大学（中国武汉）的四名大三学生组成。领导是 NameIsNotName（郭梓凡，队长）、王浩坤、陈欣琪和 Jonas Peng（我）。

多语言版本的 readme 文件即将推出。现在是 2026 年 1 月 10 日。