# 构建阶段
FROM rust:1.75-slim as builder

# 安装必要的构建依赖
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 复制依赖文件
COPY Cargo.toml Cargo.lock ./
COPY plugins/kovi-plugin-xiaoshi/Cargo.toml ./plugins/kovi-plugin-xiaoshi/

# 创建虚拟项目以缓存依赖
RUN mkdir -p src plugins/kovi-plugin-xiaoshi/src && \
    echo "fn main() {}" > src/main.rs && \
    echo "" > plugins/kovi-plugin-xiaoshi/src/lib.rs && \
    cargo build --release && \
    rm -rf src plugins/kovi-plugin-xiaoshi/src

# 复制源代码
COPY . .

# 构建应用
RUN cargo build --release

# 运行阶段
FROM debian:bookworm-slim

# 安装运行时依赖
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/target/release/kovi /usr/local/bin/kovi

# 复制配置文件（如果需要）
COPY kovi.conf.toml kovi.plugin.toml ./

# 设置环境变量
ENV RUST_LOG=info

# 暴露端口（根据实际需要调整）
EXPOSE 8080

# 运行应用
CMD ["kovi"]

