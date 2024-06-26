package logger

import "sync"

type StdLogSave struct {
	mu          sync.RWMutex
	savedOutput map[string][]string
	Group       string
	fn          func(string, func(string)) (n int, err error)
}

func CreateStdOutSave(
	savedOutput map[string][]string,
	fn func(string, func(string)) (n int, err error),
) *StdLogSave {
	return &StdLogSave{
		savedOutput: savedOutput,
		fn:          fn,
	}
}

func (so *StdLogSave) Write(p []byte) (n int, err error) {
	// jsonData, _ := json.MarshalIndent(so.savedOutput, "", "  ")
	// fmt.Println(string(jsonData))

	return so.fn(string(p), func(x string) {
		so.mu.Lock()
		defer so.mu.Unlock()

		so.savedOutput[so.Group] = append(so.savedOutput[so.Group], x)
	})
}
