var request = require( 'request' );

function main( params ) {
  return new Promise( function( resolve, reject ) {
    request( {
      method: 'POST',
      url: 'http://natty.joestelmach.com/parse',
      form: {
        value: params.value
      }
    }, function( err, result, body ) {
      resolve( {
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.parse( body )
      } );
    } )
  } );
}
