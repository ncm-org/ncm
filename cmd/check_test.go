package cmd

import (
	"testing"

	"github.com/gookit/color"
)

func TestCheckCase1(t *testing.T) {
	message := ""
	doTestCheckCase(message)
}

func TestCheckCase2(t *testing.T) {
	message := "a"
	doTestCheckCase(message)
}

func TestCheckCase3(t *testing.T) {
	message := "网"
	doTestCheckCase(message)
}

func TestCheckCase4(t *testing.T) {
	message := "test"
	doTestCheckCase(message)
}

func TestCheckCase5(t *testing.T) {
	message := "type: "
	doTestCheckCase(message)
}

func TestCheckCase6(t *testing.T) {
	message := "type: t"
	doTestCheckCase(message)
}

func TestCheckCase7(t *testing.T) {
	message := "type: subject"
	doTestCheckCase(message)
}

func TestCheckCase8(t *testing.T) {
	message := "type: subjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjectsubjec"
	doTestCheckCase(message)
}

func TestCheckCase9(t *testing.T) {
	message := "a(a): a"
	doTestCheckCase(message)
}

func TestCheckCase10(t *testing.T) {
	message := "fix: subject"
	doTestCheckCase(message)
}

func TestCheckCase11(t *testing.T) {
	message := "fix(: subject"
	doTestCheckCase(message)
}

func TestCheckCase12(t *testing.T) {
	message := "fix(): subject"
	doTestCheckCase(message)
}

func TestCheckCase13(t *testing.T) {
	message := "fix( ): subject"
	doTestCheckCase(message)
}

func TestCheckCase14(t *testing.T) {
	message := "fix(  ): subject"
	doTestCheckCase(message)
}

func TestCheckCase15(t *testing.T) {
	message := "fix( a ): subject"
	doTestCheckCase(message)
}

func TestCheckCase16(t *testing.T) {
	message := "fix( ab ): subject"
	doTestCheckCase(message)
}

func TestCheckCase17(t *testing.T) {
	message := "fix(a): subject"
	doTestCheckCase(message)
}

func TestCheckCase18(t *testing.T) {
	message := "fix(0): subject"
	doTestCheckCase(message)
}

func TestCheckCase19(t *testing.T) {
	message := "fix(sc): subject"
	doTestCheckCase(message)
}

func TestCheckCase20(t *testing.T) {
	message := "fix(scope): subject"
	doTestCheckCase(message)
}

func TestCheckCase21(t *testing.T) {
	message := "fix(scope): subject"
	doTestCheckCase(message)
}

func TestCheckCase22(t *testing.T) {
	message := "fix(scope)"
	doTestCheckCase(message)
}

func TestCheckCaseForMergeCommitMessage1(t *testing.T) {
	message := "Merge remote-tracking branch 'origin/develop' into develop"
	doTestCheckCase(message)
}

func TestCheckCaseForMergeCommitMessage2(t *testing.T) {
	message := "Merge branch 'develop' of git@github.com:ncm-org/ncm.git into develop"
	doTestCheckCase(message)
}

func TestCheckCaseForMergeCommitMessage3(t *testing.T) {
	message := "Merge branch 'develop' into feature_develop"
	doTestCheckCase(message)
}

func TestCheckCaseForMergeCommitMessage4(t *testing.T) {
	message := `Merge remote-tracking branch 'origin/feature_1.2.5' into feature_1.2.5
			
	# Conflicts:
	#	module_base/libs/lmg.aar`
	doTestCheckCase(message)
}

func doTestCheckCase(message string) {
	color.BgGray.Printf("commit message: %s\n", message)

	errs := checkMessage(message)
	if len(errs) == 0 {
		color.Green.Println("Correct message")
		return
	}
	for _, err := range errs {
		handleError(err)
	}
}
