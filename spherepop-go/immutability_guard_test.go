package spherepop

import "testing"

func TestHistoryIsImmutableAcrossPointerBoundaries(t *testing.T) {
	a := NewArbiter([]ObjectID{1, 2}, []RuleID{"identity"})
	input := PopEvent(1)
	proposal := NewProposal(input)
	*input.A = 2
	if _, err := a.Submit(proposal); err != nil { t.Fatal(err) }
	if _, ok := a.State().Committed[1]; !ok {
		t.Fatal("caller mutated a proposal through an original event pointer")
	}
	events := a.History().Events()
	*events[0].A = 2
	if _, ok := a.State().Committed[1]; !ok {
		t.Fatal("caller mutated authoritative history through Events()")
	}
	if _, ok := a.State().Committed[2]; ok {
		t.Fatal("history changed after caller modified exported event pointer")
	}
}

func TestHistoryAppendDetachesProposalEventFields(t *testing.T) {
	a := NewArbiter([]ObjectID{1, 2}, nil)
	p := NewProposal(BindEvent(1, 2, "edge"))
	if _, err := a.Submit(p); err != nil { t.Fatal(err) }
	*p.Events[0].Tag = "changed"
	if _, ok := a.State().Bound[Binding{1, 2, "edge"}]; !ok {
		t.Fatal("mutable proposal tag changed committed history")
	}
}
