package common

const (
	NumServer = 12
	Transport = "tcp"
	Localhost = "localhost:%d"
	Port      = 50051

	Space64 = "                                                                "
)

const (
	InputBufferSize = 64

	// 100 -> laggy
	// 33  -> enough
	InputSyncMil = 16
)
