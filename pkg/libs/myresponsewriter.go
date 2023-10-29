package libs

type ResponseWriterStruct struct {
	MyHeader      map[string]string
	MyWriter      ([]byte)
	MyWriteHeader (int)
}
