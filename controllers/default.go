package controllers

type MainController struct {
	baseController
}

func (this *MainController) Get() {
	if this.IsLogin == true {
		this.TplName = "index.html"
		return
	}
	return
}
