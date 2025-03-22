# minecraftd

A Minecraft server daemon.

## Setup

1. Download the [Minecraft server](https://www.minecraft.net/download/server) to `/var/lib/minecraft/server.jar`.
2. Download the `minecraftd-${version}.tar.gz` from the [latest release](https://github.com/jessestricker/minecraftd/releases/latest).
3. Extract the archive:
   ```shell
   sudo tar -xzvf minecraftd-${version}.tar.gz -C /
   ```
4. Enable and start the service:
   ```shell
   sudo systemctl daemon-reload
   sudo systemctl enable --now minecraft.service
   ```
