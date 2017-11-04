[![CircleCI](https://circleci.com/gh/jancajthaml/fio-bco.svg?style=svg&circle-token=dca7fe834e3de7b35f226069ae4729e283ff1df5)](https://circleci.com/gh/jancajthaml/fio-bco)

### How to run it

#### local lifecycle

To run fio-sync you have to install npm and nodeJS. Then run commands below, user has to have write permissions for working directory.

```
npm install

npm start <tenant_name> <tenant_accountIban> <fio_token> [wait]
```

- <tenant_name> - name of tenant in core
- <tenant_accountIban> - iban of the account that is synced to transaction core
- <fio_token> - token that is used to access account via FIO api, read-only token is sufficient
- [wait] - optional just type wait as last argument if you want to wait for fio api to be available otherwise fio-sync will end

#### dockerized lifecycle

```
make

TENANT_NAME=<tenant_name> ACCOUNT_IBAN=<tenant_accountIban> FIO_TOKEN=<fio_token> \
make run
```

- <tenant_name> - name of tenant in core
- <tenant_accountIban> - iban of the account that is synced to transaction core
- <fio_token> - token that is used to access account via FIO api, read-only token is sufficient

### What it does

Application get all transactions for given account via FIO api and store it in core. It also save the last
transaction that was synced, so when run again it get from FIO only transactions that are new.

> Note: There is a FIO limitation that request for account statement from which is transactions gathered can be requested
once per 20 seconds. If you run application twice in that window app will simply wait 20 seconds and then continue.

// TODO: remove accountNumber as input argument and get it from fio account statement, also need to store token in db
