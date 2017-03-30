package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  "strconv"  
)

// Task duration
const DURATION_ANY  = 0
const DURATION_30   = 1
const DURATION_60   = 2
const DURATION_2    = 3
const DURATION_4    = 4
const DURATION_LONG = 5

// Task energy level
const ENERGY_ANY    = 0
const ENERGY_LOW    = 1
const ENERGY_NORMAL = 2
const ENERGY_HIGH   = 3

// Shim
type  SimpleChaincode struct {
}

// Account
type Account struct {
  Id       string `json:"id"`
  First    string `json:"first"`
  Last     string `json:"last"`
  Name     string `json:"name"`
  Password string `json:"password"`
}

// Location
type Location struct {
  Id        string `json:"id"`
  AccountId string `json:"account"`
  Name      string `json:"name"`
}

// To do task
type Task struct {
  Id         string `json:"id"`
  AccountId  string `json:"account"`
  Due        int    `json:"due"`
  LocationId string `json:"location"`
  Duration   int    `json:"duration"`
  Energy     int    `json:"energy"`
  Tags       string `json:"tags"`
  Notes      string `json:"notes"`
  Complete   bool   `json:"complete"`
  Name       string `json:"name"`
  CreatedAt  int    `json:"created"`
}

// *
// Initialize
// *

func ( t *SimpleChaincode ) Init( stub shim.ChaincodeStubInterface, function string, args []string ) ( []byte, error ) {
  // Accounts
  var accounts []Account

  bytes, err := json.Marshal( accounts )

  if err != nil { 
    return nil, errors.New( "Error initializing accounts." ) 
  }

  err = stub.PutState( "toodles_accounts", bytes )

  // Locations
  // TODO: Empty array versus nil
  var locations []string

  bytes, err = json.Marshal( locations )

  if err != nil { 
    return nil, errors.New( "Error initializing locations." ) 
  }

  err = stub.PutState( "toodles_locations", bytes )

  // Tasks
  var tasks []string

  bytes, err = json.Marshal( tasks )

  if err != nil { 
    return nil, errors.New( "Error initializing tasks." ) 
  }

  err = stub.PutState( "toodles_tasks", bytes )

  return nil, nil
}

// *
// Invoke
// *

func ( t *SimpleChaincode ) Invoke( stub shim.ChaincodeStubInterface, function string, args []string ) ( []byte, error ) {
  if function == "init" { 
    return t.Init( stub, "init", args )  
  } else if function == "account_add" {
    return t.account_add( stub, args )
  } else if function == "account_edit" {
    return t.account_edit( stub, args )
  } else if function == "account_delete" {
    return t.account_delete( stub, args )
  } else if function == "task_add" {
    return t.task_add( stub, args )
  } else if function == "task_edit" {
    return t.task_edit( stub, args )
  } else if function == "task_delete" {
    return t.task_delete( stub, args )
  } else if function == "location_add" {
    return t.location_add( stub, args )
  } else if function == "location_edit" {
    return t.location_edit( stub, args )
  } else if function == "location_delete" {
    return t.location_delete( stub, args )
  } else if function == "reset_data" {
    return t.reset_data( stub, args )
  }

  return nil, errors.New( "Function with the name " + function + " does not exist." )  
}

// *
// Query
// *

func ( t *SimpleChaincode ) Query( stub shim.ChaincodeStubInterface, function string, args []string ) ( []byte, error ) {
  if function == "account_browse" {
    return t.account_browse( stub, args )
  } else if function == "account_read" {
    return t.account_read( stub, args )
  } else if function == "task_browse" {
    return t.task_browse( stub, args )
  } else if function == "task_read" {
    return t.task_read( stub, args )
  } else if function == "location_browse" {
    return t.location_browse( stub, args )
  } else if function == "location_read" {
    return t.location_read( stub, args )
  }

  return nil, errors.New( "Received unknown function invocation " + function )
}  

// *
// Account
// *

// Browse accounts
// Used to cross-assign tasks
func ( t *SimpleChaincode ) account_browse( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil { 
    return nil, errors.New( "Unable to get accounts." ) 
  }

  var accounts []Account
  var peers []Account

  // From JSON to data structure
  err = json.Unmarshal( bytes, &accounts )

  if args[0] == "all" {
    bytes, err = json.Marshal( accounts ) 
    return bytes, nil   
  }

  // Scrub passwords
  for a := 0; a < len( accounts ); a++ {  
    if accounts[a].Id != args[0] {
      accounts[a].Password = ""
      peers = append( peers, accounts[a] )
    }
  }    

  bytes, err = json.Marshal( peers )

  return bytes, nil
}

// Arguments: Name, Password
func ( t *SimpleChaincode ) account_read( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil { 
    return nil, errors.New( "Unable to get accounts." ) 
  }

  var accounts []Account

  // From JSON to data structure
  err = json.Unmarshal( bytes, &accounts )
  found := false

  // Look for match
  for _, account := range accounts {  
    // Match
    if account.Name == args[0] && account.Password == args[1] {
      // Sanitize
      account.Password = ""

      // JSON encode
      bytes, err = json.Marshal( account )
      found = true
      break
    }
  }

  // Nope
  if found != true {
    bytes, err = json.Marshal( nil )    
  }

  return bytes, nil
}

