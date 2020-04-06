# zcu102-fpga-device-plugin-demo
Try to design k8s device-plugin for embedded FPGA devices.

Note: The relevant files required by the build script are in the [original](https://github.com/Xilinx/FPGA_as_a_Service/tree/master/k8s-fpga-device-plugin) location

Thanks and Refer to：[FPGA_as_a_Service](https://github.com/xuhz/FPGA_as_a_Service)，and [this issue](https://github.com/Xilinx/FPGA_as_a_Service/issues/6) about zcu102 device-plugin.
## Design Flow:
1. Only modify the [fpga.go](https://github.com/Vae1997/zcu102-fpga-device-plugin-demo/blob/master/fpga.go) file according to the issue
2. Successfully build the binary file through the [build](https://github.com/Vae1997/zcu102-fpga-device-plugin-demo/blob/master/build) script
3. Build a device plugin image through [Dockerfile](https://github.com/Vae1997/zcu102-fpga-device-plugin-demo/blob/master/Dockerfile)
## Deployment and Verify
1. Start device plugin
```
kubectl create -f fpga-device-plugin.yml
```
The image ```vae2019/k8s_zcu102_fpga_plugin_test:latest```in fpga-device-plugin.yml is just build by [Dockerfile](https://github.com/Vae1997/zcu102-fpga-device-plugin-demo/blob/master/Dockerfile).

2. View the log output, the following shows that the plugin has been successfully deployed：
```
$kubectl logs -n kube-system fpga-device-plugin-daemonset-xxxx

time="2020-04-06T01:56:15Z" level=info msg="Starting FS watcher."
time="2020-04-06T01:56:15Z" level=info msg="Starting OS watcher."
time="2020-04-06T01:56:15Z" level=info msg="Starting to serve on /var/lib/kubelet/device-plugins/drm_minor-20191217-fpga.sock"
2020/04/06 01:56:15 grpc: Server.Serve failed to create ServerTransport:  connection error: desc = "transport: write unix /var/lib/kubelet/device-plugins/drm_minor-20191217-fpga.sock->@: write: broken pipe"
time="2020-04-06T01:56:15Z" level=info msg="Registered device plugin with Kubelet xilinx.com/fpga-drm_minor-20191217"
time="2020-04-06T01:56:15Z" level=info msg="Sending 1 device(s) [&Device{ID:a0000000.zyxclmm_drm,Health:Healthy,}] to kubelet"
time="2020-04-06T03:11:39Z" level=info msg="Receiving request a0000000.zyxclmm_drm"
```
There is an error in the log: ```write: broken pipe```, I don’t know the reason, and the log given [here](https://github.com/xuhz/FPGA_as_a_Service/tree/master/k8s-fpga-device-plugin/trunk#check-status-of-daemonset) also has this error.

3. Deploy pod to verify device plugin works
```
$kubectl create -f dp-pod.yaml
$kubectl exec -ti mypod -- ls /dev/dri
renderD128
$kubectl describe pod mypod
...
Limits:
      xilinx.com/fpga-drm_minor-20191217:  1
    Requests:
      xilinx.com/fpga-drm_minor-20191217:  1
...
```
The image ```vae2019/zcu102-fpga-verify:latest```in dp-pod.yaml is just build by [zcu102-fpga-verify/Dockerfile](https://github.com/Vae1997/zcu102-fpga-device-plugin-demo/blob/master/zcu102-fpga-verify/Dockerfile).
