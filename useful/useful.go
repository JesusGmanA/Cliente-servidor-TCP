package useful

import "../process"

type Client struct {
	ClientId int64
	Process  process.Process
}

const PORT = ":8043"
