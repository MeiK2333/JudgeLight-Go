package JudgeLight_Go

import (
	"golang.org/x/sys/unix"
)

type MemoryStatus struct {
	vmSize int
	vmRSS  int
	vmData int
	vmStk  int
	vmExe  int
	vmLib  int
}

func MemoryUsage(fd int) (MemoryStatus, error) {
	ms := MemoryStatus{}

	// get status by /proc/<pid>/status
	body := make([]byte, 4096)
	count, err := unix.Pread(fd, body, 0)
	if err != nil {
		return ms, err
	}

	// parse data by file
	for i := 0; i < count; i++ {
		switch body[i] {
		case 'V':
			goto V
		default:
			goto NEXTLINE
		}

	V:
		i++
		switch body[i] {
		case 'm':
			goto Vm
		default:
			goto NEXTLINE
		}

	Vm:
		i++
		switch body[i] {
		case 'R':
			i += 2
			goto VmRSS
		case 'D':
			i += 3
			goto VmData
		case 'S':
			goto VmS
		case 'E':
			i += 2
			goto VmExe
		case 'L':
			goto VmL
		default:
			goto NEXTLINE
		}

	VmRSS:
		i += 2
		ms.vmRSS = _getNumByVmLine(body[i:])
		goto NEXTLINE

	VmData:
		i += 2
		ms.vmData = _getNumByVmLine(body[i:])
		goto NEXTLINE

	VmS:
		i++
		switch body[i] {
		case 't':
			i++
			goto VmStk
		case 'i':
			i += 2
			goto VmSize
		default:
			goto NEXTLINE
		}

	VmStk:
		i += 2
		ms.vmStk = _getNumByVmLine(body[i:])
		goto NEXTLINE

	VmSize:
		i += 2
		ms.vmSize = _getNumByVmLine(body[i:])
		goto NEXTLINE

	VmExe:
		i += 2
		ms.vmExe = _getNumByVmLine(body[i:])
		goto NEXTLINE

	VmL:
		i++
		switch body[i] {
		case 'i':
			i++
			goto VmLib
		}

	VmLib:
		i += 2
		ms.vmLib = _getNumByVmLine(body[i:])
		goto NEXTLINE

	NEXTLINE:
		for body[i] != '\n' {
			i++
		}
	}

	return ms, nil
}

/**
in: '     42964 kB'
out: 42964
*/
func _getNumByVmLine(body []byte) int {
	offset := 0
	ans := 0
	start := false
	for {
		if !start {
			if _isDigit(body[offset]) {
				start = true
			} else {
				offset++
				continue
			}
		}
		if start {
			if _isDigit(body[offset]) {
				ans *= 10
				ans += int(body[offset] - '0')
				offset++
			} else {
				break
			}
		}
	}
	return ans
}

func _isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
