# Crypto bots

> Monorepo fr many types of bot

## Features

#### Chain related

- ERC-20 sniper bot main :
    - [ ] Control verbose
    - [x] Send
    - [x] Cancel
    - [x] Update tx
    - [x] Call sc
    - [x] Deploy sc
    - [ ] Add private key to accounts
    - [ ] wallet (locking by sniping fee funds)
    - [ ] smart contract pair (sell, buy)
    - [ ] Block observation
- [ ] Secrets using decentralized secrets store [nucypher](https://www.nucypher.com/)
- [ ] git ignore ropsten wallet

#### Ui

- [ ] Web ui for bot configuration and launch strategy

#### Exchanges related

- [ ] Listing of twitter exchange messages (WIP)
- [ ] Pump and dump telegram groups (WIP)
- [ ] Scalping strat (hammingbot ?)

##### Less prio

- Wallet features
    - [ ] Import account with memonnic using go package btc-util
    - [ ] derivation path account
    - [ ] Refacto transaction (move to model), create a full flow funciton (create / sign / broadcast)

##### Tests

- [ ] Plan what to test
- [ ] Test unit token conversions
- [ ] (e2e) contract calls
- [ ] Typing casts

##### Refacto

- [x] move errors in return to panic and handle with panic
- [ ] Tx adapter for chain with gas or no gas (pow / pos / zk-snarks)

### Backtest config

DEX : 
- [ ] SRM
- [ ] uniswapv3
- [ ] Pancakeswap
- [ ] Pangolin

Nets : 
- [x] (erc20) find testnet connection
- [ ] (bep20) find testnet connection
- k6 for load and e2e tests ?

### Bug

- [ ] Avoid multiple keystores by creating one if missing or fetch most recent (use cmd args)

### Build and deploy

1. Install

*Go* :
    - `curl -sSL https://git.io/g-install | sh -s`
    - `g install 1.17.2`

*Geth* :
    - `C compiler` : Ensure you have clang (Apple) Mingw2 clang (win) or gcc (Linux) 
    - `mage`

- [x] Solidity `pip3 install solc-select` over mage
- [ ] g version manager doc `g install latest` 
- [ ] Mage targets
- [ ] Github actions tests
- [ ] Github actions semantic release monorepo
- [ ] Badges (ci, report, go reference)
- [ ] Check the use of `azblockchain.azure-blockchain`
- [ ] truffle test

### Compatibility

CEX : 

- [ ] binance
- [ ] gate 
- [ ] kucoin 
- [ ] Mexc

DEX :
- [ ] SRM (raydium pools prio)
- [ ] Pangolin (ARC-20)
- [ ] uniswap v2 (pancake)
- [ ] uniswap v3

## Doc

### Structure

1. pkg

    Exposed libraries

1. web

    Vue 3 frontend

1. Internal

1. build

    Magefile parts and dockerfile (.df)

### Development

Used : `https://github.com/thoas/go-funk`

### Testnets

- Bsc : `https://admin.moralis.io/`
- Eth : alchemy
