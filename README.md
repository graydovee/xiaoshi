# 小诗AI机器人

小诗是一个基于大语言模型的AI聊天机器人，支持OneBot协议，具备上下文记忆能力，并通过MCP插件机制实现功能扩展。适用于QQ等IM平台的智能对话、自动化助手、虚拟女仆等场景。

---

## 功能特性
- **OneBot协议支持**：可无缝对接QQ等IM平台，实现群聊/私聊AI对话。
- **大模型驱动**：支持任意兼容OpenAI协议并支持Function Call的大模型，具备丰富的自然语言理解与生成能力（经测试，DeepSeekV3在有较长System Prompt的情况下Function Call能力较差，建议使用openAI模型）。
- **上下文记忆**：可配置记忆长度与时效，支持多会话上下文管理。
- **MCP插件扩展**：通过MCP协议可动态加载第三方插件（如B站搜索、终端指令等），让AI能力无限拓展。
- **高度可定制**：支持自定义角色设定、系统提示词、模型参数等。

---

## 快速开始

### 1. 构建镜像

```shell
make docker-build
```

### 2. 配置文件

请参考 `config.example.yaml`，复制为 `config.yaml` 并根据实际需求修改：

```shell
cp config.example.yaml config.yaml
```

主要配置项说明见下文。

### 3. 运行容器

```shell
docker run -itd \
  -v ${PWD}/config.yaml:/etc/xiaoshi/config.yaml \
  -e CONFIG_FILE_PATH=/etc/xiaoshi/config.yaml \
  graydovee/xiaoshi:latest
```

如需本地开发调试，也可直接编译运行：

```shell
make build
CONFIG_FILE_PATH=config.yaml ./bin/xiaoshi
```

---

## 说明

### 配置文件字段详解

| 字段 | 类型 | 默认值 | 说明 |
| ---- | ---- | ------ | ---- |
| memory.messageLimit | int | 32 | 单会话最大记忆消息数 |
| memory.expireSeconds | int | 600 | 记忆过期时间（秒） |
| oneBot.id | int | - | 机器人QQ号（需与OneBot服务一致） |
| oneBot.ws.addr | string | localhost | OneBot WebSocket服务地址 |
| oneBot.ws.port | int | 3001 | OneBot WebSocket服务端口 |
| oneBot.ws.token | string | "" | OneBot鉴权Token（如无可留空） |
| mcp.systemPrompt | string | 见示例 | AI角色设定与行为约束（支持多行） |
| mcp.llm.apiKey | string | - | 大模型API密钥（如OpenAI Key） |
| mcp.llm.baseUrl | string | https://api.openai.com/v1 | 大模型API地址 |
| mcp.llm.model | string | gpt-4.1 | 使用的大模型名称 |
| mcp.mcpServers | map | - | MCP插件服务配置（可多个） |
| mcp.mcpServers.*.command | string | - | 本地插件启动命令 |
| mcp.mcpServers.*.args | list | - | 启动命令参数列表 |
| mcp.mcpServers.*.url | string | - | 远程插件服务地址（如为本地插件可不填） |
| mcp.mcpServers.*.description | string | - | 插件功能描述 |
| mcp.mcpServers.*.env | map | - | 启动插件时的环境变量 |
| mcp.mcpServers.*.disabled | bool | false | 是否禁用该插件 |
| mcp.mcpServers.*.autoApprove | list | [] | 自动批准的操作列表 |

如需更多配置示例和详细注释，请参考 `config.example.yaml` 文件。

---

## MCP插件扩展

本项目支持通过MCP协议动态加载插件，扩展AI能力。

### mcp插件示例
- **bilibili-search**：B站视频搜索
- **terminal-stdio**：本地终端指令执行

### 添加自定义插件
1. 按MCP协议开发插件，或者从mcp社区获取现成插件
2. 在`config.yaml`的`mcp.mcpServers`中添加插件配置
3. 重启机器人即可自动加载

---

## 贡献与支持

欢迎提交PR、Issue或自定义插件！

- 项目主页：[https://git.graydove.cn/graydove/xiaoshi](https://git.graydove.cn/graydove/xiaoshi)
- 主要作者：graydovee
- License: MIT

