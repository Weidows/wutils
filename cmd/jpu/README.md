# Jetbrains Portable Upgrader

通过改配置实现 Portable 效果

```
- PyCharm-P
    - ch-0
        - 223.8214.51
            - bin
                - idea.properties
        - 223.8214.51.plugins
        - 223.8617.48
        - 223.8617.48.plugins
    - config
    - system
- Goland
- datagrip
```

在 `IDE/bin/idea.properties` 中添加

```properties
idea.config.path=${idea.home.path}/../../config
idea.system.path=${idea.home.path}/../../system
```

由于路径中含带版本号, 用脚本不易操作, 所以用 go 写