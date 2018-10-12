package importer

type rasterFile struct {
	systemName  string // What is the system going to refer to this as?
	urlSource   string // Where did this come from?
	shpFilePath string // What is the shpFile path to work with
}
