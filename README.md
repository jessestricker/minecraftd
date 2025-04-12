# minecraftd

A Minecraft server daemon.

## Setup

```bash
# URL of Minecraft server JAR file, see https://www.minecraft.net/download/server
server_url='https://piston-data.mojang.com/v1/objects/4707d00eb834b446575d89a61a11b5d548d8c001/server.jar'

# URL of minecraftd, see https://github.com/jessestricker/minecraftd/releases/latest
minecraftd_url='https://github.com/jessestricker/minecraftd/releases/download/v1.5.0/minecraftd-amd64-1.5.0.tar.gz'

# create user
sudo adduser --system minecraft

# setup Minecraft server directory
sudo mkdir -p /var/lib/minecraft/
sudo wget -O /var/lib/minecraft/server.jar "$server_url"
sudo chown -R minecraft: /var/lib/minecraft/

# setup minecraftd
wget -O - "$minecraftd_url" | sudo tar -x -z -v -C /

# start minecraftd
sudo systemctl daemon-reload
sudo systemctl start minecraftd.service
```
