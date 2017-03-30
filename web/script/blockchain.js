class Blockchain {

  static request( route ) {
    return new Promise( ( resolve, reject ) => {
      let xhr = new XMLHttpRequest();
      xhr.addEventListener( 'error', evt => {
        reject( xhr.statusText );
      } );
      xhr.addEventListener( 'load', evt => {
        let data = JSON.parse( xhr.responseText );
        let result = null;

        try {
          result = JSON.parse( data.result.message );
        } catch( e ) {
          result = data.result.message;
        }

        resolve( result );
      } );
      xhr.open( 'POST', Blockchain.URL, true );
      xhr.setRequestHeader( 'Content-Type', 'application/json' );
      xhr.send( JSON.stringify( {
        jsonrpc: '2.0',
        method: route.method,
        params: {
          chaincodeID: {
            name: Blockchain.CHAINCODE
          },
          ctorMsg: {
            function: route.operation,
            args: route.values
          },
          secureContext: Blockchain.USER,
          type: 1
        },
        id: 1        
      } ) );
    } );
  }
  
}

Blockchain.CHAINCODE = '9e675fb998c004215503fe445081cd9213a86e2951b13c0dbbce3fa8dd21cbd38dccb3689e3c574677a01087a93ac705cd34380acd6249a5be2ca18c900193ab';
Blockchain.URL = 'https://ba5d025f326946a1a17a37c120fd233c-vp0.us.blockchain.ibm.com:5002/chaincode';
Blockchain.USER = 'user_type1_0';

Blockchain.QUERY = 'query';
Blockchain.INVOKE = 'invoke';
