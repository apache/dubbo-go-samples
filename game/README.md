# Game Service Example

### Backend

- The example includes two services, **gate** (gateway service) and **game** (logical service)
- The two services communicate with each other RPC (both registered **provider** and **consumer**)
- The **gate** additionally starts the http service (port **8000**), which is used to manually trigger the **gate** RPC to call **game**

> After each **gate** RPC call (**Message**) **game**, **game** will synchronize the RPC call (Send) **gate** pushes the same message
> This logic is annotated by default, and the location is `go-server-game/pkg/provider.go` (line 25 ~ 30)

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.
