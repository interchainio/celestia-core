# Tendermint and Celestia

celestia-core is not meant to be used as a general purpose framework.
Instead, its main purpose is to provide certain components (mainly consensus but also a p2p layer for Tx gossiping) for the Celestia main chain.
Hence, we do not provide any extensive documentation here.

Instead of keeping a copy of the Tendermint documentation, we refer to the existing extensive and maintained documentation and specification:

 - https://docs.tendermint.com/
 - https://github.com/tendermint/tendermint/tree/master/docs/
 - https://github.com/tendermint/spec

Reading these will give you a lot of background and context on Tendermint which will also help you understand how celestia-core and [celestia-app](https://github.com/celestiaorg/celestia-app) interact with each other.

# Celestia

As mentioned above, celestia-core aims to be more focused on the Celestia use-case than vanilla Tendermint.
Moving forward we might provide a clear overview on the changes we incorporated.
For now, we refer to the Celestia specific [ADRs](./adr) in this repository as well as to the Celestia specification:

 - [celestia-adr](./adr)
 - [celestia-specs](https://github.com/celestiaorg/celestia-specs)
