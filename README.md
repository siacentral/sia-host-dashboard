**Sia Host Dashboard is deprecated and replaced by the Sia Foundation's [hostd](https://github.com/SiaFoundation/hostd)**

A simple and powerful dashboard for Sia network Hosts. Combines data for your node with data
from the blockchain to provide accurate and powerful financial and health statistics for hosts.

## Features

+ Secure
+ Easy access from anywhere within local network
+ Monitor Connection Status
+ Smart health alerts
+ View historical data
+ Track storage and data usage
+ Track earned and potential income

![host dashboard](https://siacentral.s3.filebase.com/resources/host-dashboard/host%20dashboard.png)

## Running
1. Download the latest release corresponding to your platform
2. Open command prompt or terminal
3. `cd` to the directory that the dashboard binary is located in
4. Start the dashboard by running `./dashboard`

## Command Line Flags

#### `--data-path`
Changes the path that host data is stored

```
dashboard --data-path /home/ubuntu/data
```

#### `--listen-addr`
Changes the address and port that the dashboard will respond to

```
dashboard --listen-addr localhost:8885
```

#### `--sia-addr`
Changes the address and port of the Sia API the dashboard accesses

```
dashboard --sia-addr localhost:9980
```

## Updating
1. Stop dashboard
2. Download the latest release
3. Replace old dashboard binary with downloaded one
4. Start dashboard

## Docker
The dashboard needs access to a running Sia node. To run the dashboard in a docker container
requires using host networking for Sia and Dashboard or creating a separate docker network so containers can be accessed by name.

### Host Networking
**1. Change Sia to use host networking mode using `--network="host"`**

```
docker run -d \
	--name sia \
	--network="host"\
	-v $(PWD)/data:/sia-data \
	siacentral/sia
```

**2. Start the dashboard container using the same `--network="host"`**

The Sia container name from the last step is `sia` to point the dashboard to the proper location we use `-e SIA_API_ADDR="sia:9980"`

```
docker run -d \
	--name host-dashboard \
	--network="host" \
	-v $(PWD)/data:/data \
	siacentral/host-dashboard
```

### Custom Networking

**1. Create a new docker network**
```
docker network create sia-network
```

**2. Add sia container to the new network by adding the `--network sia-network` flag to your `docker run` command**

```
docker run -d \
	--name sia \
	--network="sia-network"\
	-v $(PWD)/data:/sia-data \
	-p 127.0.0.1:9980:9980 \
	-p 9981:9981 \
	-p 9982:9982 \
	siacentral/sia
```

**3. Start the dashboard container passing in the `SIA_API_ADDR` environment variable using the Sia container's name**

The Sia container name from the last step is `sia` to point the dashboard to the proper location we use `-e SIA_API_ADDR="sia:9980"`

```
docker run -d \
	--name host-dashboard \
	--network="sia-network" \
	-p 8884:8884 \
	-v $(PWD)/data:/data \
	-e SIA_API_ADDR="sia:9980" \
	siacentral/host-dashboard
```

### Docker Compose

Below is a docker-compose example service to run Sia and Host Dashboard. Replace `SIA_API_PASSWORD` and `SIA_WALLET_PASSWORD` with your own passwords and the volume mounts with the correct volume paths

```yml
version: '3.0'
services:
  sia-host:
    image: siacentral/sia:latest
    environment:
      - SIA_API_PASSWORD=asecureapipassword
      - SIA_WALLET_PASSWORD=asecurewalletpassword
    volumes:
      - ./sia-data:/sia-data
      - ./renter-data:/renter-data
    ports:
      - "127.0.0.1:9980:9980"
      - "9981:9981"
      - "9982:9982"
      - "9983:9983"
    restart: unless-stopped
  host-dashboard:
    image: siacentral/host-dashboard:latest
    depends_on:
      - sia-host
    links:
      - sia-host
    environment:
      - SIA_API_ADDR=sia-host:9980
    volumes:
      - ./dashboard-data:/data
    ports:
      - "8884:8884"
    restart: unless-stopped
```

## Development

### Lint and fix
```
make lint
```

### Package and run for development
```
make run
```

### Build for current platform
```
make build
```

### Build for all platforms
```
make release
```
