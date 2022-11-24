# Cluster Autoscaler for Hetzner

The cluster autosclaer for [Hetzner Cloud](https://www.hetzner.com/cloud). Can scale multiple node pools. Node pools defined by [labels hetzner cloud](https://docs.hetzner.cloud/#overview-labels). Cluster Autosclaer will run as a `Deployment` in your cluster.

## How it works
Cluster Autoscaler create and delete nodes if necessary using Api Hetzner Cloud. You can find more detailed information in [FAQ](../../FAQ.md). Work is supported only within one [project](https://wiki.hetzner.de/index.php/CloudServer/en#What_are_projects.2C_and_how_can_I_use_them.3F). Each node in the pool is assigned its own [label](https://docs.hetzner.cloud/#overview-labels), it matches the name of the pool. The labels determine whether a node belongs to a pool.

The `TemplateNodeInfo` method is supported, by which the best pool for scaling is selected (when choosing the corresponding` expander` - `least-waste`, see [Here](../../FAQ.md#what-are-expanders) ) Expander ` price` - not supported.

You must implement the initialization of the node and add it to the cluster. To initialize a node, you can use the Cloud Init, which can be added to the Cluster Autoscaler configuration.

Cluster Autosclaer only removes nodes from the Hetzner Cloud, but they also need to be removed from your k8s cluster. CA can't do it. CA expects nodes to contain `spec.providerID` by which it defines the ID of node. Without this, the CA will not work correctly. Therefore you must install [cloud-controller-manager](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/) to your k8s cluster. For Hetzner Cloud is the [hcloud-cloud-controller-manager](https://github.com/hetznercloud/hcloud-cloud-controller-manager). In addition to removing nodes from the k8s cluster, the cloud-controller-manager initializes them when they are added to the cluster, creates `spec.providerID`  and useful labels `instance-type`, `region`, `zone`. If you have hybrid environment and you use cloud servers and bare-metal(dedicated) servers in the same cluster, you can use [hetzner-cloud-controller-manager](https://github.com/identw/hetzner-cloud-controller-manager), it can work both with the cloud nodes of Hetzner Cloud and with bare-metal from Hetzner Robot.

**Note**: when using hetzner-cloud-controller-manager instead of hcloud-cloud-controller-manager you need to change **provider_prefix** from `hcloud: //` to `hetzner: //`.

## Deployment
You need to create a token to access the Api Hetzner Cloud. To do this, follow the [instruction](https://docs.hetzner.cloud/#overview-getting-started).

Create secret `cluster-autoscaler-cloud-config` in `kube-system` namespace, that will contain a token and configuration for CA. [Example secret](./examples/secret.yaml):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cluster-autoscaler-cloud-config
  namespace: kube-system
stringData:
  cloud-config: |
    {
      "token":"hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb",
      "endpoint":"https://api.hetzner.cloud/v1",
      "ssh_keys": [111, 112],
      "cloud_init": "IyEvdXNyL2Jpbi9lbnYgYmFzaAojIFRlc3RlZCBpbiB1YnVudHUgMTguMDQKZXhwb3J0IERFQklBTl9GUk9OVEVORD1ub25pbnRlcmFjdGl2ZQpjb2RlbmFtZT1gbHNiX3JlbGVhc2UgLWNzYAoKIyMjIHN5c3RlbQpjYXQgPiAvZXRjL3N5c2N0bC5kLzk5LWs4cy5jb25mIDw8RU9GCiMgUkVRVUlSRUQgCiMgaHR0cHM6Ly9rdWJlcm5ldGVzLmlvL2RvY3MvY29uY2VwdHMvZXh0ZW5kLWt1YmVybmV0ZXMvY29tcHV0ZS1zdG9yYWdlLW5ldC9uZXR3b3JrLXBsdWdpbnMvI25ldHdvcmstcGx1Z2luLXJlcXVpcmVtZW50cwojIGh0dHBzOi8va3ViZXJuZXRlcy5pby9kb2NzL3NldHVwL3Byb2R1Y3Rpb24tZW52aXJvbm1lbnQvY29udGFpbmVyLXJ1bnRpbWVzLwpuZXQuaXB2NC5pcF9mb3J3YXJkPTEKbmV0LmlwdjYuY29uZi5hbGwuZm9yd2FyZGluZz0xCm5ldC5icmlkZ2UuYnJpZGdlLW5mLWNhbGwtaXB0YWJsZXM9MQoKIyBGSVhFRCBwcm9ibGVtcwpmcy5pbm90aWZ5Lm1heF91c2VyX3dhdGNoZXM9NTI0Mjg4ICMgZml4OiBmYWlsZWQgdG8gd2F0Y2ggZmlsZSAiL3Zhci9saWIvZG9ja2VyL2NvbnRhaW5lcnMvIjogbm8gc3BhY2UgbGVmdCBvbiBkZXZpY2UKRU9GCnN5c2N0bCAtLXN5c3RlbQoKY2F0ID4gL2V0Yy9tb2R1bGVzLWxvYWQuZC9icl9uZXRmaWx0ZXIuY29uZiA8PEVPRgpicl9uZXRmaWx0ZXIKRU9GCm1vZHByb2JlIGJyX25ldGZpbHRlcgoKIyMgZG9ja2VyIGFuZCBrOHMgZGVwZW5kcwphcHQtZ2V0IHVwZGF0ZSAmJiBhcHQtZ2V0IGluc3RhbGwgLXkgYXB0LXRyYW5zcG9ydC1odHRwcyBjYS1jZXJ0aWZpY2F0ZXMgY3VybCBzb2Z0d2FyZS1wcm9wZXJ0aWVzLWNvbW1vbgpjdXJsIC1mc1NMIGh0dHBzOi8vZG93bmxvYWQuZG9ja2VyLmNvbS9saW51eC91YnVudHUvZ3BnIHwgYXB0LWtleSBhZGQgLQpjdXJsIC1zIGh0dHBzOi8vcGFja2FnZXMuY2xvdWQuZ29vZ2xlLmNvbS9hcHQvZG9jL2FwdC1rZXkuZ3BnIHwgYXB0LWtleSBhZGQgLQpjYXQgPDxFT0YgPi9ldGMvYXB0L3NvdXJjZXMubGlzdC5kL2t1YmVybmV0ZXMubGlzdApkZWIgaHR0cHM6Ly9hcHQua3ViZXJuZXRlcy5pby8ga3ViZXJuZXRlcy14ZW5pYWwgbWFpbgpFT0YKY2F0IDw8RU9GID4vZXRjL2FwdC9zb3VyY2VzLmxpc3QuZC9kb2NrZXIubGlzdApkZWIgW2FyY2g9YW1kNjRdIGh0dHBzOi8vZG93bmxvYWQuZG9ja2VyLmNvbS9saW51eC91YnVudHUgJHtjb2RlbmFtZX0gc3RhYmxlCkVPRgphcHQtZ2V0IHVwZGF0ZQphcHQtZ2V0IGluc3RhbGwgLXkgIGRvY2tlci1jZT01OjE5LjAzLjR+My0wfnVidW50dS0ke2NvZGVuYW1lfSBrdWJlbGV0PTEuMTcuMi0wMCBrdWJlYWRtPTEuMTcuMi0wMCBrdWJlY3RsPTEuMTcuMi0wMAphcHQtbWFyayBob2xkIGt1YmVsZXQga3ViZWFkbSBrdWJlY3RsIGRvY2tlci1jZQoKY2F0ID4gL2V0Yy9kb2NrZXIvZGFlbW9uLmpzb24gPDxFT0YKewogICJleGVjLW9wdHMiOiBbIm5hdGl2ZS5jZ3JvdXBkcml2ZXI9c3lzdGVtZCJdLAogICJsb2ctZHJpdmVyIjogImpzb24tZmlsZSIsCiAgImxvZy1vcHRzIjogewogICAgIm1heC1zaXplIjogIjEwMG0iCiAgfSwKICAic3RvcmFnZS1kcml2ZXIiOiAib3ZlcmxheTIiLAogICJpcC1mb3J3YXJkIjogZmFsc2UsCiAgImlwLW1hc3EiOiBmYWxzZSwKICAiaXB0YWJsZXMiOiBmYWxzZSwKICAiYnJpZGdlIjogIm5vbmUiCn0KRU9GCm1rZGlyIC1wIC9ldGMvc3lzdGVtZC9zeXN0ZW0vZG9ja2VyLnNlcnZpY2UuZCB8fCB0cnVlCm1rZGlyIC1wIC9ldGMvc3lzdGVtZC9zeXN0ZW0va3ViZWxldC5zZXJ2aWNlLmQgfHwgdHJ1ZQpjYXQgPiAvZXRjL3N5c3RlbWQvc3lzdGVtL2t1YmVsZXQuc2VydmljZS5kLzIwLWV4dGVybmFsLWNsb3VkLmNvbmYgPDxFT0YKW1NlcnZpY2VdCkVudmlyb25tZW50PSJLVUJFTEVUX0VYVFJBX0FSR1M9LS1jbG91ZC1wcm92aWRlcj1leHRlcm5hbCIKRU9GCnN5c3RlbWN0bCBkYWVtb24tcmVsb2FkCnN5c3RlbWN0bCByZXN0YXJ0IGRvY2tlcgpzeXN0ZW1jdGwgcmVzdGFydCBrdWJlbGV0CgojIGNsZWFuIGRvY2tlciBpcHRhYmxlcwppcHRhYmxlcyAtdCBuYXQgLUYKaXB0YWJsZXMgLUYKCktVQkVfQVBJX0VORFBPSU5UPWlwX2FkZHJlc3M6cG9ydApLVUJFX1RPS0VOPXRva2VuCktVQkVfVE9LRU5fQ0FfQ0VSVD1zaGEyNTY6aGFzaAprdWJlYWRtIGpvaW4gJHtLVUJFX0FQSV9FTkRQT0lOVH0gLS10b2tlbiAke0tVQkVfVE9LRU59IC0tZGlzY292ZXJ5LXRva2VuLWNhLWNlcnQtaGFzaCAke0tVQkVfVE9LRU5fQ0FfQ0VSVH0K",
      "instance_type": "cx51",
      "location": "hel1",
      "image": {
        "id": 168855,
        "name": "ubuntu-18.04",
        "type": "system"
      },
      "pools": {
        "k8s-autoscaler1": {
          "node_name_prefix":"kube-worker102-1"
        },
        "k8s-autoscaler2": {
          "node_name_prefix":"kube-worker102-2",
          "ssh_keys": [113, 114],
          "instance_type": "cx41",
          "location": "nbg1"
        }
      }
    }
```
### ssh_keys
You can find out id keys in `ssh_keys` using [hcloud-cli](https://github.com/hetznercloud/cli):
```bash
$ hcloud ssh-key list
ID        NAME        FINGERPRINT
111       john_doe    5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1a
112       user1       5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1b
113       user2       5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1c
114       user3       5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1d
```
Or make a request in API:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/ssh_keys
```

### image
List available image ids:
```bash
$ hcloud image list
ID        TYPE     NAME           DESCRIPTION    IMAGE SIZE   DISK SIZE   CREATED
1         system   ubuntu-16.04   Ubuntu 16.04   -            5 GB        2 years ago
2         system   debian-9       Debian 9       -            5 GB        2 years ago
3         system   centos-7       CentOS 7       -            5 GB        2 years ago
168855    system   ubuntu-18.04   Ubuntu 18.04   -            5 GB        2 years ago
5924233   system   debian-10      Debian 10      -            5 GB        7 months ago
8356453   system   centos-8       CentOS 8       -            5 GB        4 months ago
9032843   system   fedora-31      Fedora 31      -            5 GB        3 months ago
```
API:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/images
```
You can use snapshots instead of images. Example:
```json
      "image": {
        "id": 1221313,
        "name": "my-snapshot-name",
        "type": "snapshot"
      },
```

### location
Avaliable locations:
```bash
$ hcloud datacenter list
ID   NAME        DESCRIPTION          LOCATION
2    nbg1-dc3    Nuremberg 1 DC 3     nbg1
3    hel1-dc2    Helsinki 1 DC 2      hel1
4    fsn1-dc14   Falkenstein 1 DC14   fsn1
```
API:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/datacenters
```

### instance_type
Available server types (`instance_type`):
```bash
$ hcloud server-type list
ID   NAME        CORES   MEMORY     DISK     STORAGE TYPE
1    cx11        1       2.0 GB     20 GB    local
2    cx11-ceph   1       2.0 GB     20 GB    network
3    cx21        2       4.0 GB     40 GB    local
4    cx21-ceph   2       4.0 GB     40 GB    network
5    cx31        2       8.0 GB     80 GB    local
6    cx31-ceph   2       8.0 GB     80 GB    network
7    cx41        4       16.0 GB    160 GB   local
8    cx41-ceph   4       16.0 GB    160 GB   network
9    cx51        8       32.0 GB    240 GB   local
10   cx51-ceph   8       32.0 GB    240 GB   network
11   ccx11       2       8.0 GB     80 GB    local
12   ccx21       4       16.0 GB    160 GB   local
13   ccx31       8       32.0 GB    240 GB   local
14   ccx41       16      64.0 GB    360 GB   local
15   ccx51       32      128.0 GB   600 GB   local
```
API:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/server_types
```

### cloud_init
`cloud_init` must be encoding to `Base64`. This parameter is decoded and transmitted as user_data during [server creation](https://docs.hetzner.cloud/#servers-create-a-server). `cloud_init` is the [Cloud Init](https://cloudinit.readthedocs.io/en/latest/).

The following script is used in the example:
```bash
#!/usr/bin/env bash
# Tested in ubuntu 18.04
export DEBIAN_FRONTEND=noninteractive
codename=`lsb_release -cs`

### system
cat > /etc/sysctl.d/99-k8s.conf <<EOF
# REQUIRED 
# https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/#network-plugin-requirements
# https://kubernetes.io/docs/setup/production-environment/container-runtimes/
net.ipv4.ip_forward=1
net.ipv6.conf.all.forwarding=1
net.bridge.bridge-nf-call-iptables=1

# FIXED problems
fs.inotify.max_user_watches=524288 # fix: failed to watch file "/var/lib/docker/containers/": no space left on device
EOF
sysctl --system

cat > /etc/modules-load.d/br_netfilter.conf <<EOF
br_netfilter
EOF
modprobe br_netfilter

## docker and k8s depends
apt-get update && apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
cat <<EOF >/etc/apt/sources.list.d/docker.list
deb [arch=amd64] https://download.docker.com/linux/ubuntu ${codename} stable
EOF
apt-get update
apt-get install -y  docker-ce=5:19.03.4~3-0~ubuntu-${codename} kubelet=1.17.2-00 kubeadm=1.17.2-00 kubectl=1.17.2-00
apt-mark hold kubelet kubeadm kubectl docker-ce

cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2",
  "ip-forward": false,
  "ip-masq": false,
  "iptables": false,
  "bridge": "none"
}
EOF
mkdir -p /etc/systemd/system/docker.service.d || true
mkdir -p /etc/systemd/system/kubelet.service.d || true
cat > /etc/systemd/system/kubelet.service.d/20-external-cloud.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--cloud-provider=external"
EOF
systemctl daemon-reload
systemctl restart docker
systemctl restart kubelet

# clean docker iptables
iptables -t nat -F
iptables -F

KUBE_API_ENDPOINT=ip_address:port
KUBE_TOKEN=token
KUBE_TOKEN_CA_CERT=sha256:hash
kubeadm join ${KUBE_API_ENDPOINT} --token ${KUBE_TOKEN} --discovery-token-ca-cert-hash ${KUBE_TOKEN_CA_CERT}

```
The script installs everything that the node needs (docker, cubelet, kubeadm) for initialization in the cluster. And then joins the cluster using `kubeadm join`. You must specify the correct cluster join params (KUBE_API_ENDPOINTKUBE_TOKEN,KUBE_TOKEN_CA_CERT). Pay attention to the additional parameter `--cloud-provider=external` for kubelet, the cloud-controller-manager needs it to initialize the node in your k8s cluster.

### pools

`pools` describes the configuration for each node pool. The key is the pool name, and it must match the pool name in the `--node` arguments for cluster-autoscaler. For each pool, you can override the parameters `ssh_keys`,` cloud_init`, `instance_type`,` location`, `image`. The `node_name_prefix` parameter must be unique for each pool, this is the prefix of the names of the created nodes (`{node_name_prefix} - {random_number}`)


Deploy cluster autosclaer:
```bash
kubectl apply -f examples/deploy.yaml
```

Correct the [deploy.yaml](./examples/deploy.yaml) file to your needs. Note the options `--nodes` - [deploy.yaml#nodes](./examples/deploy.yaml#L165). Its are set in the format `MIN_NODES:MAX_NODES:POOL_NAME`, where` MIN_NODES` is the minimum number of nodes in the pool, CA will delete unnecessary nodes until it reaches this minimum, you can set it to 0. `MAX_NODES` is the maximum number of nodes in the pool. `POOL_NAME` - the name of the pool, **must** be present in the` pools` config.

## Configuration
The configuration file for cluster-autoscaler is passed to the [--cloud-config](./examples/deploy.yaml#L163) parameter, the file must be in the `JSON` format. Configuration example: [secret.yaml](./examples/secret.yaml#L7).

Available parameters:

| parameter              | type           | description                                                                                                                        | default                                             |
|------------------------|----------------|------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| `token`(required)      | string         | Token for access to Api                                                                                                            | not                                                 |
| `endpoint`(required)   | string         | Api endpoint                                                                                                                       | not                                                 |
| `ssh_keys`(required)   | array          | Array of numeric id public ssh keys, see in the [ssh_keys](#ssh_keys) section                                                      | not                                                 |
| `cloud_init`(required) | string(base64) | Base64 encoded Cloud-Init                                                                                                          | not                                                 |
| `instance_type`        | string         | Server type, see in the [instance_type](#instance_type) section                                                                    | cx41                                                |
| `location`             | string         | Location. Available `nbg1`, `fsn1`, `hel1`                                                                                         | nbg1                                                |
| `image`                | object         | The system image for the server, see in the [image](#image) section                                                                | {"id":168855,"name":"ubuntu-18.04","type":"system"} |
| `pools`(required)      | object         | configuration of pools, see in the [pools](#pools) section                                                                         | not                                                 |
| `provider_prefix`      | string         | Provider prefix from node `spec.providerID`. Must be `hcloud://` or `hetzner://`. See in the [How it works](#how-it-works) section | hcloud://                                           |

An example of a configuration file is in the [Deployment](#deployment) section.