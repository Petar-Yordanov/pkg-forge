package engine

import "github.com/Petar-Yordanov/pkg-forge/manifest/parser"

type Events interface {
	OnDocStart(docIndex int)
	OnDocDone(docIndex int)

	OnEntryStart(e parser.Entry)
	OnEntrySkip(e parser.Entry, reason string)
	OnEntryDone(e parser.Entry)

	OnPreInstall(e parser.Entry, s parser.Step)
	OnInstall(e parser.Entry)
	OnUninstall(e parser.Entry)
	OnPostInstall(e parser.Entry, s parser.Step)
	OnValidation(e parser.Entry, s parser.Step)

	OnStep(e parser.Entry, s parser.Step)
	OnStepSkip(e parser.Entry, s parser.Step, reason string)

	OnError(e parser.Entry, err error)
}
