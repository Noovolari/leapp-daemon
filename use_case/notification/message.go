package notification

type MessageType int
const (
  MfaTokenRequest MessageType = iota
)

type MfaTokenRequestData struct {
  SessionId string
}

type Message struct {
  MessageType MessageType
  Data        string
}
