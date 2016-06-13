package alidayu_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/northbright/alidayu"
)

// Example of Client.Post()
func ExampleClient_Post() {
	// Create a new client.
	c := &alidayu.Client{AppKey: "", AppSecret: "", UseHTTPS: false}

	// ---------------------------------------
	// Send Verification Code in SMS.
	// ---------------------------------------

	// Set Parameters.
	params := map[string]string{}
	// It'll use default common parameters if you don't set them.
	// params["format"] = "json"
	// params["v"] = "2.0"
	// params["sign_method"] = "md5"
	params["method"] = "alibaba.aliqin.fc.sms.num.send"           // Set method to send SMS.
	params["sms_type"] = "normal"                                 // Set SMS type.
	params["sms_free_sign_name"] = ""                             // Set SMS signature.
	params["sms_param"] = `{"code":"123456", "product":"My App"}` // Set variable for SMS template.
	params["sms_template_code"] = ""                              // Set SMS template code.
	params["rec_num"] = ""                                        // Set phone number to send SMS.

	// Call Post() to post the request.
	resp, err := c.Post(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Post() error: %v\n", err)
		return
	}

	// Read HTTP Response.
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadAll() error:%v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Post() successfully\n%v\n", string(data))

	// ------------------------------------------
	// Send Verification Code in Single Call.
	// ------------------------------------------
	// Set Parameters.
	params = map[string]string{}
	params["method"] = "alibaba.aliqin.fc.tts.num.singlecall"     // Set method to make single call.
	params["tts_param"] = `{"code":"123456", "product":"My App"}` // Set variable for TTS template.
	params["called_num"] = ""                                     // Set phone number to make single call.
	params["called_show_num"] = ""                                // Set show number.
	params["tts_code"] = ""                                       // Set TTS code.

	// Call Post() to post the request.
	resp2, err := c.Post(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Post() error: %v\n", err)
		return
	}

	// Read HTTP Response.
	defer resp2.Body.Close()
	data2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadAll() error:%v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Post() successfully\n%v\n", string(data2))

	// Output:
}

// Example of Client.Exec()
func ExampleClient_Exec() {
	// Create a new client.
	c := &alidayu.Client{AppKey: "", AppSecret: "", UseHTTPS: false}

	// ---------------------------------------
	// Send Verification Code in SMS.
	// ---------------------------------------

	// Set Parameters.
	params := map[string]string{}
	// It'll use default common parameters if you don't set them.
	// params["format"] = "json"
	// params["v"] = "2.0"
	// params["sign_method"] = "md5"
	params["format"] = "xml"
	params["method"] = "alibaba.aliqin.fc.sms.num.send"           // Set method to send SMS.
	params["extend"] = "123456"                                   // Set callback parameter.
	params["sms_type"] = "normal"                                 // Set SMS type.
	params["sms_free_sign_name"] = ""                             // Set SMS signature.
	params["sms_param"] = `{"code":"123456", "product":"My App"}` // Set variable for SMS template.
	params["sms_template_code"] = ""                              // Set SMS template code.
	params["rec_num"] = ""                                        // Set phone number to send SMS.

	// Call Exec() to post the request.
	success, result, err := c.Exec(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Exec() error: %v\nsuccess: %v\nresult: %v\n", err, success, result)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Exec() successfully\nsuccess: %v\nresult: %s\n", success, result)

	// ------------------------------------------
	// Send Verification Code in Single Call.
	// ------------------------------------------
	// Set Parameters.
	params = map[string]string{}
	params["method"] = "alibaba.aliqin.fc.tts.num.singlecall"     // Set method to make single call.
	params["tts_param"] = `{"code":"123456", "product":"My App"}` // Set variable for TTS template.
	params["called_num"] = ""                                     // Set phone number to make single call.
	params["called_show_num"] = ""                                // Set show number.
	params["tts_code"] = ""                                       // Set TTS code.

	// Call Exec() to post the request.
	success, result, err = c.Exec(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Exec() error: %v\nsuccess: %v\nresult: %v\n", err, success, result)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Exec() successfully\nsuccess: %v\nresult: %s\n", success, result)

	// Output:
}
