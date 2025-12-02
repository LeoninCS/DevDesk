# 🌟 DevDesk

[![Go Version](https://img.shields.io/badge/Go-1.24.5%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4%2B-42b883?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?logo=open-source-initiative&logoColor=white)](LICENSE)
[![Stars](https://img.shields.io/github/stars/LeoninCS/DevDesk?label=⭐%20Stars&logo=nil)](https://github.com/LeoninCS/DevDesk/stargazers)
[![Issues](https://img.shields.io/github/issues/LeoninCS/DevDesk?label=❗%20Issues)](https://github.com/LeoninCS/DevDesk/issues)






---
## [功能演示](https://b23.tv/ZZQ79DS)
## 效果演示
![alt text](/picture/image-1.png)
![alt text](/picture/image-2.png)
![alt text](/picture/image-3.png)
![alt text](/picture/image-4.png)
![alt text](/picture/image-5.png)
![alt text](/picture/image-6.png)
![alt text](/picture/image-7.png)
![alt text](/picture/image-8.png)

## 📝 前言
> **持续更新中，欢迎贡献代码！**  
> **如果你有什么好的 tips 或建议，欢迎在 Issues 中提出，或者联系作者。**

- **作者 QQ**：2329988157  
- **邮箱**：xianchaoqian@foxmail.com  

---

## 🚀 项目介绍
一个拥有诸多功能的工具库，帮助你解决瞬时问题。

---

## 🔧 功能介绍
- **Markdown Editor**：支持 Markdown 多人编辑和预览，支持导出多种格式。
- **HTTP Interface Testing**：类似 Postman 的接口测试工具，支持多请求工作区和 JSON 导入导出。
- **WorkPlan**：在这里存入你的计划，让你更好地管理时间。
- **CodeShare**：代码分享功能，帮助你分享代码片段。
- **HTMLHost**：提供在线 HTML 托管和预览服务。

---

## 💻 环境要求

> 现在项目推荐使用 **Docker 一键启动**，本机无需再单独配置 Go / Node 环境。

- **Docker**：20+（支持 `docker compose`）
- **Docker Compose**：v2+  
- （可选）本地开发环境：
  - **Go**：1.24 及以上  
  - **Vue**：3.4 及以上  

---

## ⚙️ 如何启动项目

### ✅ 方式一：使用 Docker 启动（推荐）

1. **确保已安装 Docker（带 compose 功能）**
2. **克隆项目到本地：**
   ```bash
   git clone https://github.com/LeoninCS/DevDesk.git
   ```

3. **进入项目目录：**

   ```bash
   cd DevDesk
   ```
4. **使用 Docker 启动服务：**

   ```bash
   # 如果你使用的是 Docker Compose 插件
   docker compose up -d

   # 如果你本机仍是老版本 docker-compose 命令，则：
   # docker-compose up -d
   ```
5. **查看容器是否正常运行：**

   ```bash
   docker ps
   ```
6. **停止项目：**

   ```bash
   docker compose down
   # 或
   # docker-compose down
   ```

> 前后端、依赖环境都会在容器里自动拉起，无需手动执行 `start.sh`。

---

### 🧑‍💻 方式二：本地开发启动（可选）

如果你想本地调试代码（非 Docker 模式），可以使用下面的方式：

1. **确保你已经安装了 Go 和 Node（用于 Vue），并且配置好环境变量。**
2. **克隆项目到本地：**

   ```bash
   git clone https://github.com/LeoninCS/DevDesk.git
   ```
3. **进入后端目录并安装依赖：**

   ```bash
   cd DevDesk
   cd backend
   go mod tidy
   ```
4. **shell脚本启动项目**
   ``` bash
   cd DevDesk
   chmod +x start.sh
   ./start.sh
   ```

---

## 🧪 测试方法

### 使用 Docker 启动后

1. 确保容器已正常运行（`docker ps` 可以看到相关容器）。
2. 打开浏览器，访问：

   ```text
   http://localhost:5173
   ```

   > 实际端口以你的 `docker-compose.yml` 映射为准，如果你改过端口，请按自己的配置访问。

### 本地开发模式

* 前端开发：访问 Vite 启动时输出的地址（通常是 `http://localhost:5173`）。
* 后端接口：根据你的后端监听端口（例如 `http://localhost:8080`）进行调试。

---

## 📈 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=LeoninCS/DevDesk\&type=Date)](https://star-history.com/#LeoninCS/DevDesk&Date)

