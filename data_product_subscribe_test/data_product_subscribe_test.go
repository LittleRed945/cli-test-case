package dataproductsubscribe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
)

var ut = testutils.TestUtils{}
var EventCount = 10

type JSONSubData struct {
	Event     string      `json:"event"`
	Header    interface{} `json:"header"`
	Method    string      `json:"method"`
	Payload   interface{} `json:"payload"`
	Pk        []string    `json:"primarykey"`
	Product   string      `json:"product"`
	Seq       uint64      `json:"seq"`
	Subject   string      `json:"subject"`
	Table     string      `json:"table"`
	Timestamp string      `json:"timestamp"`
}

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
func FindJSON(data string) []string {
	var result []string
	stringStart := -1
	lvl := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '{' {
			if lvl == 0 {
				stringStart = i
			}
			lvl++
		} else if data[i] == '}' {
			lvl--
			if lvl == 0 {
				resultString := data[stringStart : i+1]
				result = append(result, resultString)
			}
		}
	}
	return result
}

func SubscribeDataProductCommand(productName string, subName string, partitions string, seq string) error {
	productName = ut.ProcessString(productName)
	subName = ut.ProcessString(subName)
	cmdString := "timeout 1 ../gravity-cli product sub "
	if productName != testutils.NullString {
		cmdString += productName + " "
	}
	// if subName != testutils.IgnoreString {
	// 	cmdString += "--name " + subName + " "
	// }
	if partitions != testutils.IgnoreString {
		cmdString += "--partitions " + partitions + " "
	}
	if seq != testutils.IgnoreString {
		cmdString += "--seq " + seq + " "
	}

	cmdString += "-s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(cmdString)
	if err != nil {
		return err
	}
	return nil
}

func DisplayData() error {
	resultStringList := FindJSON(ut.CmdResult.Stdout)
	fmt.Println(resultStringList)
	return nil
}

func PublishProductEvent(numbersOfEvents int) error {
	EventCount = numbersOfEvents
	for i := 0; i < EventCount; i++ {
		result := fmt.Sprintf(`{"id":%d, "name":"test%d", "kcal":%d, "price":%d}`, i+1, i+1, i*100, i+20)
		cmd := exec.Command(testutils.GravityCliString, "pub", Event, result, "-s", ut.Config.JetstreamURL)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDataProduct(dataProduct string) error {
	decription := "drink資料"
	schema := "./assets/schema.json"
	enabled := "true"
	cmd := exec.Command(testutils.GravityCliString, "product", "create", dataProduct, "--desc", decription, "--schema", schema, "--enabled="+enabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("create data product failed")
	}
	return nil
}

var DataProduct string
var Ruleset string
var Event string

const Method string = "create"
const Pk string = "id,name"
const RulesetDesc string = "drink創建事件"
const Handler string = "./assets/handler.js"
const Schema string = "./assets/schema.json"
const Enabled string = "true"

func CreateDataProductRuleset(dataProduct string, ruleset string) error {
	DataProduct = dataProduct
	Ruleset = ruleset
	Event = ruleset
	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", DataProduct, Ruleset, "--event", Event, "--method", Method, "--desc", RulesetDesc, "--pk", Pk, "--handler", Handler, "--schema", Schema, "--enabled="+Enabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("data product add ruleset failed")
	}
	return nil
}
func ValidateSubResult(partitionString string, seqString string) error {
	var err error
	var partitions int64 = -1
	var seq uint64 = 1
	if partitionString != testutils.IgnoreString {
		partitions, err = strconv.ParseInt(partitionString, 10, 64)
		if err != nil {
			return fmt.Errorf("轉換partition成int失敗: %s", err.Error())
		}
	}
	if seqString != testutils.IgnoreString {
		seq, err = strconv.ParseUint(seqString, 10, 64)
		if err != nil {
			return fmt.Errorf("轉換seq成uint失敗: %s", err.Error())
		}
	}
	if partitions != -1 {
		if ut.CmdResult.Err != nil && ut.CmdResult.Err.(*exec.ExitError).ExitCode() != 124 {
			return errors.New("sub指令失敗")
		}
		return nil
	}
	resultStringList := FindJSON(ut.CmdResult.Stdout)
	numbersOfEvents := uint64(EventCount) - seq + 1
	if uint64(EventCount+1) < seq {
		numbersOfEvents = 0
	}
	if uint64(len(resultStringList)) != numbersOfEvents {
		errString := fmt.Sprintf("Event數量與發佈數量不符合, 預期數量: %d, 獲取數量: %d", numbersOfEvents, len(resultStringList))
		return errors.New(errString)
	}
	for i, jsonData := range resultStringList {
		i := uint64(i)
		i = i + seq - 1
		var UnmarshalResult JSONSubData
		if err := json.Unmarshal([]byte(jsonData), &UnmarshalResult); err != nil {
			return errors.New("json unmarshal failed" + err.Error())
		}

		payloadString := fmt.Sprintf(`{"id":%d, "name":"test%d", "kcal":%d, "price":%d}`, i+1, i+1, i*100, i+20)
		expectPayload := ut.FormatJSONData(payloadString)

		if err := ut.AssertStringEqual(UnmarshalResult.Event, Event); err != nil {
			return err
		}
		if err := ut.AssertStringEqual(UnmarshalResult.Product, DataProduct); err != nil {
			return err
		}
		payloadByte, _ := json.Marshal(UnmarshalResult.Payload)
		resultPayload := string(payloadByte)
		if err := ut.AssertStringEqual(resultPayload, expectPayload); err != nil {
			return err
		}
		pkExpect := strings.Split(Pk, ",")
		for i := 0; i < len(pkExpect); i++ {
			if err := ut.AssertStringEqual(UnmarshalResult.Pk[i], pkExpect[i]); err != nil {
				return err
			}
		}
		// if err := ut.AssertUIntEqual(UnmarshalResult.Seq, i+1); err != nil {
		// 	return err
		// }
		// if err := ut.AssertStringEqual(UnmarshalResult.Method, Method); err != nil {
		// 	return err
		// }
	}
	return nil
}

func SubCommandFail() error {
	if ut.CmdResult.Err == nil || ut.CmdResult.Err.(*exec.ExitError).ExitCode() == 124 {
		return errors.New("使用該指令應該要失敗")
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有 data product "'(.*?)'"$`, CreateDataProduct)
	ctx.Given(`^已有 data product 的 ruleset "'(.*?)'" "'(.*?)'"$`, CreateDataProductRuleset)
	ctx.Given(`^已 publish "'(.*?)'" 筆 Event$`, PublishProductEvent)
	ctx.When(`^訂閱data product "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'"`, SubscribeDataProductCommand)
	ctx.Then(`^Cli 回傳 "'(.*?)'" 內 "'(.*?)'" 後所有事件資料$`, ValidateSubResult)
	ctx.Then(`^Cli 回傳指令失敗$`, SubCommandFail)
	ctx.Then(`^顯示資料`, DisplayData)
}
