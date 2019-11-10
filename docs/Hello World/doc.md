# Hello World

## 安装

### 依赖

- protoc-gen-micro

    ```bash
    go get github.com/micro/protoc-gen-micro
    ```

- protoc-gen-go

    ```bash
    brew install protoc-gen-go
    ```

### micro

- go-micro

  可能会出现的问题：
  - [#748](https://github.com/micro/go-micro/issues/748)
  - [#748](https://github.com/micro/go-micro/issues/817)

    ```bash
    export GO111MODULE=on
    go get github.com/micro/go-micro
    ```

- micro

    ```bash
    go get github.com/micro/micro
    ```

    ```bash
    # MacOS
    curl -fsSL https://micro.mu/install.sh | /bin/bash

    # Linux
    wget -q https://micro.mu/install.sh -O - | /bin/bash

    # Windows
    powershell -Command "iwr -useb https://micro.mu/install.ps1 | iex"
    ```

## 第一个服务

[tutorials](https://github.com/micro-in-cn/tutorials)

### Web Dashboard

```bash
micro web
```

[http://127.0.0.1:8082/](http://127.0.0.1:8082/)
