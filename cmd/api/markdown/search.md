Endpoint finds entity by hash (block, address, namespace and tx)

### Block

Block will be found by its hash. Hash example: `652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF`.
Hash should be hexadecimal and has a length of 64.

#### Example response 

```json
{
    "type": "block",
    "result": {
        "id": 1,
        "hash": "652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF",
        // ... rest fields from response.Block type
    }
}
```

### Tx

Tx will be found by its hash. Hash example: `652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF`.
Tx should be hexadecimal and has a length of 64.

#### Example response 

```json
{
    "type": "tx",
    "result": {
        "id": 1,
        "hash": "652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF",
        // ... rest fields from response.Tx type
    }
}
```

### Address

The Address will be found by its hash.
Hash example: `celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60`.
Address has prefix `celestia` and has length 47.
Also, it should be decoded by `bech32`.

#### Example response 

```json
{
    "type": "address",
    "result": {
        "id": 1,
        "hash": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
        "height": 100,
        "balance": "6525472354"
    }
}
```

### Namespace

Namespace can be found by base64 hash and identity pair version + namespace id. 
Hash example: `U3dhZ2dlciByb2Nrcw==`. 
Identity pair example: `014723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02`

#### Example response 

```json
{
    "type": "namespace",
    "result": {
        "id": 1,
        "hash": "U3dhZ2dlciByb2Nrcw==",
        "version": 1,
        "namespace_id": "4723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02"
        // ... rest fields from response.Namespace type
    }
}
```
