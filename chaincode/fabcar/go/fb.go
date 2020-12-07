import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/cd1/utils-golang"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)
type area struct {
	Doctype  string
	Name     string
	Division string
	District string
	Thana    string
	Key      string
}
type candidate struct {
	Doctype   string
	Election  string
	Name      string
	AreaName  string
	Totalvote int
	Sign      string
	Key       string
}

type election struct {
	Doctype string
	Name    string
	Key     string
}

type commission struct {
	Doctype  string
	Name     string
	Email    string
	Password string
	Key      string
}

type vote struct {
	Doctype string
	UserKey string
	ElectionName string
	Key string
}

type user struct {
	Doctype           string
	Nid               int64
	Mobile            string
	Birthdate         string
	Email             string
	PresentDivision   string
	PresentDistrict   string
	PresentThana      string
	PermanentDivision string
	PermanentDistrict string
	PermanentThana    string
	Password          string
	Key               string
	Token             string
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	s.initLedger(APIstub); //initLedger ta call to dibi -_-
	return shim.Success(nil)
} ///
func (s *SmartContract) addArea(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	addArea := area{}
	addArea.Doctype = "area"
	addArea.Name = args[0]
	addArea.Division = args[1]
	addArea.District = args[2]
	addArea.Thana = args[3]
	key := utils.RandomString()
	addArea.Key = key

	jsonUser, _ := json.Marshal(addArea)

	_ = APIstub.PutState(key, jsonUser)
	return shim.Success(nil)
}

/////  Init comission

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	commission := []commission{
		commission{Doctype: "commission", Name: "Romana", Email: "rme@gmail.com", Password: "1234567", Key: utils.RandomString()},
		commission{Doctype: "commission", Name: "Mahjabin", Email: "nullcharacter97@gmail.com", Password: "1234567", Key: utils.RandomString()},
		commission{Doctype: "commission", Name: "Alfa", Email: "alfa@gmail.com", Password: "1234567", Key: utils.RandomString()},
	}

	i := 0
	for i < len(commission) {
		fmt.Println("i is ", i)
		CommissionBytes, _ := json.Marshal(commission[i])
		APIstub.PutState(commission[i].Key, CommissionBytes)
		fmt.Println("Added", commission[i].Key)
		i = i + 1
	}

	s.addArea(APIstub, []string{"Sylhet-1", "Sylhet", "Sylhet", "Kutwali"})
	s.addArea(APIstub, []string{"Sylhet-2", "Sylhet", "Sunamgonj", "Osmani Nagar"})
	s.addArea(APIstub, []string{"Bahubol-2", "Sylhet", "Habiganj", "Nabiganj"})

	return shim.Success(nil)

}

/////
/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "register" {
		return s.register(APIstub, args)
	} else if function == "login" {
		return s.login(APIstub, args)
	} else if function == "getData" {
		return s.getDataFromArgs(APIstub, args)
	} else if function == "setData" {
		return s.setData(APIstub, args)
	} else if function == "addElection" {
		return s.addElection(APIstub, args)
	} else if function == "addCandidate" {
		return s.addCandidate(APIstub, args)
	} else if function == "getCandidates" {
		return s.getCandidates(APIstub, args)
	} else if function == "admin" {
		return s.admin(APIstub, args)
	} else if function == "addCandidate" {
		return s.addCandidate(APIstub, args)
	} else if function == "addElecion" {
		return s.addElection(APIstub, args)
	} else if function == "addArea" {
		return s.addArea(APIstub, args)
	} else if function == "addVote" {
		return s.addVote(APIstub, args)
	} else if function == "latestElection" {
		return s.latestElection(APIstub, args)
	} else if function == "getHistory" {
		return s.getHistory(APIstub, args)
	}else if function == "getCandidateList" {
		return s.getCandidateList(APIstub, args)
	}else if function == "hasAlreadyVoted" {
		return s.hasAlreadyVoted(APIstub, args)
	}


	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) register(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments, required 11, given " + strconv.Itoa(len(args)))
	}

	mUser := user{}

	mUser.Doctype = "user"
	mUser.Nid, _ = strconv.ParseInt(args[0], 10, 64)
	mUser.Mobile = args[1]
	mUser.Birthdate = args[2]
	mUser.Email = args[3]
	mUser.PresentDivision = args[4]
	mUser.PresentDistrict = args[5]
	mUser.PresentThana = args[6]
	mUser.PermanentDivision = args[7]
	mUser.PermanentDistrict = args[8]
	mUser.PermanentThana = args[9]
	mUser.Password = args[10]
	token := utils.RandomString()
	mUser.Token = token

	h := sha256.New()
	h.Write([]byte(mUser.Password))
	mUser.Password = fmt.Sprintf("%x", h.Sum(nil))

	userKey := utils.RandomString()
	mUser.Key = userKey

	jsonUser, _ := json.Marshal(mUser)

	_ = APIstub.PutState(userKey, jsonUser)
	return shim.Success(nil)

}

