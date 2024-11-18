---
Title: Docker on Raspberry Pi
Description: Test how a really long description looks like in the metadata of a markdown file on a Raspberry Pi
Created: 2024-10-29
Modified: 2024-10-29
StyleSheets:
  - /static/styles/utils/code.css
Scripts:
  - /static/scripts/utils/copy-code.js
---

# Docker on Raspberry Pi

## Install Docker

```bash
curl -sSL https://get.docker.com | sh
```

## Add user to Docker group hello-world

```bash
sudo usermod -aG docker pi
sudo usermod -aG docker pi
```

## Reboot

```bash
sudo reboot
sudo reboot
sudo reboot
```

## Test Docker

```bash
docker run hello-world
docker run hello-world
docker run hello-world
docker run hello-world
```

## References

- [Docker on Raspberry Pi](https://www.raspberrypi.org/blog/docker-comes-to-raspberry-pi/)
