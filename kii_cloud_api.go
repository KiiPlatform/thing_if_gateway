// Provides APIs to access to Kii Cloud and
// Thing Interaction Framework (thing-if).
package kii

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Represents Application in Kii Cloud.
type App struct {
	AppID       string
	AppKey      string
	AppLocation string
}

// Obtain Host name of the Application endpoint.
func (ka *App) HostName() string {
	lowerLoc := strings.ToLower(ka.AppLocation)
	switch lowerLoc {
	case "jp":
		return "api-jp.kii.com"
	case "us":
		return "api.kii.com"
	case "cn":
		return "api-cn3.kii.com"
	case "sg":
		return "api-sg.kii.com"
	default:
		return lowerLoc
	}
}

// Obtain thing-if endpoint base url.
func (ka *App) ThingIFBaseUrl() string {
	return fmt.Sprintf("https://%s/thing-if/apps/%s", ka.HostName(), ka.AppID)
}

// Obtain Kii Cloud endpoint base url.
func (ka *App) KiiCloudBaseUrl() string {
	return fmt.Sprintf("https://%s/api/apps/%s", ka.HostName(), ka.AppID)
}

// Layout position of the Thing
type LayoutPosition int

const (
	ENDNODE LayoutPosition = iota
	STANDALONE
	GATEWAY
)

// Obtain Layout postion of the Thing in string.
func (lp LayoutPosition) String() string {
	switch lp {
	case ENDNODE:
		return "END_NODE"
	case STANDALONE:
		return "STANDALONE"
	case GATEWAY:
		return "GATEWAY"
	default:
		log.Fatal("never reache here")
		return "invalid layout"
	}
}

func executeRequest(request http.Request) (respBody []byte, error error) {

	client := &http.Client{}
	resp, err := client.Do(&request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("body: " + string(bodyStr))

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return bodyStr, nil
	} else {
		err = errors.New(string(bodyStr))
		return nil, err
	}
}

// Struct for requesting Gateway Onboard.
type OnboardGatewayRequest struct {
	VendorThingID   string                 `json:"vendorThingID"`
	ThingPassword   string                 `json:"thingPassword"`
	ThingType       string                 `json:"thingType"`
	LayoutPosition  string                 `json:"layoutPosition"`
	ThingProperties map[string]interface{} `json:"thingProperties"`
}

// Struct for receiving response of Gateway Onboard.
type OnboardResponse struct {
	ThingID      string       `json:"thingID"`
	AccessToken  string       `json:"accessToken"`
	MqttEndpoint MqttEndpoint `json:"mqttEndpoint"`
}

// Struct represents MQTT endpoint.
type MqttEndpoint struct {
	InstallationID string `json:"installationID"`
	Host           string `json:"host"`
	MqttTopic      string `json:"mqttTopic"`
	Username       string `json:"userName"`
	Password       string `json:"password"`
	PortSSL        int    `json:"portSSL"`
	PortTCP        int    `json:"portTCP"`
}

// Struct represents API author.
// Can be Gateway, EndNode or KiiUser, depending on the token.
type APIAuthor struct {
	Token string
	App   App
}

// Struct for requesting end node token
type EndNodeTokenRequest struct {
	ExpiresIn string `json:"expires_in,omitempty"`
}

// Struct for receiving response of end node token
type EndNodeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	ThingID      string `json:"id"`
	RefreshToken string `json:"refresh_token"`
}

// Struct of predefined fileds for requesting Thing Registration.
type RegisterThingRequest struct {
	VendorThingID   string `json:"_vendorThingID"`
	ThingPassword   string `json:"_password"`
	ThingType       string `json:"_thingType,omitempty"`
	LayoutPosition  string `json:"_layoutPosition,omitempty"`
	Vendor          string `json:"_vendor,omitempty"`
	FirmwareVersion string `json:"_firmwareVersion,omitempty"`
	Lot             string `json:"_lot,omitempty"`
	StringField1    string `json:"_stringField1,omitempty"`
	StringField2    string `json:"_stringField2,omitempty"`
	StringField3    string `json:"_stringField3,omitempty"`
	StringField4    string `json:"_stringField4,omitempty"`
	StringField5    string `json:"_stringField5,omitempty"`
	NumberField1    int64  `json:"_numberField1,omitempty"`
	NumberField2    int64  `json:"_numberField2,omitempty"`
	NumberField3    int64  `json:"_numberField3,omitempty"`
	NumberField4    int64  `json:"_numberField4,omitempty"`
	NumberField5    int64  `json:"_numberField5,omitempty"`
}

