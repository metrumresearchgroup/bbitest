package bbitest

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

type testConfig struct {
	bbiOption string
	goldenExt string
}

var testConfigs = []testConfig{
	{
		"--json",
		".json",
	},
	{
		"",
		".txt",
	},
}

var SummaryHappyPathTestMods = []string{
	"acop",             // basic model
	"12",               // bootstrap model with no $COV step
	"example2_saemimp", // two est methods SAEM => IMP (Has large_condition_number and eigenvalue_issues)
	"example2_itsimp",  // two est methods ITS => IMP (No Prior)
	"example2_bayes",   // Bayes (5 est methods, from NONMEM examples)
	"iovmm",            // Mixture model. Also has parameter_near_boundary and final_zero_gradient heuristics.
	"acop-iov",         // fake model with 62 OMEGAS (fake iov)
}

func TestSummaryHappyPath(tt *testing.T) {
	t := wrapt.WrapT(tt)

	for _, mod := range SummaryHappyPathTestMods {
		for _, tc := range testConfigs {

			commandAndArgs := []string{
				"nonmem",
				"summary",
				filepath.Join(SUMMARY_TEST_DIR, mod, mod+".lst"),
			}

			if tc.bbiOption != "" {
				commandAndArgs = append(commandAndArgs, tc.bbiOption)
			}

			output, err := executeCommand(context.Background(), "bbi", commandAndArgs...)

			t.R.NoError(err)
			t.R.NotEmpty(output)

			gtd := GoldenFileTestingDetails{
				outputString:   output,
				goldenFilePath: filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden"+tc.goldenExt),
			}

			if os.Getenv("UPDATE_SUMMARY") == "true" {
				UpdateGoldenFile(t, gtd)
			}

			RequireOutputMatchesGoldenFile(t, gtd)
		}
	}
}

type testModWithArg struct {
	mod        string
	bbiArg     string
	errorRegEx string
}

var SummaryArgsTestMods = []testModWithArg{
	{ // from bbr example project. Has a PRDERR that causes shrinkage file to be missing.
		"66",
		"--no-shk-file",
		`\-\-no\-shk\-file`,
	},
	{ // copy of acop with .grd deleted
		"acop_no_grd",
		"--no-grd-file",
		`\-\-no\-grd\-file`,
	},
	{ // Bayesian model testing --ext-file flag. Also has a large condition number.
		"1001",
		"--ext-file=1001.1.TXT",
		`\-\-ext\-file`,
	},
}

func TestSummaryArgs(tt *testing.T) {
	t := wrapt.WrapT(tt)

	for _, tm := range SummaryArgsTestMods {
		for _, tc := range testConfigs {

			mod := tm.mod

			commandAndArgs := []string{
				"nonmem",
				"summary",
				filepath.Join(SUMMARY_TEST_DIR, mod, mod+".lst"),
			}

			if tc.bbiOption != "" {
				commandAndArgs = append(commandAndArgs, tc.bbiOption)
			}

			// try without flag and get error
			output, err := executeCommandNoErrorCheck(context.Background(), "bbi", commandAndArgs...)
			t.R.NotNil(err)
			errorMatch, _ := regexp.MatchString(tm.errorRegEx, output)
			t.R.True(errorMatch)

			// append flag and get success
			commandAndArgs = append(commandAndArgs, tm.bbiArg)
			output, err = executeCommand(context.Background(), "bbi", commandAndArgs...)

			t.R.NoError(err)
			t.R.NotEmpty(output)

			gtd := GoldenFileTestingDetails{
				outputString:   output,
				goldenFilePath: filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden"+tc.goldenExt),
			}

			if os.Getenv("UPDATE_SUMMARY") == "true" {
				UpdateGoldenFile(t, gtd)
			}

			RequireOutputMatchesGoldenFile(t, gtd)
		}
	}
}

type SummaryErrorCase struct {
	testPath string
	errorMsg string
}

var SummaryErrorCases = []SummaryErrorCase{
	{
		"acop", // points to directory instead of file
		noSuchFileError,
	},
	{
		"acop/aco", // misspelled filename
		noSuchFileError,
	},
	{
		"aco", // non-existing directory
		noSuchFileError,
	},
	{
		"acop/acop.ls", // no file at that extension
		wrongExtensionError,
	},
	{
		"acop/acop.ext", // wrong (but existing) file
		wrongExtensionError,
	},
}

func TestSummaryErrors(tt *testing.T) {
	t := wrapt.WrapT(tt)

	for _, tc := range SummaryErrorCases {

		commandAndArgs := []string{
			"nonmem",
			"summary",
			filepath.Join(SUMMARY_TEST_DIR, tc.testPath),
		}

		// try without flag and get error
		output, err := executeCommandNoErrorCheck(context.Background(), "bbi", commandAndArgs...)
		t.R.NotNil(err)
		errorMatch, _ := regexp.MatchString(tc.errorMsg, output)
		t.R.True(errorMatch)
	}
}

func TestSummaryHappyPathNoExtension(tt *testing.T) {
	t := wrapt.WrapT(tt)

	mod := "acop" // just testing one model

	commandAndArgs := []string{
		"nonmem",
		"summary",
		filepath.Join(SUMMARY_TEST_DIR, mod, mod), // adding no extension should work
	}

	output, err := executeCommand(context.Background(), "bbi", commandAndArgs...)

	t.R.NoError(err)
	t.R.NotEmpty(output)

	gtd := GoldenFileTestingDetails{
		outputString:   output,
		goldenFilePath: filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden.txt"),
	}

	RequireOutputMatchesGoldenFile(t, gtd)
}
