package controllers

import (
	"gatssh/models"
	"github.com/tealeg/xlsx"
	"fmt"
	"strconv"
	"bytes"
)

type QueryGatSshTaskResults struct {
	baseController
}

type taskResults struct {
	TaskDetails  []models.TaskDetail `json:"taskDetails"`
	CurrentCount int                 `json:"currentCount"`
	LastID       int                 `json:"lastId"`
}

func (this *QueryGatSshTaskResults) QueryTaskResults() {
	if this.IsLogin != true {
		this.ServeJSON(40000, "Please login...")
		return
	}

	taskId := this.GetString("taskId")

	startID, err := this.GetInt("startId")
	if err != nil {
		startID = 0
	}

	taskDetails, currentCount, lastID, err := models.QueryTaskDetails(taskId, startID)
	if err != nil {
		this.ServeJSON(40000, err)
		return
	}

	tr := &taskResults{
		TaskDetails:  taskDetails,
		CurrentCount: currentCount,
		LastID:       lastID,
	}

	this.ServeJSON(20000, tr)
	return
}

func (this *QueryGatSshTaskResults) DownloadExcel() {

	if this.IsLogin != true {
		this.ServeJSON(40000, "Please login...")
		return
	}

	taskId := this.GetString("taskId")

	taskDetails, _, _, err := models.QueryTaskDetails(taskId, 0)

	if err != nil {
		this.ServeJSON(40000, err)
		return
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "IP"
	cell = row.AddCell()
	cell.Value = "Code"
	cell = row.AddCell()
	cell.Value = "Stdout"
	cell = row.AddCell()
	cell.Value = "Stderr"

	for _, td := range taskDetails {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = td.Ip
		cell = row.AddCell()
		cell.Value = strconv.Itoa(td.ResultCode)
		cell = row.AddCell()
		cell.Value = td.ResultContent
		cell = row.AddCell()
		cell.Value = td.ResultErr
	}

	buf := new(bytes.Buffer)

	file.Write(buf)

	this.Ctx.Output.Header("Content-Disposition", "attachment;filename=taskDetails.xlsx")

	this.Ctx.Output.Body(buf.Bytes())

	return
}
