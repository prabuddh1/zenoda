# zenoda
**zenoda** is a blockchain built using Cosmos SDK (v0.50.11), Tendermint and Ignite CLI (v28.7.0).


## Zenoda Chain Architecture

![Alt Text](/enoda1.png)

## Zenoda Chain Thoery of above Architecture
**Zenoda Chain** is a custom built chain that incorporates several innovative changes as compared to existing Cosmos based chains. Below we describe the flow for **Zenoda** chain (markings done as per the above architecture diagram):

1. Zenoda uses the default validator set and genesis accounts and leverages the power of default x/bank, x/mint, x/auth, x/staking modules from Cosmos-SDK & Ignite Cli.
    **1.1** Defines dynamic paramters for genesis (Set of 10 predefined wallets that will be **governance layer wallets**).
    **1.2** Defines dynamic paramters for genesis (Inflation rate).
    *Example:*
    ### Params

    ```json
    {
    "inflation_rate": "0.05",
    "predefined_wallets": [
        "cosmos1lhahcqzx45mssr9wfknx48hy4truyz9p2wj3ht",
        "cosmos1g6k8qf0zksqruq8exv0duw3p9fn33aeffdprl6",
        "cosmos1fmatdzaqv75sxte07nt67arx398v2cdudv25fl",
        "cosmos1hfc7lhageqsln6kj2e225q6hj2muv27ck7ajzx",
        ...
    ]
    }

2. EGV token is created within the custom x/rewards that will serve for transaction count based governance.

3. Pre-distribution of 1000 EGV tokens to **governance layer wallets**.

4. Transaction tracking (individual & overall network) & EGV reward distribution.
    **[Reward calculated as: (individual_address_transactions / total_network_transactions) * (inflation_rate * total_supply)]**

5. Governance module that handles proposal, voting, upgrades based on network contribution.
    **[Voting weights calculated as: (individual_address_transactions / total_network_transactions)]**

6. Governance upgrade incorporation based on voting results to update parameters like **Inflation Rate & Governance Layer Wallets.**


## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your **Zenoda** blockchain in development.

### Configure

**Zenoda** chain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers both Vue and React options for frontend scaffolding:

For a Vue frontend, use: `ignite scaffold vue`
For a React frontend, use: `ignite scaffold react`
These commands can be run within your scaffolded blockchain project. 

For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

After a draft release is created, make your final changes from the release page and publish it.

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