// Struct for receiving response of end node token
type RegisterThingResponse struct {
	ThingID        string `json:"_thingID"`
	VendorThingID  string `json:"_vendorThingID"`
	ThingType      string `json:"_thingType"`
	LayoutPosition string `json:"_layoutPosition"`
	Created        int    `json:"_created"`
	Disabled       bool   `json:"_disabled"`
}

// Struct for request registration of KiiUser.
// At least one of LoginName, EmailAddress or PhoneNumber must be provided.
type KiiUserRegisterRequest struct {
	LoginName           string `json:"loginName,omitempty"`
	DisplayName         string `json:"displayName,omitempty"`
	Country             string `json:"country,omitempty"`
	Locale              string `json:"locale,omitempty"`
	EmailAddress        string `json:"emailAddress,omitempty"`
	PhoneNumber         string `json:"phoneNumber,omitempty"`
	PhoneNumberVerified bool   `json:"phoneNumberVerified,omitempty"`
	Password            string `json:"password"`
}

// Struct for receiving registration of KiiUser.
type KiiUserRegisterResponse struct {
	UserID              string `json:"userID"`
	LoginName           string `json:"loginName"`
	DisplayName         string `json:"displayName"`
	Country             string `json:"country"`
	Locale              string `json:"locale"`
	EmailAddress        string `json:"emailAddress"`
	PhoneNumber         string `json:"phoneNumber"`
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	HasPassword         bool   `json:"_hasPassword"`
}

// Struct for requesting login of KiiUser
type KiiUserLoginRequest struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
}

// Struct for receiving response of login
type KiiUserLoginResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

// Struct for posting command
// Issuer can be group or user.
// If user, must be "user:<user-id>".
type PostCommandRequest struct {
	Issuer           string                   `json:"issuer"`
	Actions          []map[string]interface{} `json:"actions"`
	Schema           string                   `json:"schema"`
	SchemaVersion    int                      `json:"schemaVersion"`
	FiredByTriggerID string                   `json:"firedByTriggerID,omitempty"`
	Titlle           string                   `json:"title,omitempty"`
	Description      string                   `json:"description,omitempty"`
	Metadata         map[string]interface{}   `json:"metadata,omitempty"`
}

// Struct for receiving response of posting command
type PostCommandResponse struct {
	CommandID string `json:"commandID"`
}

// Struct for requesting Onboard by Thing Owner.
type OnboardByOwnerRequest struct {
	ThingID        string `json:"thingID"`
	ThingPassword  string `json:"thingPassword"`
	Owner          string `json:"owner"`
	LayoutPosition string `json:"layoutPosition,omitempty"` // pattern: GATEWAY|STANDALONE|ENDNODE, STANDALONE by default
}

// Struct for updating command results
type UpdateCommandResultsRequest struct {
	ActionResults []map[string]interface{} `json:"actionResults"`
}

// Login as Anonymous user.
// When there's no error, APIAuthor is returned.
func AnonymousLogin(app App) (*APIAuthor, error) {
	type AnonymousLoginRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}
	type AnonymousLoginResponse struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	reqObj := AnonymousLoginRequest{
		ClientID:     app.AppID,
		ClientSecret: app.AppKey,
		GrantType:    "client_credentials",
	}
	reqJson, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/oauth2/token", app.KiiCloudBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	}

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("body: " + string(bodyStr))

	var respObj AnonymousLoginResponse
	err = json.Unmarshal(bodyStr, &respObj)
	if err != nil {
		return nil, err
	}
	au := APIAuthor{
		Token: respObj.AccessToken,
		App:   app,
	}
	return &au, nil
}

