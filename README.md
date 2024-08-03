[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-0969da?style=flat-square&logo=opensource)](https://opensource.org/licenses/MIT)

```bash
     _    ____   ___  _        ____ _     ___
    / \  / ___| / _ \| |      / ___| |   |_ _|
   / _ \ \___ \| | | | |     | |   | |    | |
  / ___ \ ___) | |_| | |___  | |___| |___ | |
 /_/   \_\____/ \___/|_____|  \____|_____|___|

```

Welcome to **asol**, a CLI for interacting with the Solana blockchain 🎉

## Installation

```bash
# Clone the repository
git clone https://github.com/Almazatun/asol
# Enter into the directory
cd asol/
# Install the dependencies
$ make install
```

```bash
# Build
$ make build
# Installation binary to [/go/bin/asol]
$ go install
```

## 🌟 Features

```bash
# Create accounts in the Solana blockchain
$ asol wallet OR asol wallet --list=<MAX=100>
# Get account balance
$ asol balance
# Transfer SOL from one account to another [{ "publicKey": "some_address", "amount": "some_amount" }, ...]
$ asol transfer OR asol transfer --path=<json file with list of accounts>
```
