package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/crhym3/go-endpoints/endpoints"

	"appengine"
	"appengine/urlfetch"
)

const (
	API_KEY = "<your api key here>"
)

func init() {
	registerCloudEndpoints()
}

func registerCloudEndpoints() {
	// Create instance of our service
	service := &MessengerService{}

	// Register our service as an endpoint service
	api, err := endpoints.RegisterService(service, "messenger", "v1", "Messenger API", true)
	if err != nil {
		panic(err.Error())
	}

	// Register Echo method
	info := api.MethodByName("Echo").Info()
	info.Name, info.HttpMethod, info.Path, info.Desc =
		"Echo", "POST", "Echo", "Echoes your message."

	// Register Send method
	info = api.MethodByName("Send").Info()
	info.Name, info.HttpMethod, info.Path, info.Desc =
		"Send", "POST", "Send", "Sends your message to the specified device."

	// Start handling requests
	endpoints.HandleHttp()
}

// Empty struct to attach our endpoint methods to.
type MessengerService struct {
}

// Message just wraps a string in a struct.  We could add other fields here if we wanted.
type Message struct {
	Content string
	To      string
}

// Echo repeats the message you send back to you.
func (ms *MessengerService) Echo(r *http.Request, m *Message, resp *Message) error {
	resp.Content = "Recieved message: " + m.Content
	return nil
}

// GCMMsg is an message which can be encoded to JSON the way the GCM server likes
type GCMMsg struct {
	Data            map[string]string `json:"data"`
	RegistrationIDs []string          `json:"registration_ids"`
}

// Creates a GCMMsg from a Message
func createGCMMsg(msg *Message) *GCMMsg {
	title := "WUSTL HACK GCM MESSAGE"

	// Prepare GCM message
	gcmm := &GCMMsg{}
	gcmm.Data = map[string]string{"title": title, "message": msg.Content}
	gcmm.RegistrationIDs = []string{msg.To}

	return gcmm
}

// Encodes a Message as JSON compatible with GCM
func encodeMsg(msg *Message) io.Reader {
	gcmm := createGCMMsg(msg)

	// Encode our message as JSON
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(gcmm)

	// Create an io.Reader to return
	return bytes.NewReader(buf.Bytes())
}

func createRequest(msg *Message) (*http.Request, error) {
	// Encode Message into JSON compatible with GCM
	body := encodeMsg(msg)

	// Create our request and handle errors
	req, err := http.NewRequest("POST",
		"https://android.googleapis.com/gcm/send",
		body)

	if err != nil {
		return nil, err
	}

	// Add Auth headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "key="+API_KEY)

	return req, nil
}

// Send crafts a msg to send to GCM based on the Message passed in.  It returns the
// GCM message which was sent to the server and logs the server's response.
func (ms *MessengerService) Send(r *http.Request, m *Message, res *Message) error {
	// Create an http.Request from our input message.
	req, err := createRequest(m)
	if err != nil {
	   return err
	}

	// Send the request and handle errors
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Return the response from the GCM server as a Message
	content, _ := ioutil.ReadAll(resp.Body)
	res.Content = fmt.Sprintf("Response from GCM server: %#v", string(content))
	return nil
}
