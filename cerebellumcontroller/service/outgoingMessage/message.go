package outgoingMessage

import (
	"errors"
	"iai/cerebellumController/dto/common"
	"iai/cerebellumController/dto/message"
	"iai/cerebellumController/util"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func SendMessage(host string, msg *message.Message, src string, dst string, taskid string) error {
	log.Println(host, msg.File, src, dst, taskid, util.ReGenFileNameOutgoing(msg.File, taskid))
	c := fiber.AcquireClient()
	a := c.Post(host + "/api/v1/message/incoming")
	args := fiber.AcquireArgs()
	args.Set("src", src)
	args.Set("dst", dst)
	args.Set("taskid", taskid)
	args.Set("metadata", msg.MetaData)
	newName := util.ReGenFileNameOutgoing(msg.File, taskid)
	a.SendFile(newName, "file").MultipartForm(args)
	a.Parse()
	fiber.ReleaseArgs(args)
	os.Remove(newName)
	var resp common.Response
	code, _, errs := a.Struct(&resp)
	log.Println(resp)
	if code != 200 {
		return errors.New("recv non 200 code")
	}
	if errs != nil {
		error_detail := ""
		for _, err := range errs {
			error_detail += err.Error() + "\n"
		}
		return errors.New(error_detail)
	}
	if resp.Code != 0 {
		return errors.New(resp.Msg)
	}
	return nil

}