// Arguments: ID, First, Last, Name, Password
func ( t *SimpleChaincode ) account_edit( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil { 
    return nil, errors.New( "Unable to get accounts." ) 
  }

  var accounts []Account

  // From JSON to data structure
  err = json.Unmarshal( bytes, &accounts )

  // Look for match
  for a := 0; a < len( accounts ); a++ {  
    // Match
    if accounts[a].Id == args[0] {
      accounts[a].First = args[1]
      accounts[a].Last = args[2]      
      accounts[a].Name = args[3]
      accounts[a].Password = args[4]      
      break
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( accounts )
  err = stub.PutState( "toodles_accounts", bytes )

  return nil, nil
}

// Arguments: ID, First, Last, Name, Password
func ( t *SimpleChaincode ) account_add( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil { 
    return nil, errors.New( "Unable to get accounts." ) 
  }

  var account Account

  // Build JSON values
  id := "\"id\": \"" + args[0] + "\", "
  first := "\"first\": \"" + args[1] + "\", "
  last := "\"last\": \"" + args[2] + "\", "
  name := "\"name\": \"" + args[3] + "\", "
  password := "\"password\": \"" + args[4] + "\""

  // Make into a complete JSON string
  // Decode into a single account value
  content := "{" + id + first + last + name + password + "}"
  err = json.Unmarshal( []byte( content ), &account )

  var accounts []Account

  // Decode JSON into account array
  // Add latest account
  err = json.Unmarshal( bytes, &accounts )
  accounts = append( accounts, account )

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( accounts )
  err = stub.PutState( "toodles_accounts", bytes )

  return nil, nil
}

// Arguments: ID
func ( t *SimpleChaincode ) account_delete( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil { 
    return nil, errors.New( "Unable to get accounts." ) 
  }

  var accounts []Account

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &accounts )

  for a := 0; a < len( accounts ); a++ {  
    // Match
    if accounts[a].Id == args[0] {
      accounts = append( accounts[:a], accounts[a + 1:]... )
    }
  }  

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( accounts )
  err = stub.PutState( "toodles_accounts", bytes )

  return nil, nil
}

// *
// Task
// *

// Arguments: Account ID
func ( t *SimpleChaincode ) task_browse( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil { 
    return nil, errors.New( "Unable to get tasks." ) 
  }

  // Back door for viewing all data
  if( args[0] == "all" ) {
    return bytes, nil
  }

  var tasks []Task
  var items []Task

  // From JSON to data structure
  err = json.Unmarshal( bytes, &tasks )

  // Look for match
  for _, task := range tasks {  
    // Match
    if task.AccountId == args[0] {
      items = append( items, task )      
    }
  }

  // JSON encode
  bytes, err = json.Marshal( items )

  return bytes, nil
}

// Arguments: Id
func ( t *SimpleChaincode ) task_read( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil { 
    return nil, errors.New( "Unable to get tasks." ) 
  }

  var tasks []Task

  // From JSON to data structure
  err = json.Unmarshal( bytes, &tasks )
  found := false

  // Look for match
  for _, task := range tasks {  
    // Match
    if task.Id == args[0] {
      // JSON encode
      bytes, err = json.Marshal( task )
      found = true
      break
    }
  }

  // Nope
  if found != true {
    bytes, err = json.Marshal( nil )    
  }

  return bytes, nil
}

// Arguments: ID, Account ID, Due, Location ID, Duration, Energy, Tags, Notes, Complete, Name
func ( t *SimpleChaincode ) task_edit( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil { 
    return nil, errors.New( "Unable to get tasks." ) 
  }

  var tasks []Task

  // From JSON to data structure
  err = json.Unmarshal( bytes, &tasks )

  // Look for match
  for t := 0; t < len( tasks ); t++ {  
    // Match
    if tasks[t].Id == args[0] {
      // String arguments to integer values
      // TODO: Deal with errors
      due, _ := strconv.Atoi( args[2] )
      duration, _ := strconv.Atoi( args[4] )
      energy, _ := strconv.Atoi( args[5] )

      // No ternary operator
      complete := false
      if args[8] == "true" {
        complete = true
      }

      tasks[t].AccountId = args[1]
      tasks[t].Due = due
      tasks[t].LocationId = args[3]
      tasks[t].Duration = duration
      tasks[t].Energy =  energy
      tasks[t].Tags = args[6]
      tasks[t].Notes = args[7]
      tasks[t].Complete = complete
      tasks[t].Name = args[9]
      break
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( tasks )
  err = stub.PutState( "toodles_tasks", bytes )

  return nil, nil
}

// Arguments: ID, Account ID, Name, Due, Location, Duration, Energy, Created
func ( t *SimpleChaincode ) task_add( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil { 
    return nil, errors.New( "Unable to get tasks." ) 
  }

  var task Task

  // Build JSON values
  id := "\"id\": \"" + args[0] + "\", "
  account := "\"account\": \"" + args[1] + "\", "
  name := "\"name\": \"" + args[2] + "\", "
  due := "\"due\": " + args[3] + ", "
  location := "\"location\": \"" + args[4] + "\", "
  duration := "\"duration\": " + args[5] + ", "  
  energy := "\"energy\": " + args[6] + ", "
  tags := "\"tags\": \"\", "
  notes := "\"notes\": \"\", "
  complete := "\"complete\": false, "
  created := "\"created\": " + args[7]

  // Make into a complete JSON string
  // Decode into structure instance
  content := "{" + id + account + name + due + location + duration + energy + tags + notes + complete + created + "}"
  err = json.Unmarshal( []byte(content), &task )

  var tasks []Task

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &tasks )
  tasks = append( tasks, task )

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( tasks )
  err = stub.PutState( "toodles_tasks", bytes )

  return nil, nil
}

// Arguments: ID
func ( t *SimpleChaincode ) task_delete( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil { 
    return nil, errors.New( "Unable to get tasks." ) 
  }

  var tasks []Task

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &tasks )

  for t := 0; t < len( tasks ); t++ {  
    // Match
    if tasks[t].Id == args[0] {
      tasks = append( tasks[:t], tasks[t + 1:]... )
    }
  }  

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( tasks )
  err = stub.PutState( "toodles_tasks", bytes )

  return nil, nil
}

// *
// Location
// *

// Arguments: Account ID
func ( t *SimpleChaincode ) location_browse( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil { 
    return nil, errors.New( "Unable to get locations." ) 
  }

  // Back door for viewing all data
  if( args[0] == "all" ) {
    return bytes, nil
  }

  var locations []Location
  var items []Location

  // From JSON to data structure
  err = json.Unmarshal( bytes, &locations )

  // Look for match
  for _, location := range locations {  
    // Match
    if location.AccountId == args[0] {
      items = append( items, location )      
    }
  }

  // JSON encode
  bytes, err = json.Marshal( items )

  return bytes, nil
}

// Arguments: Id
func ( t *SimpleChaincode ) location_read( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil { 
    return nil, errors.New( "Unable to get locations." ) 
  }

  var locations []Location

  // From JSON to data structure
  err = json.Unmarshal( bytes, &locations )
  found := false

  // Look for match
  for _, location := range locations {  
    // Match
    if location.Id == args[0] {
      // JSON encode
      bytes, err = json.Marshal( location )
      found = true
      break
    }
  }

  // Nope
  if found != true {
    bytes, err = json.Marshal( nil )    
  }

  return bytes, nil
}

// Arguments: ID, Account ID, Name
func ( t *SimpleChaincode ) location_edit( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil { 
    return nil, errors.New( "Unable to get locations." ) 
  }

  var locations []Location

  // From JSON to data structure
  err = json.Unmarshal( bytes, &locations )

  // Look for match
  for g := 0; g < len( locations ); g++ {  
    // Match
    if locations[g].Id == args[0] {
      locations[g].AccountId = args[1]
      locations[g].Name = args[2]
      break
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( locations )
  err = stub.PutState( "toodles_locations", bytes )

  return nil, nil
}

// Arguments: ID, Account ID, Name
func ( t *SimpleChaincode ) location_add( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil { 
    return nil, errors.New( "Unable to get locations." ) 
  }

  var location Location

  // Build JSON values
  id := "\"id\": \"" + args[0] + "\", "
  account := "\"account\": \"" + args[1] + "\", "
  name := "\"name\": \"" + args[2] + "\""

  // Make into a complete JSON string
  // Decode into structure instance
  content := "{" + id + account + name + "}"
  err = json.Unmarshal( []byte(content), &location )

  var locations []Location

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &locations )
  locations = append( locations, location )

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( locations )
  err = stub.PutState( "toodles_locations", bytes )

  return nil, nil
}

// Arguments: ID
func ( t *SimpleChaincode ) location_delete( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil { 
    return nil, errors.New( "Unable to get locations." ) 
  }

  var locations []Location

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &locations )

  for g := 0; g < len( locations ); g++ {  
    // Match
    if locations[g].Id == args[0] {
      locations = append( locations[:g], locations[g + 1:]... )
    }
  }  

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( locations )
  err = stub.PutState( "toodles_locations", bytes )

  return nil, nil
}

// *
// Reset data
// *

func ( t *SimpleChaincode ) reset_data( stub shim.ChaincodeStubInterface, args []string ) ( []byte, error ) {
  stub.PutState( "toodles_accounts", []byte( args[0] ) )
  stub.PutState( "toodles_locations", []byte( args[1] ) )  
  stub.PutState( "toodles_tasks", []byte( args[2] ) )    

  return nil, nil  
}


// Main
func main() {
  err := shim.Start( new( SimpleChaincode ) )

  if err != nil { 
    fmt.Printf( "Error starting chaincode: %s", err ) 
  }
}
