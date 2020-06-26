package file

import (
	"github.com/coredns/coredns/plugin/file/tree"
	"github.com/coredns/coredns/plugin/transfer"

	"github.com/miekg/dns"
)

// Transfer implements the transfer.Transfer interface.
func (f File) Transfer(zone string, serial uint32) (<-chan []dns.RR, error) {
	z, ok := f.Zones.Z[zone]
	if !ok || z == nil {
		return nil, transfer.ErrNotAuthoritative
	}
	// get soa and apex
	apex, err := z.ApexIfDefined()
	if err != nil {
		return nil, err
	}

	ch := make(chan []dns.RR)
	go func() {
		ch <- apex

		z.Walk(func(e *tree.Elem, _ map[uint16][]dns.RR) error { ch <- e.All(); return nil })

		ch <- []dns.RR{apex[0]}
		close(ch)
	}()

	return ch, nil
}
