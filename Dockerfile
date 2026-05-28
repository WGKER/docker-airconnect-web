FROM 1activegeek/docker-airconnect:1.9.3

# 复制已经编译好的二进制文件
COPY airconnect-webui /app/airconnect-webui

# 赋予执行权限
RUN chmod +x /app/airconnect-webui

# 暴露 WebUI 端口
EXPOSE 8087

# 同时启动 WebUI + 主程序
CMD ["/bin/sh", "-c", "/app/airconnect-webui & /airupnp-arm64-static ${AIRUPNP_VAR}"]
