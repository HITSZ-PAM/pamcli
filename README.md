# pamcli

## 使用

首先在pam上注册应用，并分配读权限到具体的密码保险箱

```
export PAM_ADDR=https://10.95.58.126  
export PAM_CLIENT_ID=I6Kyx9jW798ALzw9yt6O0Ua4sLUAqlU0  
export PAM_CLIENT_SECRET=pc9OSe5mR4D0mWpO1uq5qlj2I38e6CZ0q8ekMUUFQ6G1DUdK1nP052f6yznLM8sgTM2EBT9GIc5eWw3D37Ifkh2012812eiB272fJhpv488EdEqtsBl9LILB7yMMmbE1  
```

然后按pamcli://username/<accountid>、pamcli://password/<accountid>格式注册变量

```
export ORACLE_ACCOUNT_USERNAME=pamcli://username/6825644932833738753  
export ORACLE_ACCOUNT_PASSWORD=pamcli://password/6825644932833738753  
export LINUX_ACCOUNT_USERNAME=pamcli://username/6825638948082024449  
export LINUX_ACCOUNT_PASSWORD=pamcli://password/6825638948082024449  
```

执行应用程序  
./pamcli run -- sh example.sh  