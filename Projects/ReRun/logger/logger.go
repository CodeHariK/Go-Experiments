package logger

type StdLogSave struct {
	savedOutput map[string][]string
	Group       string
	fn          func(p []byte) (n int, err error)
}

func CreateStdOutSave(
	savedOutput map[string][]string,
	fn func(p []byte) (n int, err error),
) *StdLogSave {
	return &StdLogSave{
		savedOutput: savedOutput,
		fn:          fn,
	}
}

func (so *StdLogSave) Write(p []byte) (n int, err error) {
	so.savedOutput[so.Group] = append(so.savedOutput[so.Group], string(p))

	// jsonData, _ := json.MarshalIndent(so.savedOutput, "", "  ")
	// fmt.Println(string(jsonData))

	return so.fn(p)
}
