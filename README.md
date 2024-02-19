# bc-migration-tool

This project is a tool designed to assist BC(BNB beacon chain) validators in migrating to BSC(BNB smart chain) after BC-fusion upgrade. 

The migration process generally involves the following stages: 
1. creating a validator on the BSC side;
2. migrate delegations from BC side to BSC side.

## Prerequisites
- go 1.21+

## Installation

### Install this tool

```shell
git clone https://github.com/bnb-chain/bc-migration-tool.git
cd bc-migration-tool
go build -o bin/bc-migration-tool main.go
```

### Install `bnbcli`

For more information of the installation and usage of `bnbcli`, you can also refer to [bnbcli](https://github.com/bnb-chain/node/blob/master/README.md).

#### Option 1

Download the binaries which are suitable for your platform from the [latest release](https://github.com/bnb-chain/node/releases/latest).

#### Option 2

```shell
git clone https://github.com/bnb-chain/node.git
cd node
make build
cp build/bnbcli ${workspace}/bin
```

### Install `geth`

#### Option 1

Download the binaries which are suitable for your platform from the [latest release](https://github.com/bnb-chain/bsc/releases/latest).

#### Option 2

```shell
git clone https://github.com/bnb-chain/bsc.git
cd bsc
make geth
cp build/bin/geth ${workspace}/bin
```

## Preparation

Before you start, you need to prepare your BC account and BSC account.

### 1. Set up BC account

We require you to import your BC account into `bnbcli`. If you have already set up the account, you can skip this step.

#### Existing account

Your account info will be stored at `~/.bnbcli/keys` by default. You can check your account info by running the following command:

```shell
${workspace}/bin/bnbcli keys show <your-account-name> --home ${HOME}/.bnbcli
```

#### Import account by mnemonic

If you have mnemonic, you can import your account by running the following command:

```shell
$ ${workspace}/bin/bnbcli keys add <your-account-name> --recover --home ${HOME}/.bnbcli
Enter a passphrase for your key:
Repeat the passphrase:
> Enter your recovery seed phrase:
```

You will be asked to set a password for this account and input your mnemonic. After that, you will get your account info.

#### Import account from ledger

If you have a ledger, you can import your account by running the following command:

```shell
${workspace}/bin/bnbcli keys add <your-account-name> --ledger --home ${HOME}/.bnbcli
```

### 2. Set up BSC account

You can set up your BSC account by providing keystore file or private key lately. 

### 3. Set up BLS account

To create a validator, you need to provide a BLS address. And this address should not be the same as the one you ever used on the BC side.

You can create a new BLS account by running the following command:

```shell
${workspace}/bin/geth bls account new --datadir ${datadir}
```

If you don't have a wallet folder under `${datadir}`, this will help you create a new wallet at first. The wallet file will be stored at `${datadir}/bls`.

And if you have already created a wallet, you need to provide the password for your wallet.

After that, you will get your BLS account info like this:

```shell
INFO[0003] Successfully imported validator key(s)        prefix=local-keymanager publicKeys=0x000000000000
```

The key info is `0x000000000000` in this case.

## Usage

### Create a validator on the BSC side

**Note**
If you are an old validator operator on the BC side, please make a validator mapping signature by following steps.
This can help the user to verify that the validator on the BSC side is the same as the one on the BC side.
And a user can redelegate to the new validator on the BSC side without waiting for the unbonding period via 
cross-chain redelegation (refer to the following sections).

#### Local Key
```shell
${workspace}/bin/bnbcli \
  validator-ownership \
  sign-validator-ownership \
  --bsc-operator-address ${NEW_VALIDATOR_OPERATOR_ADDR_ON_BSC} \
  --from ${ACCOUNT_NAME} \
  --chain-id ${BC_CHAIN_ID} \
```

#### Ledger Key
```shell
${workspace}/bin/bnbcli \
  validator-ownership \
  sign-validator-ownership \
  --bsc-operator-address ${NEW_VALIDATOR_OPERATOR_ADDR_ON_BSC} \
  --from ${BSC_OPERATOR_NAME} \
  --chain-id ${CHAIN_ID} \
  --ledger

```

- `${workspace}/bin/bnbcli`: The path to the `bnbcli` binary executable.

- `--to ${NEW_VALIDATOR_OPERATOR_ADDR_ON_BSC}`: Specifies the BSC address to which the new validator operator address will be mapped.

- `--chain-id ${BC_CHAIN_ID}`: Specifies the chain ID for the BC(BNB beacon chain). By default, the mainnet chain ID is `Binance-Chain-Tigris`.

- `--from ${ACCOUNT_NAME}`: Specifies the account name from which the sign will be performed.

And you will get the output like this:
```
TX JSON: {"type":"auth/StdTx","value":{"msg":[{"type":"migrate/ValidatorOwnerShip","value":{"bsc_operator_address":"RXN7r5XZlaljqzp8msZvx6Y6124="}}],"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Ahr+LlBMLgiUFkP75kIuJW1YHrsTy39GeOdV+IaTREDN"},"signature":"AL5mj52s0+tcdoEb6c6PAmqBixuv3XEmrLW3Y1kvUeYgG3RqVvWU/dIVcfxiHHwLGXlcn0X1v00jFrpLIsxtqA==","account_number":"0","sequence":"0"}],"memo":"","source":"0","data":null}}
Sign Message:  {"account_number":"0","chain_id":"Binance-GGG-Ganges","data":null,"memo":"","msgs":[{"bsc_operator_address":"0x45737baf95d995a963ab3a7c9ac66fc7a63ad76e"}],"sequence":"0","source":"0"}
Sign Message Hash:  0x8f7179e7969e497b5f3c006535e55c2fa5bea5d118a8008eddce3fccd1675673
Signature: 0x00be668f9dacd3eb5c76811be9ce8f026a818b1bafdd7126acb5b763592f51e6201b746a56f594fdd21571fc621c7c0b19795c9f45f5bf4d2316ba4b22cc6da8
PubKey: 0x021afe2e504c2e08941643fbe6422e256d581ebb13cb7f4678e755f886934440cd
```
The `Signature` is your `VALIDATOR_OWNER_SHIP_SIGNATURE`


First, you need to modify the `config/config.yml` file.
```yaml
BscRpcUrl: "https://bsc-dataseed.binance.org/"
BlsDataDir: "~/.geth/bls/wallet"
ValidatorInfo: {
    Delegation: "2000000000000000000000", // no less than 2000 BNB
    ConsensusAddress: "0x0000000000000000000000000000000000000000", // cannot be the same as the one you ever used on the BC side
    Description: {
        "moniker": "moniker", // only support alphanumeric characters and the length should be between 3 and 9
        "identity": "${VALIDATOR_OWNER_SHIP_SIGNATURE}", // if you are a old validator operator on the BC side, please make a validator mapping signature. or you can leave it blank or any string
        "website": "website",
        "details": "details"
    },
    Commission: {
        "rate": 0, // 10000 == 100%
        "maxRate": 5000, // 10000 == 100% and no more than 5000
        "maxChangeRate": 5000 // 10000 == 100% and no more than 5000
    }
}
```

- `BscRpcUrl`: The RPC URL of the BSC node. e.g. `https://bsc-dataseed.binance.org/`.
- `BlsDataDir`: The path to the BLS wallet folder. e.g. `~/.geth/bls/wallet`.
- `ValidatorInfo`: The validator info you want to set on the BSC side. You can refer to [here](https://docs.binance.org/smart-chain/validator/validator.html#validator-info) for more information.

When you are ready, you can create a validator on the BSC side by running the following command:

```shell
${workspace}/bin/bc-migration-tool create-validator \
    --bls-password ${BLS_PASSWORD} \
    --bls-pubkey ${BLS_PUBKEY} \
    --operator-account ${OPERATOR_ACCOUNT} \
    --ledger/--private-key/--keystore-path \
    --index ${LEDGER_ACCOUNT_INDEX}
```

- `${workspace}/bin/bc-migration-tool`: The path to the `bc-migration-tool` binary executable.
- `--bls-password ${BLS_PASSWORD}`: The password for your BLS wallet.
- `--bls-pubkey ${BLS_PUBKEY}`: The BLS public key you want to use.
- `--operator-account ${OPERATOR_ACCOUNT}`: The BSC address you want to use as the operator address.
- `--ledger/--private-key/--keystore-path`: The way you want to provide your BSC account. You can choose one of them.
- `--index`: To specify the account you want to use on ledger. Default is zero, means your first account.

### Migrate delegations from BC side to BSC side

There are two options to migrate delegations from the BC side to the BSC side.
- Option 1: Cross-chain redelegation (recommended)
- Option 2: Undelegation and cross-chain transfer

#### Option 1: Cross-chain redelegation

You can migrate your delegations on the BC side to BSC side without an unbonding period (7 days) by running the following command:

```shell
${workspace}/bin/bnbcli staking bsc-stake-migration \
    --chain-id ${BC_CHAIN_ID} \
    --side-chain-id ${BSC_CHAIN_NAME} \
    --from ${ACCOUNT_NAME} \
    --validator ${VALIDATOR_ADDR} \
    --address-smart-chain-validator ${BSC_VALIDATOR_ADDR} \
    --address-smart-chain-beneficiary ${BSC_BENEFICIARY_ADDR} \
    --amount ${AMOUNT}:BNB \
    --node ${BC_NODE_URL} --trust-node \
    --home ${HOME}/.bnbcli
```

- `${workspace}/bin/bnbcli`: The path to the `bnbcli` binary executable.

- `--chain-id ${BC_CHAIN_ID}`: Specifies the chain ID for the BC(BNB beacon chain). By default, the mainnet chain ID is `Binance-Chain-Tigris`.

- `--side-chain-id ${BSC_CHAIN_NAME}`: Specifies the chain ID for the BSC(BNB Smart Chain). By default, the mainnet chain ID is `bsc`.

- `--from ${ACCOUNT_NAME}`: Specifies the account name from which the unbonding operation will be performed.

- `--validator ${VALIDATOR_ADDR}`: Specifies the validator address for which the unbonding operation is intended.

- `--address-smart-chain-validator ${BSC_VALIDATOR_ADDR}`: Specifies the BSC validator address for which to delegate to.

- `--address-smart-chain-beneficiary ${BSC_BENEFICIARY_ADDR}`: Specifies the BSC delegator address as the beneficiary address. 
   Note make sure the address is correct, otherwise the fund will be lost forever.

- `--amount ${AMOUNT}:BNB`: Specifies the amount to be unbonded, along with the asset type (BNB in this case). Note that the decimal is 10.

- `--node ${BC_NODE_URL} --trust-node`: Specifies the BC node URL and instructs `bnbcli` to trust the node.

- `--home ${HOME}/.bnbcli`: Specifies the bnbcli home directory where the key file imported above is stored.

You will be asked to input the password for your account.

Be noted: As a validator operator, you should inform your delegators for migrations as well, as they will still be bonded to your validator.

#### Option 2: Undelegation and cross-chain transfer

##### Undelegate

You can undelegate your delegations on the BC side by running the following command:

```shell
${workspace}/bin/bnbcli staking bsc-unbond \
    --chain-id ${BC_CHAIN_ID} \
    --side-chain-id ${BSC_CHAIN_NAME} \
    --from ${ACCOUNT_NAME} \
    --validator ${VALIDATOR_ADDR} \
    --amount ${AMOUNT}:BNB \
    --node ${BC_NODE_URL} --trust-node \
    --home ${HOME}/.bnbcli
```

- `${workspace}/bin/bnbcli`: The path to the `bnbcli` binary executable.

- `--chain-id ${BC_CHAIN_ID}`: Specifies the chain ID for the BC(BNB beacon chain). By default, the mainnet chain ID is `Binance-Chain-Tigris`.

- `--side-chain-id ${BSC_CHAIN_NAME}`: Specifies the chain ID for the BSC(BNB Smart Chain). By default, the mainnet chain ID is `bsc`.

- `--from ${ACCOUNT_NAME}`: Specifies the account name from which the unbonding operation will be performed.

- `--validator ${VALIDATOR_ADDR}`: Specifies the validator address for which the unbonding operation is intended.

- `--amount ${AMOUNT}:BNB`: Specifies the amount to be unbonded, along with the asset type (BNB in this case). Note that the decimal is 10.

- `--node ${BC_NODE_URL} --trust-node`: Specifies the BC node URL and instructs `bnbcli` to trust the node.

- `--home ${HOME}/.bnbcli`: Specifies the bnbcli home directory where the key file imported above is stored.

You will be asked to input the password for your account.

Be noted: As a validator operator, you should inform your delegators for manually unbond, as they will still be bonded to your validator.

#####  Cross-chain transfer

After the unbonding period (7 days), you can cross-chain transfer the funds to BSC by running the following command:

```shell
${workspace}/bin/bnbcli bridge transfer-out \
    --to ${BSC_ADDRESS} \
    --chain-id ${BC_CHAIN_ID} \
    --from ${ACCOUNT_NAME} \
    --amount ${AMOUNT}:BNB \
    --expire-time ${EXPIRE_TIME} \
    --node ${BC_NODE_URL} --trust-node \
    --home ${HOME}/.bnbcli
```

- `${workspace}/bin/bnbcli`: The path to the `bnbcli` binary executable.

- `--to ${BSC_ADDRESS}`: Specifies the BSC address to which the funds will be transferred. Usually, it should be your operator address on the BSC side.

- `--chain-id ${BC_CHAIN_ID}`: Specifies the chain ID for the BC(BNB beacon chain). By default, the mainnet chain ID is `Binance-Chain-Tigris`.

- `--from ${ACCOUNT_NAME}`: Specifies the account name from which the unbonding operation will be performed.

- `--amount ${AMOUNT}:BNB`: Specifies the amount to be unbonded, along with the asset type (BNB in this case). Note that the decimal is 10.

- `--expire-time ${EXPIRE_TIME}`: Specifies the expiry time(unix) of the cross-chain transfer. The expiration time should be greater than the current time.

- `--node ${BC_NODE_URL} --trust-node`: Specifies the BC node URL and instructs `bnbcli` to trust the node.

- `--home ${HOME}/.bnbcli`: Specifies the bnbcli home directory where the key file imported above is stored.

You will be asked to input the password for your account.



## Contributing
