Feature: Data Product list

Scenario:
    Given 已開啟服務 nats
    Given 已開啟服務 dispatcher
#Scenario
    Scenario: 針對data product 的 event list，成功情境
    When 創建 "'<ProductAmount>'" 個 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'"
    Then Cli 回傳第 "'<ProductAmount>'" 個 "'<ProductName>'" 建立成功
    When 對"'<ProductName>'" 創建 "'<RulesetAmount>'" 個 ruleset
    Then ruleset 創建成功
    When 對Event做 "'<EventAmount>'" 次 publish 
    Then publish 成功
    When 使用gravity-cli 列出所有 data product
    Then 回傳 data product ProductAmount = "'<ProductAmount>'", ProductName = "'<ProductName>'", Description = "'<Description>'", Enabled="'<Enabled>'", RulesetAmount="'<RulesetAmount>'", EventAmount="'<EventAmount>'"
    Examples:
        |  ID   | ProductName | Description | Enabled | RulesetAmount | EventAmount | ProductAmount |
        |  M(1) | [a]x240     | description | [false] | 0             | 0           | 1             | #pass
        |  M(2) | drink       | [a]x32768   | [true]  | 0             | 0           | 1             | #pass 
        |  M(3) | drink       | " "         | [true]  | 1             | 500         | 1             | #1000000
        |  M(4) | drink       | ""          | [true]  | 1             | 1           | 100           | # pass

#Scenario
    Scenario: 針對data product 的 event list，未建立任何data product情境
    When 使用gravity-cli 列出所有 data product
    Then 回傳 Error: No available products

    Examples:
        |  ID   | ProductAmount |
        |  M(1) | 0             | #