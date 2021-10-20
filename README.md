# Crypto bots

> Monorepo fr many types of bot

## Features

- Listing of twitter exchange messages (WIP)
- Pump and dump telegram groups (WIP)
- [ ] Scalping
- ERC-20 sniper bot :
    - [ ] Send / call SC
    - [ ] wallet (locking by sniping fee funds)
    - [ ] smart contract pair (sell, buy)
- [ ] Secrets using decentralized secrets store [nucypher](https://www.nucypher.com/)

### Backtest config

DEX : SRM / uniswapv3
- [ ] find testnet connection
- [ ] pull ERC-20 data
- k6 for load and e2e tests ?

### Build and deploy

1. Install

*Geth* :
    - `C compiler` : clang or gcc 
    - `Go` : g install latest (1.17.2)
    - `mage`

*Solidity*: `pip3 install solc-select && solc-select install 0.5.16 && solc-select use 0.5.16`

- [ ] Solidity `pip3 install solc-select`
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