func (s *SmartContract) addVote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	candidateKey := args[0]
	userKey := args[1]
	currentCandidateData, _ := APIstub.GetState(candidateKey)
	var currentCandidate candidate
	_ = json.Unmarshal(currentCandidateData, &currentCandidate)
	currentCandidate.Totalvote = currentCandidate.Totalvote + 1
	jsonCurrentCandidate, _ := json.Marshal(currentCandidate)
	_ = APIstub.PutState(candidateKey, jsonCurrentCandidate)

	elec := s.getLatestElection(APIstub,args)
	myVote := vote{"vote", userKey, elec.Name, utils.RandomString()}
	myVoteData, _ := json.Marshal(myVote)
	_ = APIstub.PutState(myVote.Key, myVoteData)

	return shim.Success(nil)
}

func (s *SmartContract) addElection(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, required 1, given " + strconv.Itoa(len(args)))
	}
	addElection := election{}
	addElection.Doctype = "election"
	addElection.Name = args[0]
	key := utils.RandomString()
	addElection.Key = key

	jsonElection, _ := json.Marshal(addElection)

	_ = APIstub.PutState(key, jsonElection)
	return shim.Success(nil)
}

func (s *SmartContract) getHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	voteQuery := newCouchQueryBuilder().addSelector("Doctype", "vote").getQueryString()
	historyData, _ :=getJSONQueryResultForQueryString(APIstub, voteQuery)
	return shim.Success(historyData)
}

func (s *SmartContract) getLatestElection(APIstub shim.ChaincodeStubInterface, args []string) (election) {
	electionQuery := newCouchQueryBuilder().addSelector("Doctype", "election").getQueryString()
	electionIterator, _ := APIstub.GetQueryResult(electionQuery)
	var latestElection election

	for electionIterator.HasNext() {
		electionData, _ := electionIterator.Next()
		_ = json.Unmarshal(electionData.Value, &latestElection)
	}
	return latestElection
}


func (s *SmartContract) latestElection(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	latestElection := s.getLatestElection(APIstub, args)
	jsonLatestElection, _ := json.Marshal(latestElection)
	return shim.Success(jsonLatestElection)
}

func (s *SmartContract) addCandidate(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments, required 4, given " + strconv.Itoa(len(args)))
	}
	addCandidate := candidate{}
	addCandidate.Election = args[0]
	addCandidate.Name = args[1]
	addCandidate.AreaName = args[2]
	addCandidate.Sign = args[3]
	addCandidate.Totalvote = 0
	Key := utils.RandomString()
	addCandidate.Key = Key

	jsonCandidate, _ := json.Marshal(addCandidate)

	_ = APIstub.PutState(Key, jsonCandidate)
	return shim.Success(nil)
}

