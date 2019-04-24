package session

type SessionMatcher interface {
	Match(sess *Session) (*ResponseMatch, error)
}
