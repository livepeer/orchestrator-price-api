# Livepeer Price API - API Documentation

An API for curating the off-chain orchestrator pricing information and exposing the data via various endpoints.

<br>

## Getting Started

* Clone the repo: ```git clone https://github.com/livepeer/orchestrator-price-api``` and ```cd``` into it. 

### Using Docker

* Create a file ```.env``` and set environment variables (sample present in ```livepeer-pricing-tool-api/.env.sample```)
* Run your docker daemon.
* Build the docker image: ```docker build --tag api:latest .``` .
* Run Demux: ```docker run -v "$(pwd)"/data:/data -p 9000:9000 --env-file ./.env api:latest```

<br>

## Table of Contents

  - [/orchestratorStats](#orchestratorstats)
    - [Parameters](#parameters)
    - [Examples](#examples)
  - [/priceHistory](#pricehistory)
    - [Parameters](#parameters-1)
    - [Examples](#examples-1)

<br>

### /orchestratorStats

GET request returning an object array consisting of latest statistics for all the orchestrators.

#### Parameters

##### Query Parameters

| Name | Values | Description |
|--|--|--|
| excludeUnavailable | ```True``` [default] / ```False``` | This query parameter exludes the unavailable orchestrators from the list of returned orchestrators |

#### Examples

- Request
    - [https://nyc.livepeer.com/orchestratorStats](https://nyc.livepeer.com/orchestratorStats)
    - [https://nyc.livepeer.com/orchestratorStats?excludeUnavailable=False](https://nyc.livepeer.com/orchestratorStats?excludeUnavailable=False)


- Response

```
[
  {
    "Address": "0xe9e284277648fcdb09b8efc1832c73c09b5ecf59",
    "ServiceURI": "https://livepeer.production-ue1.staked.cloud:8935",
    "LastRewardRound": 1774,
    "RewardCut": 50000,
    "FeeShare": 0,
    "DelegatedStake": 2.0817318069491477e+24,
    "ActivationRound": 1611,
    "DeactivationRound": 1.157920892373162e+77,
    "Active": true,
    "Status": "Registered",
    "PricePerPixel": 27840.575,
    "UpdatedAt": 1591592904
  },
  {
    "Address": "0xda43d85b8d419a9c51bbf0089c9bd5169c23f2f9",
    "ServiceURI": "https://70ea2ff5-004f-4ed6-98b4-60f6c1b48444.livepeer.herd.run:8935",
    "LastRewardRound": 1774,
    "RewardCut": 25000,
    "FeeShare": 500000,
    "DelegatedStake": 1.3120167580279702e+24,
    "ActivationRound": 1620,
    "DeactivationRound": 1.157920892373162e+77,
    "Active": true,
    "Status": "Registered",
    "PricePerPixel": 83862.63,
    "UpdatedAt": 1591592904
  },
  {
    "Address": "0xa5e37e0ba14655e92deff29f32adbc7d09b8a2cf",
    "ServiceURI": "https://7933316d-8d09-4b34-aebe-fe0654a3b7ca.livepeer.herd.run:8935",
    "LastRewardRound": 1774,
    "RewardCut": 50000,
    "FeeShare": 450000,
    "DelegatedStake": 9.355364168371407e+23,
    "ActivationRound": 1611,
    "DeactivationRound": 1.157920892373162e+77,
    "Active": true,
    "Status": "Registered",
    "PricePerPixel": 83862.63,
    "UpdatedAt": 1591592904
  }
]
```

<br/>

### /priceHistory

GET request returning the price history corresponding to an orchestrator in the form of an object list containing timestamps and price per pixel for the respective timestamps.

#### Parameters

##### Path Parameters

- Orchestrator Address ```/priceHistory/:OrchestratorAddress```
    - This path parameter constitutes of the the orchestrator address, for which the price history is required. The addresses can be fetched via [/orchestratorStats](#orchestratorStats) endpoint.

##### Query Parameters

| Name | Values | Description |
|--|--|--|
| limit | ```integer``` | Limits the number of rows in reponse to the specified value |
| offset | ```integer``` | When coupled with ```limit```, can be used for pagination |
| from | ```integer``` | Specifies the start value of a date/time range in unix epoch time format |
| till | ```integer``` | Specifies the end value of a date/time range in unix epoch time format |


#### Examples

- Request
    - [https://nyc.livepeer.com/priceHistory/0xda43d85b8d419a9c51bbf0089c9bd5169c23f2f9](https://livepeer-pricing-tool.com/priceHistory/0xda43d85b8d419a9c51bbf0089c9bd5169c23f2f9)
    - [https://nyc.livepeer.com/priceHistory/0xda43d85b8d419a9c51bbf0089c9bd5169c23f2f9?limit=10&offset=20](https://livepeer-pricing-tool.com/priceHistory/0xda43d85b8d419a9c51bbf0089c9bd5169c23f2f9?limit=10&offset=20)
    - [https://nyc.livepeer.com/priceHistory/0xda43d85b8d419a9c51bbf0089c9bd5169c23f2f9?from=1616218200&till=1616563800](https://livepeer-pricing-tool.com/priceHistory/0xda43d85b8d419a9c51bbf0089c9bd5169c23f2f9?from=1616218200&till=1616563800)


- Response

```
[
  {
    "PricePerPixel": 81836.82,
    "Time": 1591596521
  },
  {
    "PricePerPixel": 83862.63,
    "Time": 1591592904
  },
  {
    "PricePerPixel": 83763,
    "Time": 1591589287
  },
  {
    "PricePerPixel": 83098.8,
    "Time": 1591585672
  },
  {
    "PricePerPixel": 82434.6,
    "Time": 1591582056
  }
]
```