func (s *SmartContract) login(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("you need 3 arguments, but you have given " + strconv.Itoa(len(args)))
	}

	Nid, _ := strconv.ParseInt(args[0], 10, 64)
	Birthdate := args[1]
	Password := args[2]

	h := sha256.New()
	h.Write([]byte(Password))
	Password = fmt.Sprintf("%x", h.Sum(nil))

	queryString := newCouchQueryBuilder().addSelector("Nid", Nid).addSelector("Birthdate", Birthdate).addSelector("Password", Password).getQueryString()

	userData, _ := firstQueryValueForQueryString(APIstub, queryString)
	return shim.Success(userData)
}
func (s *SmartContract) admin(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("you need 2 arguments, but you have given " + strconv.Itoa(len(args)))
	}

	Email := args[0]
	Password := args[1]

	adminQuery := newCouchQueryBuilder().addSelector("Email", Email).addSelector("Password", Password).getQueryString()

	adminData, _ := firstQueryValueForQueryString(APIstub, adminQuery)
	return shim.Success(adminData)
}

func (s *SmartContract) getData(APIstub shim.ChaincodeStubInterface, args ...string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	data, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error("There was an error")
	}

	return shim.Success(data)
}

//////

func (s *SmartContract) hasAlreadyVoted(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	userKey:=args[0]
	currentUserData,_:=APIstub.GetState(userKey)
	var currentUser user
	_=json.Unmarshal(currentUserData,&currentUser)
	election:=s.getLatestElection(APIstub,args)
	currentStatus:=newCouchQueryBuilder().addSelector("Doctype","vote").addSelector("UserKey",userKey).addSelector("ElectionName",election.Name).getQueryString()
	currentStatusData,_:=getJSONQueryResultForQueryString(APIstub,currentStatus)
	return shim.Success(currentStatusData)
}
////


func (s *SmartContract) getCandidateList(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	election := s.getLatestElection(APIstub, args)
	currentElection:=election.Name
	candidateQuery:=newCouchQueryBuilder().addSelector("Doctype","candidate").addSelector("Election",currentElection).getQueryString()
	candidateList, _ := getJSONQueryResultForQueryString(APIstub, candidateQuery)
	return shim.Success(candidateList)
}
//
func (s *SmartContract) getCandidates(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	Key := args[0]

	currentUserData, _ := APIstub.GetState(Key)
	var currentUser user
	_ = json.Unmarshal(currentUserData, &currentUser)

	presentThana := currentUser.PresentThana

	areaQuery := newCouchQueryBuilder().addSelector("Doctype", "area").addSelector("Thana", presentThana).getQueryString()

	currentAreaData, _ := firstQueryValueForQueryString(APIstub, areaQuery)
	var currentArea area
	_ = json.Unmarshal(currentAreaData, &currentArea)
	elec := s.getLatestElection(APIstub, args)
	candidateQuery := newCouchQueryBuilder().addSelector("Doctype", "candidate").addSelector("AreaName", currentArea.Name).addSelector("Election", elec.Name).getQueryString()

	candidatesData, _ := getJSONQueryResultForQueryString(APIstub, candidateQuery)
	fmt.Println(candidatesData)

	//print the output
	fmt.Println(string(candidatesData))

	return shim.Success(candidatesData)
}

func (s *SmartContract) getCurrentElection(APIstub shim.ChaincodeStubInterface, args []string) string {
	key := args[0]

	currentElectionData, _ := APIstub.GetState(key)
	var currentElection election
	_ = json.Unmarshal(currentElectionData, &currentElection)

	currentElectionName := currentElection.Name
	return currentElectionName

}

////

func (s *SmartContract) getDataFromArgs(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	data, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error("There was an error")
	}

	return shim.Success(data)
}

func (s *SmartContract) setData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	key := args[0]
	val := args[1]

	err := APIstub.PutState(key, []byte(val))
	if err != nil {
		return shim.Error("There was an error")
	}

	str := "operation successful"

	return shim.Success([]byte(str))
}

func MockInvoke(stub *shim.MockStub, function string, args []string) sc.Response {
	input := args
	output := make([][]byte, len(input)+1)
	output[0] = []byte(function)
	for i, v := range input {
		output[i+1] = []byte(v)
	}

	fmt.Println("final arguments: ", output) // [[102 111 111] [98 97 114]]

	return stub.MockInvoke("1", output)
}

