package JudgeLight_Go

type MemoryStatus struct {
	vmSize int
	vmRss  int
	vmData int
	vmStk  int
	vmExe  int
	vmLib  int
}

func MemoryUsage(pid int) (MemoryStatus, error) {
	ms := MemoryStatus{}
	// TODO
	// get status by /proc/<pid>/status
	return ms, nil
}
