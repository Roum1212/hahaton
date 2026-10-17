package aggregate_event_chat_command

type Append struct {
	sequence uint64
}

func (x Append) GetSequence() uint64 {
	return x.sequence
}

func NewAppend(sequence uint64) Append {
	return Append{
		sequence: sequence,
	}
}
