<p align="center">
  <img src="./attachments/wave.svg" width="100%" alt="Header Waves">
</p>

# Table of contents

- [About me](#about-me)
- [Technologies I placed my hands on](#installation)
  - [IDEs](#ides)
  - [Unix-like Systems](#unix-like-systems)
  - [Linux Distros](#linux-distros)
  - [Programming Languages](#programming-languages)
  - [Networking](#networking)
  - [Other tools](#other-tools)
- [Projects](#projects)
  - [golang-url-shortener](#golang-url-shortener)
  - [server-checker](#server-checker)
  - [gymctl](#gymctl-coming-soon)
- [Area of interests](#area-of-interests)
- [Contacts](#contacts)

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
<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/c/c2/Font_Awesome_5_brands_freebsd.svg/640px-Font_Awesome_5_brands_freebsd.svg.png" height="45"/>
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
<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/08/Cisco_logo_blue_2016.svg/640px-Cisco_logo_blue_2016.svg.png" height="25"/>
<img src="https://assets.streamlinehq.com/image/private/w_300,h_300,ar_1/f_auto/v1/icons/logos/mikrotik-9fxkorisqhnq5ppw9ory2o.png/mikrotik-w4y9rth430h5bcfzp9in8i.png?_a=DATAiZAAZAA0" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/svg/opnsense.svg" height="45"/>
<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/fd/VyOS_logo.svg/250px-VyOS_logo.svg.png?_=20251130143906" height="45"/>
<img src="https://infowatch.com/themes/infowatch/img/new-style/logo-arma.svg" height="45"/>
</p>

### Other tools

<p alignh="left">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg" height="45"/>
<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Git_icon.svg/640px-Git_icon.svg.png" height="45"/>
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg" height="45"/>
<img src="https://upload.wikimedia.org/wikipedia/commons/c/c6/Wireshark_icon_new.png" height="45"/>
</p>

## Projects

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
