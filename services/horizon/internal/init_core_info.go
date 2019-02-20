package horizon

func initStellarCoreInfo(app *App) {
	app.UpdateStellarCoreInfo()
}

func init() {
	appInit.Add("stellarCoreInfo", initStellarCoreInfo, "app-context", "log")
}
