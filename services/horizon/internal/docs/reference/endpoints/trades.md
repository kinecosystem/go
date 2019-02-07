---
title: Trades
---

People on the Stellar network can make [offers](../resources/offer.md) to buy or sell assets. When an offer is fully or partially fulfilled, a [trade](../resources/trade.md) happens.

Trades can be filtered for specific orderbook, defined by an asset pair: `base` and `counter`. 

This endpoint can also be used in [streaming](../streaming.md) mode, making it possible to listen for new trades as they occur on the Stellar network.
If called in streaming mode Horizon will start at the earliest known trade unless a `cursor` is set. In that case it will start from the `cursor`. You can also set `cursor` value to `now` to only stream trades created since your request time.

## Request

```
GET /trades?base_asset_type={base_asset_type}&base_asset_code={base_asset_code}&base_asset_issuer={base_asset_issuer}&counter_asset_type={counter_asset_type}&counter_asset_code={counter_asset_code}&counter_asset_issuer={counter_asset_issuer}&resolution={resolution}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `base_asset_type` | optional, string | Type of base asset | `native` |
| `base_asset_code` | optional, string | Code of base asset, not required if type is `native` | `USD` |
| `base_asset_issuer` | optional, string | Issuer of base asset, not required if type is `native` | 'GA2HGBJIJKI6O4XEM7CZWY5PS6GKSXL6D34ERAJYQSPYA6X6AI7HYW36' |
| `counter_asset_type` | optional, string | Type of counter asset  | `credit_alphanum4` |
| `counter_asset_code` | optional, string | Code of counter asset, not required if type is `native` | `BTC` |
| `counter_asset_issuer` | optional, string | Issuer of counter asset, not required if type is `native` | 'GD6VWBXI6NY3AOOR55RLVQ4MNIDSXE5JSAVXUTF35FRRI72LYPI3WL6Z' |
| `offer_id` | optional, string | filter for by a specific offer id | `283606` |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order, in terms of timeline, in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request
```sh 
curl https://horizon.stellar.org/trades?base_asset_type=native&counter_asset_code=SLT&counter_asset_issuer=GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP&counter_asset_type=credit_alphanum4&limit=2&order=desc
```

## Response

The list of trades. `base` and `counter` in the records will match the asset pair filter order. If an asset pair is not specified, the order is arbitrary.

### Example Response
```json
{
  "_links": {
    "self": {
      "href": "https://horizon.stellar.org/trades?base_asset_type=native\u0026counter_asset_code=SLT\u0026counter_asset_issuer=GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP\u0026counter_asset_type=credit_alphanum4\u0026cursor=\u0026limit=2\u0026order=desc"
    },
    "next": {
      "href": "https://horizon.stellar.org/trades?base_asset_type=native\u0026counter_asset_code=SLT\u0026counter_asset_issuer=GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP\u0026counter_asset_type=credit_alphanum4\u0026cursor=68836785177763841-0\u0026limit=2\u0026order=desc"
    },
    "prev": {
      "href": "https://horizon.stellar.org/trades?base_asset_type=native\u0026counter_asset_code=SLT\u0026counter_asset_issuer=GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP\u0026counter_asset_type=credit_alphanum4\u0026cursor=68836918321750017-0\u0026limit=2\u0026order=asc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "base": {
            "href": "https://horizon.stellar.org/accounts/GBZXCJIUEPDXGHMS64UBJHUVKV6ETWYOVHADLTBXJNJFUC7A7RU5B3GN"
          },
          "counter": {
            "href": "https://horizon.stellar.org/accounts/GBHKUQDYXGK5IEYORI7DZMMXANOIEHHOF364LNT4Q7EWPUL7FOO2SP6D"
          },
          "operation": {
            "href": "https://horizon.stellar.org/operations/68836918321750017"
          }
        },
        "id": "68836918321750017-0",
        "paging_token": "68836918321750017-0",
        "ledger_close_time": "2018-02-02T00:20:10Z",
        "base_offer_id": "695254",
        "base_account": "GBZXCJIUEPDXGHMS64UBJHUVKV6ETWYOVHADLTBXJNJFUC7A7RU5B3GN",
        "base_amount": "0.1217566",
        "base_asset_type": "native",
        "counter_offer_id": "O30064775169",        
        "counter_account": "GBHKUQDYXGK5IEYORI7DZMMXANOIEHHOF364LNT4Q7EWPUL7FOO2SP6D",
        "counter_amount": "0.0199601",
        "counter_asset_type": "credit_alphanum4",
        "counter_asset_code": "SLT",
        "counter_asset_issuer": "GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP",
        "base_is_seller": true,
        "price": {
          "N": 10,
          "D": 61
        }
      },
      {
        "_links": {
          "base": {
            "href": "https://horizon.stellar.org/accounts/GCUODCZAU6CSXEKKWZZNWQXDITIWLWCDK6M4IZ7H5PACLC3QAWEJSOTR"
          },
          "counter": {
            "href": "https://horizon.stellar.org/accounts/GBHKUQDYXGK5IEYORI7DZMMXANOIEHHOF364LNT4Q7EWPUL7FOO2SP6D"
          },
          "operation": {
            "href": "https://horizon.stellar.org/operations/68836785177763841"
          }
        },
        "id": "68836785177763841-0",
        "paging_token": "68836785177763841-0",
        "ledger_close_time": "2018-02-02T00:18:00Z",
        "base_offer_id": "695244",
        "base_account": "GCUODCZAU6CSXEKKWZZNWQXDITIWLWCDK6M4IZ7H5PACLC3QAWEJSOTR",
        "base_amount": "0.0000050",
        "base_asset_type": "native",
        "counter_offer_id": "4611686044197195777",   
        "counter_account": "GBHKUQDYXGK5IEYORI7DZMMXANOIEHHOF364LNT4Q7EWPUL7FOO2SP6D",
        "counter_amount": "0.0000009",
        "counter_asset_type": "credit_alphanum4",
        "counter_asset_code": "SLT",
        "counter_asset_issuer": "GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP",
        "base_is_seller": false,
        "price": {
          "N": 2,
          "D": 11
        }
      }
    ]
  }
}
```

## Example Streaming Event
```
{ 
  _links: 
   { self: { href: '' },
     base: 
      { href: '/accounts/GDJNMHET4DTS7HUHU7IG5DB274OSMHUYA7TRRKOD6ZABHPUW5YWJ4SUD' },
     counter: 
      { href: '/accounts/GCALYDRCCJEUPMV24TAX2N2N3IBX7NUUYZNM7I5FQS5GIEQ4A7EVKUOP' },
     operation: { href: '/operations/47261068505915393' } },
  id: '47261068505915393-0',
  paging_token: '47261068505915393-0',
  ledger_close_time: '2018-09-11T19:42:04Z',
  offer_id: '734529',
  base_account: 'GDJNMHET4DTS7HUHU7IG5DB274OSMHUYA7TRRKOD6ZABHPUW5YWJ4SUD',
  base_amount: '0.0175999',
  base_asset_type: 'credit_alphanum4',
  base_asset_code: 'BOC',
  base_asset_issuer: 'GCTS32RGWRH6RJM62UVZ4UT5ZN5L6B2D3LPGO6Z2NM2EOGVQA7TA6SKO',
  counter_account: 'GCALYDRCCJEUPMV24TAX2N2N3IBX7NUUYZNM7I5FQS5GIEQ4A7EVKUOP',
  counter_amount: '0.0199998',
  counter_asset_type: 'credit_alphanum4',
  counter_asset_code: 'ABC',
  counter_asset_issuer: 'GCTS32RGWRH6RJM62UVZ4UT5ZN5L6B2D3LPGO6Z2NM2EOGVQA7TA6SKO',
  base_is_seller: true,
  price: { n: 2840909, d: 2500000 }
}
```
## Possible Errors

- The [standard errors](../errors.md#Standard_Errors).
