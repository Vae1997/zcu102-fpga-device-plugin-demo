# Dockerfile of k8s-fpga-device-plugin
FROM ubuntu:16.04
# Build from build script
COPY bin/k8s-fpga-device-plugin /usr/local/bin/
# Run plugin in docker
ENTRYPOINT ["k8s-fpga-device-plugin"]
