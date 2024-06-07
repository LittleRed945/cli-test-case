package dataproductlist

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
)

const (
	blankString1 = "[null]"
	blankString2 = "[space]"
)

var ut = testutils.TestUtils{}

func TestFeatures(t *testing.T) {
	err := ut.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"./"},
			StopOnFailure: ut.Config.StopOnFailure,
			TestingT:      t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func CreateDataProductCommand(productAmount int, dataProduct string, description string, enabled string) error {
	dataProductNameBase := ut.ProcessString(dataProduct)
	dataProductName := dataProductNameBase
	for i := 0; i < productAmount; i++ {
		if i != 0 {
			dataProductName = dataProductNameBase + "_" + strconv.Itoa(i)
		}

		// TODO: 空格輸入不合預期
		if description != testutils.IgnoreString {
			description = ut.ProcessString(description)
		}
		// [space] [null] 處理, match到就把--desc參數去掉

		enabledString := ""
		if enabled != testutils.IgnoreString {
			if enabled == testutils.TrueString {
				enabledString += "--enabled"
			}
		}
		var cmd *exec.Cmd
		if description == blankString1 || description == blankString2 {
			cmd = exec.Command(testutils.GravityCliString, "product", "create", dataProductName, "--schema", "./assets/schema.json", enabledString)
		} else {
			cmd = exec.Command(testutils.GravityCliString, "product", "create", dataProductName, "--desc", description, "--schema", "./assets/schema.json", enabledString)
		}
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func AddRulesetCommand(RulesetAmount int, dataProduct string) error {
	dataProduct = ut.ProcessString(dataProduct)
	for i := 0; i < RulesetAmount; i++ {
		ruleset := dataProduct + "Created"
		if i != 0 {
			ruleset += strconv.Itoa(i)
		}
		cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", ruleset, "--enabled", "--method", "create", "--schema", "./assets/schema.json", "--pk", "id")
		err := cmd.Run()
		if err != nil {
			return errors.New(cmd.String())
		}
	}
	return nil
}

func PublishProductEvent(eventAmount int) error {
	for i := 0; i < eventAmount; i++ {
		event := "drinkCreated"
		payload := fmt.Sprintf(`{"id":%d, "name":"test%d", "kcal":0, "price":0}`, i, i)
		cmd := exec.Command(testutils.GravityCliString, "pub", event, payload)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("publish event failed: %s", err.Error())
		}
	}
	return nil
}

func ProductListCommand() error {
	cmd := testutils.GravityCliString + " product list"
	return ut.ExecuteShell(cmd)
}

func ProductListCommandSuccess(productAmount int, dataProduct string, rulesetAmount string, eventAmount string, description string, enabled string) error {
	outStr := ut.CmdResult.Stdout

	outStrList := strings.Split(outStr, "\n")
	if len(outStrList) != productAmount+3 { // 2 is header + 1 is the final empty line
		return errors.New("CLI returns error message: ProductAmount mismatches")
	}

	if strings.Compare(enabled, testutils.TrueString) == 0 {
		enabled = "enabled"
	} else {
		enabled = "disabled"
	}

	for i := 1; i <= productAmount; i++ {
		product := outStrList[2+i-1]
		dataProduct = ut.ProcessString(dataProduct)
		dataProductName := dataProduct
		if i != productAmount {
			dataProductName = dataProduct + "_" + strconv.Itoa(i)
		}

		productItem := strings.Fields(product)

		if productItem[0] != dataProductName {
			return errors.New("CLI returns error message: list ProductName mismatches")
		}

		index := 0
		if description != blankString1 && description != blankString2 {
			description = ut.ProcessString(description)
			if productItem[1+index] != description {
				return errors.New("CLI returns error message: list Description mismatches")
			}
			index++
		}

		if productItem[1+index] != enabled {
			return errors.New("CLI returns error message: list Enabled mismatches")
		}

		if i == productAmount {
			if productItem[2+index] != rulesetAmount {
				return errors.New("CLI returns error message: list RulesetAmount mismatches")
			}

			if productItem[3+index] != eventAmount {
				return errors.New("CLI returns error message: list EventAmount mismatches")
			}

		} else {
			if productItem[2+index] != "0" {
				return errors.New("CLI returns error message: list RulesetAmount mismatches")
			}

			if productItem[3+index] != "0" {
				return errors.New("CLI returns error message: list EventAmount mismatches")
			}

		}
	}

	return nil
}

func ProductListCommandFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return fmt.Errorf("List command should fail")
}

func CheckError(errMsg string) error {
	outStr := ut.CmdResult.Stderr
	if strings.Contains(outStr, errMsg) {
		return nil
	}
	return errors.New(errMsg)
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^Nats has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create "'(.*?)'" data products with "'(.*?)'" using parameters "'(.*?)'" "'(.*?)'"$`, CreateDataProductCommand)
	ctx.Given(`^Create "'(.*?)'" rulesets for "'(.*?)'"$`, AddRulesetCommand)
	ctx.Given(`^Publish the event "'(.*?)'" times$`, PublishProductEvent)
	ctx.When(`^Listing all data products using gravity-cli$`, ProductListCommand)
	ctx.Then(`^The CLI returns "'(.*?)'" data products. The final product has the name "'(.*?)'", with "'(.*?)'" rulesets, and a total of "'(.*?)'" events published. Each data product has a description of "'(.*?)'" and an enabled status of "'(.*?)'".$`, ProductListCommandSuccess)
	ctx.Then(`^CLI returns exit code 1$`, ProductListCommandFail)
	ctx.Then(`^The error message should be "'(.*?)'"$`, CheckError)
}
