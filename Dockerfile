FROM debian:stable-slim

# COPY source destination
COPY goserver /bin/goserver

ENV PORT=8888
CMD ["/bin/goserver"]

