package subspec

import (
	"testing"
)

func TestSubtreeSpecification(t *testing.T) {
	// Verify parsing of valid string-based SubSpec values
	for idx, raw := range []any{
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter and:{item:1.3.6.1.4.1,or:{item:cn,item:2.5.4.7}}}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter or:{item:1.3.6.1.4.1,not:item:1.3.6.1.5.5,and:{item:cn,item:2.5.4.7}}}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter item:1.3.6.1.4.1.56521}`,
		`{minimum 1, maximum 1}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter not:item:1.3.6.1.4.1.56521}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", specificExclusions { chopBefore "n=14", chopAfter "n=555", chopAfter "n=74,n=6" }, minimum 1, maximum 1, specificationFilter item:1.3.6.1.4.1.56521}`,
		`{}`,
	} {
		if v, err := New(raw); err != nil {
			t.Errorf("%s[%d] failed: %v", t.Name(), idx, err)
			return
		} else if got := v.String(); got != raw {
			t.Errorf("%s[%d] failed:\n\twant: '%s'\n\tgot: '%s'",
				t.Name(), idx, raw, got)
		}
	}
}

func TestSubtreeSpecification_codecov(t *testing.T) {
	New(nil)
	New(``)
	New(`X`)
	New(byte(33))
	New(`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum -1, maximum 1, specificationFilter or:{item:1.3.6.1.4.1,not:item:1.3.6.1.5.5,and:{item:cn,item:2.5.4.7}}}`)

	_, _, _ = subtreeBase(rune(11))
	_, _, _ = subtreeBase(`value:...`)
	_, _ = subtreeRefinement(nil)

	_, _ = subtreeRefinement("any:{...}")

	var spec SubtreeSpecification
	spec.Base = "cn=1,cn=2,cn=3"
	spec.ChopSpecification.Exclusions = SpecificExclusions{
		SpecificExclusion{}}
	spec.ChopSpecification = ChopSpecification{}
	spec.SpecificationFilter = RefinementAnd{}
	_ = spec.String()

	_, _, _ = subtreeExclusions("{", 0)
	_, _, _ = subtreeExclusions("{chopBefore:cn=y,chopAfter:cn=x}", 0)
	_, _, _ = deconstructExclusions("{chopAfter:cn=x}", 0)

	var orref RefinementOr
	orref.Push(nil)
	_ = orref.String()
	orref.Index(2)
	orref.isRefinement()
	oi1, _ := parseOr("item:2.6.5.0")
	orref.Push(oi1)
	orref.Push("item:2.6.5.5")
	orref = append(orref, RefinementItem(``))

	var andref RefinementAnd
	andref.Push(nil)
	_ = andref.String()
	andref.Index(2)
	andref.isRefinement()
	ai1, _ := parseAnd("item:2.6.5.0")
	andref.Push(ai1)
	andref.Push("item:2.6.5.5")
	andref = append(andref, RefinementItem(``))

	var excls SpecificExclusions
	_ = excls.String()

	var excl SpecificExclusion
	_ = excl.String()

	var iref RefinementItem
	iref.Choice()
	_ = iref.String()
	iref.Len()
	iref.Index(1)
	iref.isRefinement()

	var nref RefinementNot
	_ = nref.String()
	nref.Len()
	nref.Index(1)
	nref.isRefinement()

	var ivref invalidRefinement
	ivref.Index(2)
	ivref.isRefinement()
	ivref.IsZero()
	ivref.Len()
	ivref.Choice()
	_ = ivref.String()

	checkSubtreeEncaps(`fjhdjk`)
	checkSubtreeEncaps(`{..`)
	subtreeExclusions(`F`, 0)
	subtreeExclusions(`a`, 0)

	parseItem("item:something")
	parseItem("item:")
	parseItem(":something")
	parseItem(":")
	parseItem("")

	parseNot("x")
	parseComplexRefinement("and", "{bogus}")

}
