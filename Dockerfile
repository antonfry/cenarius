FROM debian:bullseye-slim
WORKDIR /
COPY cmd/cenarius/cenarius-linux /cenarius
COPY conf/conf.toml /conf.toml
EXPOSE 8080
CMD ["/cenarius", "-m", "server", "-conf", "/conf.toml"]