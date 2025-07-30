# Ravly — Fast Local File Discovery & Transfer Tool

Ravly is a lightweight developer tool for discovering devices and securely transferring files over a local/private network without requiring internet or third-party services. Inspired by tools like AirDrop, Ravly is designed to be cross-platform, terminal-native, and developer-friendly.

⚠️ **Note**: Ravly is in active development. The interface and features are not yet stable, and many capabilities are still in progress.

---

## ✨ Features

- 🔐 **QUIC File Transfers** — Encrypted and fast with TLS.
- 🔎 **Local Peer Discovery** — Uses UDP broadcasts to discover nearby devices.
- 🧠 **gRPC-based Control Plane** — Daemon stores config, handles requests from CLI.
- 📦 **Daemon Mode**— One-time setup, then just use ravly CLI from anywhere.
- 🧪 Easy CLI interface to interact with the running daemon
---

## 🧱 Architecture
```
┌───────────────────────┐       ┌────────────────────────────┐
│       CLI             │       │        Ravly Daemon        │
│                       │ gRPC  │                            │
│  ravly run ...        ├──────►│  - Config Store            │
│  ravly send ...       │       │  - Discovery Service (UDP) │
│  ravly discover       │       │  - Transfer Service (QUIC) │
│                       │       └──────────▲─────────────────┘
└───────────────────────┘                  │
                               ┌───────────┴────────────┐
                               │     System Service     │
                               └────────────────────────┘
```

Ravly is built on a **daemon + CLI** model:

### 🔧 Ravly Daemon

The `ravlyd` daemon is responsible for:

- Configuration Storage (in-memory for now, persistent planned)
- Discovery Service (UDP broadcast to find other Ravly daemons)
- Transfer Server (QUIC-based file receiving service)

**It's designed to run in the background. You can control it through the CLI.**

> The daemon is automatically launched in the background when needed (e.g., via `ravly run`), so the user doesn't have to manually manage it.

### 💻⚙️ How It Works
1. Each device runs `ravly run`, which starts a gRPC daemon.

2. The daemon:
     - Stores configuration (e.g. name, ports, directory).
     - Starts UDP peer discovery (to find other Ravly users on LAN).
     - Listens for incoming files via **QUIC**.

3. You can then run:
   - `ravly discover` to find nearby peers.
   - `ravly send <ip> <file-path>` to send files securely via QUIC.

---

## 🏁 🛠️ CLI Commands
### `ravly run` Starts the Ravly daemon with configurable settings.
```bash
ravly run \
  --name "MacBook" \
  --transfer-port 9898 \
  --discovery-port 9999 \
  --recv-dir ~/Downloads \
  --auto-accept
```
### This does the following:
- Launches the gRPC server for CLI control.
- Starts listening for incoming transfers on QUIC port.
- Broadcasts presence over UDP (discovery).
- Applies and stores the given config in-memory.

| Flag               | Description                         | Default       |
| ------------------ | ----------------------------------- | ------------- |
| `--name`           | Display name of the device          | hostname      |
| `--transfer-port`  | Port to receive files over QUIC     | `9898`        |
| `--discovery-port` | UDP port for peer discovery         | `9999`        |
| `--recv-dir`       | Directory to save incoming files    | `~/Downloads` |
| `--auto-accept`    | Automatically accept incoming files (todo) | `false`       |


### `ravly discover` Scans the local network for other Ravly peers.
```bash
ravly discover
```
Output:
```bash
🌟 Ravly Peer Discovery
-----------------------
Discovered peers:
🏠 10.0.0.23   MacBook-Pro   (seen 2s ago)
```
Discovery is done using UDP broadcasts. All devices running ravly run on the same network will appear.

### `ravly send <ip> <file>` Transfers a file to another Ravly peer using QUIC.
```bash
ravly send 10.0.0.23 ./photo.jpg
```