// Let Gateway onboard to the cloud.
// When there's no error, OnboardResponse is returned.
func (au *APIAuthor) OnboardGateway(request OnboardGatewayRequest) (*OnboardResponse, error) {
	var ret OnboardResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/onboardings", au.App.ThingIFBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/vnd.kii.onboardingWithVendorThingIDByThing+json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Request access token of end node of gateway.
// Notes the APIAuthor should be a Gateway.
// When there's no error, EndNodeTokenResponse is returned.
func (au APIAuthor) GenerateEndNodeToken(gatewayID string, endnodeID string, request EndNodeTokenRequest) (*EndNodeTokenResponse, error) {
	var ret EndNodeTokenResponse
	url := fmt.Sprintf("%s/things/%s/end-nodes/%s/token", au.App.KiiCloudBaseUrl(), gatewayID, endnodeID)

	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Add an end node thing to gateway
// Notes that the APIAuthor should be a Gateway
func (au APIAuthor) AddEndNode(gatewayID string, endnodeID string) error {
	url := fmt.Sprintf("%s/things/%s/end-nodes/%s", au.App.KiiCloudBaseUrl(), gatewayID, endnodeID)

	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)
	if err != nil {
		return err
	}

	_, err1 := executeRequest(*req)
	return err1
}

// Register Thing.
// The request must consist of the predefined fields(see RegisterThingRequest).
// If you want to add the custom fileds, you can simply make RegisterThingRequest as anonymous field of your defined request struct, like:
//  type MyRegisterThingRequest struct {
//    RegisterThingRequest
//    MyField1             string
//  }
// Where there is no error, RegisterThingResponse is returned
func (au APIAuthor) RegisterThing(request interface{}) (*RegisterThingResponse, error) {
	var ret RegisterThingResponse

	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/things", au.App.KiiCloudBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/vnd.kii.ThingRegistrationRequest+json")
	req.Header.Set("X-Kii-AppID", au.App.AppID)
	req.Header.Set("X-Kii-AppKey", au.App.AppKey)

	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Update Thing state.
// Notes that the APIAuthor should be already initialized as a Gateway or EndNode
func (au APIAuthor) UpdateState(thingID string, request interface{}) error {

	reqJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/targets/thing:%s/states", au.App.ThingIFBaseUrl(), thingID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	_, err1 := executeRequest(*req)
	return err1
}

// Login as KiiUser.
// If there is no error, KiiUserLoginResponse is returned.
// Notes that after login successfully, api doesn't update token of APIAuthor,
// you should update by yourself with the token in response.
func (au *APIAuthor) LoginAsKiiUser(request KiiUserLoginRequest) (*KiiUserLoginResponse, error) {
	var ret KiiUserLoginResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s/api/oauth2/token", au.App.HostName())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Kii-AppID", au.App.AppID)
	req.Header.Set("X-Kii-AppKey", au.App.AppKey)
	log.Printf("login request body:%s", string(reqJson))
	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}

}

// Register KiiUser
// If there is no error, KiiUserRegisterResponse is returned.
func (au *APIAuthor) RegisterKiiUser(request KiiUserRegisterRequest) (*KiiUserRegisterResponse, error) {
	var ret KiiUserRegisterResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/users", au.App.KiiCloudBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Kii-AppID", au.App.AppID)
	req.Header.Set("X-Kii-AppKey", au.App.AppKey)
	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}

}

// Post command to Thing.
// Notes that it requires Thing already onboard.
// If there is no error, PostCommandRequest is returned.
func (au APIAuthor) PostCommand(thingID string, request PostCommandRequest) (*PostCommandResponse, error) {
	var ret PostCommandResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/targets/THING:%s/commands", au.App.ThingIFBaseUrl(), thingID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)
	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Update command results
func (au APIAuthor) UpdateCommandResults(thingID string, commandID string, request UpdateCommandResultsRequest) error {
	reqJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/targets/thing:%s/commands/%s/action-results", au.App.ThingIFBaseUrl(), thingID, commandID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	_, err = executeRequest(*req)
	return err
}

func (au *APIAuthor) OnboardThingByOwner(request OnboardByOwnerRequest) (*OnboardResponse, error) {
	var ret OnboardResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/onboardings", au.App.ThingIFBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/vnd.kii.OnboardingWithThingIDByOwner+json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}
