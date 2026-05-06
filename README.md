<p align="center">
  <img src="./attachments/wave.svg" width="100%" alt="Header Waves">
</p>

# Table of contents

- [About this app](#about-this-app)
  - [Quickstart](#quickstart)
  - [Project Architecture](#project-architecture)
  - [Ansible usage](#ansible-usage)
  - [Kubernetes deploy](#kubernetes-deploy)
- [About me](#about-me)
- [Technologies I placed my hands on](#installation)
  - [IDEs](#ides)
  - [Unix-like Systems](#unix-like-systems)
  - [Linux Distros](#linux-distros)
  - [Programming Languages](#programming-languages)
  - [Networking](#networking)
  - [Other tools](#other-tools)
- [Pet-projects](#pet-projects)
  - [golang-url-shortener](#golang-url-shortener)
  - [server-checker](#server-checker)
  - [gymctl](#gymctl-coming-soon)
- [Area of interests](#area-of-interests)
- [Contacts](#contacts)

## About this app

Implemented RESTful API web service with two endpoints (`/info` and `/info/weather`) for retrieving weather forecast via VisualCrossing API.

### Quickstart

1. **Set up environment:** Create a `.env` file in the root directory:

```env
API_KEY=your_visual_crossing_api_key # default: EMPTY
PORT=your_port                       # default: 8000
AUTHOR=egor.volkov                   # me :O
VERSION=0.3.0                        # default: 1.0.0
SERVICE=weather                      # service name
```

2. **Run the service**:

```bash
go run main.go
```

### Project architechture

### Entities

- **weather.go**: `TemperatureStats` (average, median, min, max), `DayData` (avg temp, min temp, max temp)
- **service.go**: `ServiceInfo` (version, serviceName, author)

### Use cases

- **weather_usecase.go**
  - `GetStats`: app's business logic (average, median, min and max temperature estimation)

### Adapters

- **http_handler.go**
  - `GetWeather`: handler to form URL from user request
- **weather_api.go**
  - `GetRawWeatherData`: function to send url request to VisualCrossing API and retrieve data

### Infrastructure

- **main.go**
  - `main`: reading variables, dependency injection, starting Gin router
  - `getEnv`: additional function for reading vars from .env

## Ansible usage

### Step 1: Install Dependencies

Download required collections from GitHub/Galaxy:

```bash
ansible-galaxy install -r requirements.yml
```

### Step 2: Configure Inventory

Edit inventory.ini with your server IPs:

```ini
[k8s_nodes]
worker1 ansible_host=10.184.0.24 ansible_user=worker ansible_ssh_private_key_file=~/.ssh/yadro_worker1_ed25519
worker2 ansible_host=10.184.0.23 ansible_user=worker ansible_ssh_private_key_file=~/.ssh/yadro_worker2_ed25519
master ansible_host=10.184.0.25 ansible_user=master ansible_ssh_private_key_file=~/.ssh/yadro_master_ed25519


[k8s_nodes:vars]
ansible_python_interpreter=/usr/bin/python3
```

### Step 3: Run Playbook

Execute the automation:

```bash
ansible-playbook -i inventory.ini playbook.yml -K

-K will prompt for the sudo password (of remote nodes)
```

### 4. Roles Overview

k8s_node_prep

- Disables SWAP
- Loads overlay and br_netfilter modules
- Sets net.ipv4.ip_forward = 1
- Terminates conflicting apt background processes

cri_o

- Default version: 1.33.
- Adds OpenSUSE OBS repositories.
- Installs and starts the crio service.

kube_tools

- Default version: 1.33.
- Uses Yandex Mirrors for stability.
- Installs kubeadm, kubelet, kubectl.

### 5. Testing with Molecule

Molecule is used to test the cri_o role in a Docker container.

#### Setup Environment

```bash
python3 -m venv venv
source venv/bin/activate
pip install molecule molecule-plugins[docker] ansible-core
```

#### Run Tests

Go to the role directory and start the test cycle:

```
cd ansible/roles/cri_o
molecule test
```

What happens: Molecule creates a jrei/systemd-ubuntu container, installs sudo via prepare.yml, and applies the role. The test verifies syntax, installation, and idempotency.

## Kubernetes deploy

The application is deployed in a high-availability cluster using **Traefik** as the Ingress Controller.

### 1. High Availability

- **Replicas**: The `weather-deployment` runs **2 replicas** distributed across worker nodes. This ensures that the service remains available even if a single pod or node fails.
- **Namespaces**: Resources are logically isolated:
  - `traefik`: Infrastructure
  - `weather-app`: Application, Services, and Secrets.

### 2. Networking & External Access

Since the cluster is running on private IPs, external access is provided via **NodePort** services:

- **HTTP**: `http://<NODE_IP>:31300`
- **HTTPS**: `https://<NODE_IP>:31301`
- **Traefik Dashboard**: `http://<NODE_IP>:31300/dashboard/`

### 3. Ingress & Security (HTTPS)

The project includes a full TLS termination setup:

- **TLS Secret**: A self-signed certificate is stored in a Kubernetes Secret to encrypt traffic.
- **HTTPS Redirection**: A Traefik `Middleware` using **RedirectRegex** is implemented. It automatically redirects users from HTTP to HTTPS, specifically handling the transition between NodePorts (`31300` -> `31301`).
- **TLSStore**: A custom `TLSStore` is configured to serve our certificate by default, preventing Traefik from using its internal "Default Cert" when accessing via IP addresses.

### 4. Applied Approach (Why this way?)

| Feature                | Choice               | Why?                                                                                                        |
| :--------------------- | :------------------- | :---------------------------------------------------------------------------------------------------------- |
| **Ingress Controller** | **Traefik**          | Modern, lightweight, and supports CRDs (IngressRoute, Middleware) for fine-grained traffic control.         |
| **Service Type**       | **NodePort**         | Ideal for clusters without an integrated Cloud LoadBalancer. It maps ports across all cluster nodes.        |
| **Redirection**        | **Regex Middleware** | Standard HTTPS redirects often fail with custom ports. Regex allows precise rewriting of the URL and Port.  |
| **Security**           | **TLS Termination**  | Security is handled at the Ingress level (edge), so the Go application doesn't need to manage certificates. |

### Deployment Steps

1. Create **namespaces**:

```bash
kubectl create namespace traefik
kubectl create namespace weather-app
```

2. **Install Traefik** (via Helm):

```bash
helm install traefik traefik/traefik \
--namespace traefik \
--set ports.web.nodePort=31300 \
--set ports.websecure.nodePort=31301 \
--set service.type=NodePort
```

3. Generate **self-signed certificate**:

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
-keyout tls.key -out tls.crt \
-subj "/CN=weather.local"
```

4. Add **k8s secrets** (VisualCrossing API-KEY and certificate):

```bash
kubectl create secret tls weather-tls-secret --key tls.key --cert tls.crt -n weather-app
kubectl create secret generic api-key-secret \
--from-literal=api-key=<YOUR-API-KEY> \
-n weather-app
```

5. Apply all **manifests**:

```bash
cd k8s/manifests
kubectl apply -f . -R
```

## About me

Greetings, my name is Volkov Egor.

I'm currently studying CyberSec (10.03.01) at the Moscow Power Engineering Institute and at the same time working at InfoWatch in the field of network security (NGFW support engineer).

I`m currently working on a tui-based tool (Golang/bubbletea) for validating physical appliances (acknowledging hardware components' model, current load analysis and comparing with specification standard, generating html/yaml/txt reports).

My hobby is doing sports such as lifting weights (going to the gym), running and very rarely playing volleyball.

## Technologies I placed my hands on

### IDEs

<p alignh="left">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/vscode/vscode-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/vscodium/vscodium-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/pycharm/pycharm-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/goland/goland-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/vim/vim-original.svg" height="45"/>
</p>

### Unix-like Systems

<p alignh="left">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/linux/linux-original.svg" height="45"/>
<img src="https://www.vectorlogo.zone/logos/freebsd/freebsd-icon.svg" height="45"/>
</p>

### Linux Distros

<p alignh="left">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/centos/centos-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/fedora/fedora-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/debian/debian-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ubuntu/ubuntu-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/linuxmint/linuxmint-original.svg" height="45"/>
</p>

### Programming Languages

<p alignh="left">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/bash/bash-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/python/python-original.svg" height="45"/>
</p>

### Networking

<p alignh="left">
<img src="https://eltex-co.com/images/rules/01.svg" height="50" width="50"/>
<img src="https://www.vectorlogo.zone/logos/cisco/cisco-ar21.svg" height="25"/>
<img src="https://assets.streamlinehq.com/image/private/w_300,h_300,ar_1/f_auto/v1/icons/logos/mikrotik-9fxkorisqhnq5ppw9ory2o.png/mikrotik-w4y9rth430h5bcfzp9in8i.png?_a=DATAiZAAZAA0" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/svg/opnsense.svg" height="45"/>
<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/fd/VyOS_logo.svg/250px-VyOS_logo.svg.png?_=20251130143906" height="45"/>
<img src="https://infowatch.com/themes/infowatch/img/new-style/logo-arma.svg" height="45"/>
</p>

### Other tools

<p alignh="left">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg" height="45"/>
<img src="https://git-scm.com/images/logos/downloads/Git-Icon-1788C.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg" height="45"/>
<img src="https://upload.wikimedia.org/wikipedia/commons/c/c6/Wireshark_icon_new.png" height="45"/>
</p>

## Pet-projects

### golang-url-shortener

Just a project-based learning, nothing serious:

- Gin (handlers for creating (hashing + storing with original values) and for providing created short url)
- Redis (to store map short url<->original url)

That's it

### server-checker

server-checker is a terminal-based hardware validation tool for Linux appliances and security platforms.

It collects system-level hardware data, validates it against predefined or custom YAML templates, and generates structured audit-ready reports.

![example.gif](https://media3.giphy.com/media/v1.Y2lkPTc5MGI3NjExZ3MwdDM0anMyanFjcDVycnMweHZ6bnRqancwOGgzeWc4bnE2Y2djNSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/ALBDzTOeA11mgHbeqE/giphy.gif)

Report output structure:

```bash
HealthCheck-YYYY-MM-DD_HH-MM-SS/
├── report.txt
├── report.yaml
├── report.html
└── logs/
    ├── dmesg.txt
    ├── lspci-k.txt
    ├── lsusb.txt
    ├── lsblk.txt
    ├── product_name.txt
    ├── product_serial.txt
    ├── sys_vendor.txt
    └── smart_<device>.txt
```

### gymctl (coming soon)

Golang client/server app to let users track their workouts and progress.

## Area of interests

- Linux infrastructure
- Infrastructure automation
- Observability and diagnostics
- Backend tooling in Go

## Contacts

Email: wolkega@mail.ru  
Telegram: @w3athr  
GitHub: https://github.com/w3athr
