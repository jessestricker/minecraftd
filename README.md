# minecraftd

A Minecraft server daemon.

## Setup

```bash
# URL of Minecraft server JAR file, see https://www.minecraft.net/download/server
server_url='https://piston-data.mojang.com/v1/objects/4707d00eb834b446575d89a61a11b5d548d8c001/server.jar'

# URL of minecraftd, see https://github.com/jessestricker/minecraftd/releases/latest
minecraftd_url='https://github.com/jessestricker/minecraftd/releases/download/v1.1.1/minecraftd-1.1.1.tar.gz'

sudo adduser --system minecraft
sudo mkdir -p /var/lib/minecraft/
sudo wget -O /var/lib/minecraft/server.jar "$server_url"
sudo chown -R minecraft: /var/lib/minecraft/

pushd $(mktemp -d)
wget -O minecraftd.tar.gz "$minecraftd_url"
sudo tar -xzvf minecraftd.tar.gz -C /
popd

sudo systemctl daemon-reload
sudo systemctl enable --now minecraft.service
```
