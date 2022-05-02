# Bastion

Demo Bastion

**Server**

```shell
go run cmd/server/main.go -server 0.0.0.0:8043
```

**Agent**

```shell
go run cmd/agent/agent.go -server 0.0.0.0:8043
```

**Client**

```shell
go run cmd/cli/main.go -d <YOUR_AGENT_HOSTNAME> -t /bin/bash
```