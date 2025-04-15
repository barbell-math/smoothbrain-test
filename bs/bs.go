package main

import (
	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func main() {
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterUpdateDepsTarget()
	sbbs.RegisterGoMarkDocTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.GoTargets{
		GenericTestTarget:  true,
		GenericBenchTarget: true,
		GenericFmtTarget:   true,
	})
	sbbs.RegisterMergegateTarget(sbbs.MergegateTargets{
		CheckFmt:             true,
		CheckReadmeGomarkdoc: true,
		CheckDepsUpdated:     true,
		CheckUnitTests:       true,
	})
	sbbs.Main("bs")
}
