## USER-REGISTRATION for MiniChatGroup
*Writer: Jonas Peng ; Contact: jonaspong@outlook.com ; Github: jonas-peng*

**请检查 README_CN.md 以获取中文版！！**

* You can get a text example, please check **Test_Examples.txt**, you can copy commands on the Linux terminal to test whether my micro-service works well.

* if you need develope another micro-service or a front end, please check **For_Other_Microservices_and_FrontEnd_Developement.txt** for details.
 
* Please call http://localhost:8000/q/swagger for a swagger ui to debug all functions !
 
### 1, Overview
This is a readme file for microservice as a part for a project named MiniChatGroup as a little try for the first cooperation of our team.

**It is built in Kratos with a docker file originally on Linux**

This microservice works to  store user information into a database, which is cretical for our project as a chat platform(not decentralized). I am responsible for this microservice. I write some APIs for calling to connect with other microsevices so that you can sign up, log in, delete and update your user account.

### 2, Development Environment Set-up:
*if you had not them yet, google to set them up all!!!*
#### 2.1, System Requirements
* OS: Linux(Ubuntu 24+ recommended)
* Go: version 1.22 or higher
* Database: PostgreSQL 16+
#### 2.2, Build Tools (Required for Code Generation)
*You need to install the following CLI tools to generate Proto codes and Dependency Injection:*

| Tool | Purpose | Install Command |
| :--- | :--- | :--- |
| **Kratos CLI** | Project management | `go install github.com/go-kratos/kratos/cmd/kratos/v2@latest` |
| **Wire** | Dependency Injection | `go install github.com/google/wire/cmd/wire@latest` |
| **Protoc** | Protobuf Compiler | `sudo apt install -y protobuf-compiler` |
| **Protoc Go Plugins** | Go code generation | (See setup script below) |

#### 2.3, Core Libraries (Go Modules)

*These dependencies are managed via `go.mod`. Run `go mod tidy` to install.*

- **Framework**: `github.com/go-kratos/kratos/v2`
- **ORM**: `gorm.io/gorm` & `gorm.io/driver/postgres`
- **Auth**: `github.com/golang-jwt/jwt/v5` & `golang.org/x/crypto`
- **Config**: `google.golang.org/protobuf`

#### 2.4, Tech Stacks

- **Language**: Go 1.22+
- **Framework**: [Kratos v2](https://github.com/go-kratos/kratos) (Microservice Framework)
- **Database**: PostgreSQL 16
- **ORM**: GORM (with AutoMigrate support)
- **Authentication**: JWT (JSON Web Token) + Bcrypt (Password Hashing)
- **Dependency Injection**: Google Wire
- **Protocol**: HTTP / gRPC

#### 2.5, Kratos Structure 

- `api/`: Protobuf definitions (API Contracts).
- `cmd/`: Application entry point (`main.go`, `wire_gen.go`).
- `internal/biz/`: **Domain Layer**. Core business logic & Domain Objects (DO).
- `internal/data/`: **Infrastructure Layer**. Database implementation & Persistent Objects (PO).
- `internal/service/`: **Interface Layer**. Handles HTTP/gRPC requests and DTO conversion.

### 3, Quick Start

#### 3.1, Prerequisites
Ensure you have the following installed:
- Go environment
- PostgreSQL (Create database: `benutzer_db`)

### 3.2, Run the Service
You can start the service using the Kratos CLI or standard Go commands:

```
# Install dependencies
go mod tidy

# Run the server
kratos run
# OR
cd cmd/user && go run .
```


### Conclusion
you can also use a one-click bash to set up all dependencies, on Linux, you need to open a terminal and copy bashes below into the terminal.


```bash
chmod +x ./ persetup.sh

./presetup.sh

```


Wish you good luck!

Please contact to me to my email or leave comments on Github.

Our team is a made of 4 third-year college students from **Wuhan University of Technology (Wuhan, China)**. The leader is NameIsNotName(Guo Zifan, the leader), Wang Haokun, Chen Xinqi and Jonas Peng(me).

A Multilingual version of the readme file would come soon. It is Jan 10, 2026 now.