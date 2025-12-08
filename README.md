# xiaoshi-kovi-plugin

一个基于 [Kovi](https://github.com/Threkork/kovi) 框架的智能聊天机器人插件，具备长短期记忆、RAG 检索增强、MCP 工具调用等能力。

## ✨ 功能特性

### 🤖 LLM 集成
- 支持 OpenAI 兼容 API（包括 DeepSeek、硅基流动等）
- 灵活的参数配置：temperature、top_p、max_tokens 等

![LLM 对话示例](doc/images/hello.png)

### 🧠 记忆管理
- **短期记忆**：基于会话的上下文记忆，支持历史条数和超时设置
- **长期记忆（RAG）**：基于 PostgreSQL + pgvector 的向量数据库，支持语义检索

### 📊 智能记忆评估
- 自动评估对话价值（0-100分）
- 根据评分智能决定记忆保留时长：
  - 0-25 分：不保存（噪音/废弃对话）
  - 26-60 分：保留 1 周（短期任务）
  - 61-85 分：保留 1 月（中期状态/偏好）
  - 86-100 分：永久保存（核心身份信息）

### 🔧 MCP 工具调用
- 支持 Model Context Protocol (MCP)
- 支持多种传输方式：`stdio`、`sse`、`streamable-http`
- 可接入搜索、文件系统等外部工具

**MCP 工具调用示例：**

![search1api MCP 示例](doc/images/search1api.png)

![trendRadar MCP 示例](doc/images/trendRadar.png)

### 💬 消息处理
- 私聊：直接回复用户消息
- 群聊：通过 @机器人 触发回复

## 📝 配置说明

插件配置文件位于 Kovi 的 data 目录下：`data/xiaoshi-kovi-plugin/config.json`

### config.json 配置示例

```json
{
  "llm": {
    "model": "gpt-4",
    "url": "https://api.openai.com/v1",
    "apikey": "your-api-key",
    "temperature": null,
    "top_p": null,
    "max_tokens": null,
    "presence_penalty": null,
    "frequency_penalty": null
  },
  "db": {
    "postgres": {
      "host": "localhost",
      "port": "5432",
      "username": "postgres",
      "password": "your-password",
      "database": "xiaoshi",
      "vector": {
        "lists": 100
      }
    }
  },
  "memory": {
    "history_limit": 20,
    "history_timeout": 600,
    "prompt": "你是一个可爱的虚拟女仆",
    "rag": {
      "enabled": true,
      "embedding": {
        "model": "Qwen/Qwen3-Embedding-0.6B",
        "url": "https://api.siliconflow.cn/v1/embeddings",
        "apikey": "your-embedding-api-key"
      },
      "top_n": 3,
      "window_size": 2,
      "max_memory_tokens": 1000,
      "cleanup_days": 30,
      "memory_evaluation": {
        "enabled": true,
        "model": "deepseek-chat",
        "url": "https://api.deepseek.com/v1",
        "apikey": "your-evaluation-api-key"
      }
    }
  },
  "mcp": {
    "enabled": true,
    "path": "mcp.json",
    "max_tool_iterations": 10
  }
}
```

### 配置项说明

| 配置项 | 说明 |
|--------|------|
| `llm.model` | 主对话模型名称 |
| `llm.url` | LLM API 地址 |
| `llm.apikey` | LLM API 密钥 |
| `db.postgres.*` | PostgreSQL 数据库连接配置（需安装 pgvector 扩展） |
| `memory.history_limit` | 短期记忆保留的最大消息条数 |
| `memory.history_timeout` | 短期记忆超时时间（秒） |
| `memory.prompt` | 系统提示词 |
| `memory.rag.enabled` | 是否启用 RAG 长期记忆 |
| `memory.rag.embedding.*` | 向量嵌入模型配置 |
| `memory.rag.top_n` | RAG 检索返回的锚点数量 |
| `memory.rag.window_size` | 每个锚点的上下文窗口大小 |
| `memory.rag.memory_evaluation.*` | 记忆评估模型配置 |
| `mcp.enabled` | 是否启用 MCP 工具调用 |
| `mcp.path` | MCP 配置文件路径（相对于 config.json） |
| `mcp.max_tool_iterations` | 单次对话最大工具调用轮数 |

### mcp.json 配置示例

MCP 配置文件用于定义可用的外部工具服务：

```json
{
  "mcpServers": {
    "search1api": {
      "transport": "stdio",
      "command": "npx",
      "args": ["-y", "search1api-mcp"],
      "env": {
        "SEARCH1API_KEY": "your-search-api-key"
      }
    },
    "trendRadar": {
      "transport": "streamable-http",
      "url": "https://your-mcp-server.com/mcp"
    }
  }
}
```

**支持的传输方式：**
- `stdio`: 通过子进程标准输入输出通信
- `sse`: Server-Sent Events
- `streamable-http`: HTTP 流式传输

## 🗄️ 数据库准备

RAG 功能需要 PostgreSQL 数据库并安装 pgvector 扩展：

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

插件首次启动时会自动创建所需的数据表。

## 📜 License

MIT

