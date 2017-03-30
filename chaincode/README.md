go build ./


{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID": {
      "path": "https://github.com/krhoyt/Blockhead/todo/chaincode"
    },
    "ctorMsg": {
      "function": "init",
      "args": [
        "init"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 0
}

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "ae85a2bd7da4eaa2cf99ccabeac6099cab202c7cfd90824df2c96a0c629a2dd56814b3c74a4c1207b61007779725e0cfc7f362f8d1138bd45796197c3308d68a"
    },
    "ctorMsg": {
      "function": "create_account",
      "args": [
      	"abc-123", 
      	"krhoyt", 
      	"abc123"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 0
}

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "ae85a2bd7da4eaa2cf99ccabeac6099cab202c7cfd90824df2c96a0c629a2dd56814b3c74a4c1207b61007779725e0cfc7f362f8d1138bd45796197c3308d68a"
    },
    "ctorMsg": {
      "function": "get_accounts",
      "args": [
        "all"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 0
}

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "ae85a2bd7da4eaa2cf99ccabeac6099cab202c7cfd90824df2c96a0c629a2dd56814b3c74a4c1207b61007779725e0cfc7f362f8d1138bd45796197c3308d68a"
    },
    "ctorMsg": {
      "function": "find_account",
      "args": [
        "krhoyt",
        "abc123"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 0
}