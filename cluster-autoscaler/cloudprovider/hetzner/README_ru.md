# Cluster Autoscaler для Hetzner

Автоматическое масштабирование кластера для [Hetzner Cloud](https://www.hetzner.com/cloud). Может масштабировать несколько пулов узлов. Пулы задаются по меткам, которыми можно отмечать узлы в Hetzner Cloud. Запускается в кластере как `Deployment`.

## Как работает
Cluster Autoscaler создает узлы при необходимости через Api Hetzner Cloud, и также удаляет их, когда узлы больше не нужны. Более подробную информацию вы можете узнать в [FAQ](../../FAQ.md). Поддерживается работа только внутри одного [проекта](https://wiki.hetzner.de/index.php/CloudServer/en#What_are_projects.2C_and_how_can_I_use_them.3F). Каждому узлу пула назначается своя [метка](https://docs.hetzner.cloud/#overview-labels), которая совпадает с именем пула. По меткам определяется принадлежность узла к пулу.

Поддерживается метод `TemplateNodeInfo`, по которому выбирается лучший пул для масштабирования(при выборе соотвесвтующего `expander` - `least-waste`, более подробно смотрите [Здесь](../../FAQ.md#what-are-expanders)). Expander `price` - не поддерживается.

Реализовать инициализацию ноды и добавление её в кластер, вам потребуется самому. Для этого вы можете использовать **Cloud-Init**, который можно добавить в конфигурацию Cluster Autoscaler. Ниже будет приведен пример такого cloud-init.

При удалении нод Cluster Autosclaer'ом из Hetzner Cloud, нужно их удалять и из вашего кластера, CA этого не делает. Помимо этого CA ожидает, что ноды содержат `spec.providerID` по которым он определяет `ID` узла. Без этого CA не будет корректно работать. Поэтому вы должны установить в кластер [cloud-controller-manager](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/). Для Hetzner Cloud есть готовый [hcloud-cloud-controller-manager](https://github.com/hetznercloud/hcloud-cloud-controller-manager), который требуется установить в кластер. Кроме удаления нод из кластера, cloud-controller-manager их инициализирует при добавлении в кластер, создавая `spec.providerID` и полезные метки `instance-type`, `region`, `zone`. Если у вас гибридная среда, и вы в одном кластере используете облачные и bare-metal сервера, вы можете воспользоваться [hetzner-cloud-controller-manager](https://github.com/identw/hetzner-cloud-controller-manager), он умеет работать как c облачными узлами Hetzner Cloud так и с bare-metal из Hetzner Robot.

**Обратите внимание:** при использовании hetzner-cloud-controller-manager вместо hcloud-cloud-controller-manager вам нужно сменить  **provider_prefix** с `hcloud://` на `hetzner://`.

## Деплой
Создайте токен для доступа к Api Hetzner Cloud, [ЗДЕСЬ](https://docs.hetzner.cloud/#overview-getting-started) можно найти более подробную информацию.

Создайте секрет `cluster-autoscaler-cloud-config` в namespace `kube-system`, который будет содержать токен и конфигурацию для CA. Например ([Example secret](./examples/secret.yaml)):
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
Id ключей в `ssh_keys` вы можете узнать с помощью [hcloud-cli](https://github.com/hetznercloud/cli):
```bash
$ hcloud ssh-key list
ID        NAME        FINGERPRINT
111       john_doe    5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1a
112       user1       5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1b
113       user2       5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1c
114       user3       5d:10:31:45:74:3e:af:3a:a5:ee:b6:81:10:24:87:1d
```
Или сделать запрос в Api:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/ssh_keys
```

### image
Аналогичным образом вы можете уточнить `id`, `name`, `type` для параметра `image`.
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
Api:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/images
```
Есть возможность использовать снапшоты вместо образов, для этого достаточно правильно указать id и имя, и сменить `type` с `system` на `snapshot`.
```json
      "image": {
        "id": 1221313,
        "name": "my-snapshot-name",
        "type": "snapshot"
      },
```

### location
Доступные локации:
```bash
$ hcloud datacenter list
ID   NAME        DESCRIPTION          LOCATION
2    nbg1-dc3    Nuremberg 1 DC 3     nbg1
3    hel1-dc2    Helsinki 1 DC 2      hel1
4    fsn1-dc14   Falkenstein 1 DC14   fsn1
```
Api:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/datacenters
```

### instance_type
Доступные типы серверов (`instance_type`):
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
Api:
```bash
curl -H "Authorization: Bearer hIsFGNoKJG7qMooeHl6JpyD8UQtI0QWJIS9jp2XCHoYjPf8ofPR6h5v7WwWckvUb" https://api.hetzner.cloud/v1/server_types
```

### cloud_init

Параметр `cloud_init` должен быть закодирован в `Base64`. Данный параметр декодируется и передается как user_data при [создании сервера](https://docs.hetzner.cloud/#servers-create-a-server), другими словами это обычный [Cloud Init](https://cloudinit.readthedocs.io/en/latest/). В данном примере используется скрипт:
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

Как говорилось выше, вы должны сами позаботиться о добавлении ноды в кластер. В данном примере в конце скрипта с помощью `kubeadm join` как раз нода добавляется в кластер. Также обратите внимание на дополнительный параметр `--cloud-provider=external` для kubelet, он необходим для инициализации ноды через cloud-controller-manager.

### pools

`pools` описывает конфигурацию для каждого пула нод. Ключ является именем пула, и он должен совпадать с именем пула в аргментах `--node` для cluster-autoscaler. Для каждого пула вы можете переопределить параметры `ssh_keys`, `cloud_init`, `instance_type`, `location`, `image`.  Параметр `node_name_prefix` должен быть уникальным для каждого пула, это префикс имен создаваемых нод (`{node_name_prefix}-{random_number}`)

Деплой Cluster Autosclaer'а:
```bash
kubectl apply -f examples/deploy.yaml
```
Файл [deploy.yaml](./examples/deploy.yaml) потребуется исправить под свои нужды. Обратите внимание на параметры `--nodes` - [deploy.yaml#nodes](./examples/deploy.yaml#L165). Они задаются в формате `MIN_NODES:MAX_NODES:POOL_NAME`, где `MIN_NODES` - минимальное количество нод в пуле, CA будет удалять не нужные ноды пока не достигнет этого минимума, можно поставить в 0. `MAX_NODES` - максимальное количество нод в пуле. `POOL_NAME` - имя пула, **должно** присутсвтовать в `pools` конфига.

## Конфигурация
Конфигурационный файл для cluster-autoscaler передается в параметр [--cloud-config](./examples/deploy.yaml#L163), файл должен быть в формате `JSON`. Пример конфигурационного: [secret.yaml](./examples/secret.yaml#L7).

Возможные параметры:

| параметр                 | тип            | описание                                                                                                                            | умолчание                                           |
|--------------------------|----------------|-------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| `token`(обязателен)      | string         | Токен для доступа к API                                                                                                             | нет                                                 |
| `endpoint`(обязателен)   | string         | Url API                                                                                                                             | нет                                                 |
| `ssh_keys`(обязателен)   | array          | Массив числовых id публичных ssh ключей, смотрите в разделе [ssh_keys](#ssh_keys)                                                   | нет                                                 |
| `cloud_init`(обязателен) | string(base64) | Cloud-Init закодированный в base64                                                                                                  | нет                                                 |
| `instance_type`          | string         | Тип сервера, сомтрите в разделе [instance_type](#instance_type)                                                                     | cx41                                                |
| `location`               | string         | Локация, доступны `nbg1`, `fsn1`, `hel1`                                                                                            | nbg1                                                |
| `image`                  | object         | Образ системы для сервера, смотрите в разделе [image](#image)                                                                       | {"id":168855,"name":"ubuntu-18.04","type":"system"} |
| `pools`(обязателен)      | object         | конфигураця пулов, смотрите в разделе [pools](#pools)                                                                               | нет                                                 |
| `provider_prefix`        | string         | Provder prefix из `spec.providerId` ноды. Может быть `hcloud://` либо `hetzner://`. Смотри в разделе [Как работает](#как-работает)  | hcloud://                                           |

Пример конфигурационного файла есть в разделе [Деплой](#деплой)