FROM 1activegeek/docker-airconnect:1.9.3

# 1. 复制你的 webui 二进制
COPY airconnect-webui /app/
RUN chmod +x /app/airconnect-webui

# 2. 暴露端口
EXPOSE 8087

# 3. 给 s6 加一个长期服务（longrun），随 /init 自动启动
# 路径：/etc/s6-overlay/s6-rc.d/<服务名>/
RUN mkdir -p /etc/s6-overlay/s6-rc.d/webui \
    && mkdir -p /etc/s6-overlay/s6-rc.d/user/contents.d

# 启动脚本：用 s6 的环境执行 webui
RUN echo -e '#!/command/with-contenv bash\n\nexec /app/airconnect-webui' > /etc/s6-overlay/s6-rc.d/webui/run \
    && chmod +x /etc/s6-overlay/s6-rc.d/webui/run

# 服务类型：longrun = 长期运行、崩溃自动重启
RUN echo "longrun" > /etc/s6-overlay/s6-rc.d/webui/type

# 加入 user 启动组 → 容器启动 /init 时自动拉起 webui
RUN echo "webui" > /etc/s6-overlay/s6-rc.d/user/contents.d/webui

ENTRYPOINT ["/init"]
