package JudgeLight_Go

import (
	"golang.org/x/sys/unix"
)

func SetProcLimit(
	cpuTimeLimit int,
	memoryLimit int,
	stackLimit int,
	outputSizeLimit int,
) error {
	var rl = unix.Rlimit{}

	// cpu time limit (MS)
	if cpuTimeLimit != 0 {
		rl.Cur = uint64(cpuTimeLimit/1000 + 1)
		rl.Max = rl.Cur
		if err := unix.Setrlimit(unix.RLIMIT_CPU, &rl); err != nil {
			return err
		}
	}

	// memory limit (KiB)
	// see: https://ux.stackexchange.com/questions/13815/files-size-units-kib-vs-kb-vs-kb
	if memoryLimit != 0 {
		rl.Cur = uint64(memoryLimit * 1024)
		rl.Max = rl.Cur
		if err := unix.Setrlimit(unix.RLIMIT_DATA, &rl); err != nil {
			return err
		}

		rl.Cur = rl.Cur * 2
		rl.Max = rl.Cur
		if err := unix.Setrlimit(unix.RLIMIT_AS, &rl); err != nil {
			return err
		}
	}

	// stack limit (KiB)
	if stackLimit != 0 {
		rl.Cur = uint64(stackLimit * 1024)
		rl.Max = rl.Cur
		if err := unix.Setrlimit(unix.RLIMIT_STACK, &rl); err != nil {
			return err
		}
	}

	// output size limit (B)
	if outputSizeLimit != 0 {
		rl.Cur = uint64(outputSizeLimit)
		rl.Max = rl.Cur
		if err := unix.Setrlimit(unix.RLIMIT_FSIZE, &rl); err != nil {
			return err
		}
	}
	return nil
}

func SetProcUser(uid int, gid int) error {
	// set uid
	if uid != 0 {
		if err := unix.Setuid(uid); err != nil {
			return err
		}
	}

	// set gid
	if gid != 0 {
		if err := unix.Setgid(gid); err != nil {
			return err
		}
	}
	return nil
}

func SetProcStream(inputFilePath string, outputFilePath string, errorFilePath string) error {
	// set input stream
	if inputFilePath != "" {
		// open file
		if fd, err := unix.Open(inputFilePath, unix.O_RDONLY, 666); err != nil {
			return err
		} else {
			// dup pipe
			if err := unix.Dup2(fd, unix.Stdin); err != nil {
				return err
			}
		}
	}

	// set output stream
	if outputFilePath != "" {
		if fd, err := unix.Open(outputFilePath, unix.O_WRONLY|unix.O_CREAT, 0666); err != nil {
			return err
		} else {
			if err := unix.Dup2(fd, unix.Stdout); err != nil {
				return err
			}
		}
	}

	// set error stream
	if errorFilePath != "" {
		if fd, err := unix.Open(errorFilePath, unix.O_WRONLY|unix.O_CREAT, 0666); err != nil {
			return err
		} else {
			if err := unix.Dup2(fd, unix.Stderr); err != nil {
				return err
			}
		}
	}
	return nil
}
