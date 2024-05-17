Feature: Data Product list

Scenario:
    Given 已開啟服務nats
    Given 已開啟服務dispatcher
    Given 已有gravity cli 工具  
#Scenario
    Scenario: 針對data product 的 event list，成功情境
    When 創建 "'<ProductAmount>'" 個 data product "'<ProductName>'" 使用參數 "'<Description>'" "'<Enabled>'"
    Then Cli 回傳 "'<ProductName>'" 建立成功
    When "'<ProductName>'" 創建 "'<RulesetAmount>'" 個 ruleset "'<Ruleset>'" 使用參數 "'<Method>'" "'<Event>'"
    Then ruleset 創建成功
    When 對 "'<Event>'" 做 "'<EventAmount>'" 次 publish 
    Then publish 成功
    When 使用gravity-cli 列出所有 data product
    Then 回傳 data product ProductName = "'<ProductName>'", Description = "'<Description>'", Enabled="'<Enabled>'", RulesetAmount="'<RulesetAmount>'", EventAmount="'<EventAmount>'"
    Examples:
        |  ID   | ProductName       | Description         | Enabled | RulesetAmount | EventAmount | ProductAmount |
        |  M(1) | max_len_str(256)  | description         | [false] | 0             | 0           | 1             |
        |  M(2) | drink             | ""                  | [true]  | 1             | 1           | 100           |
        |  M(3) | drink             | " "                 | [true]  | 1             | 1000000     | 1             |
        |  M(4) | drink             | max_len_str(32768)  | [true]  | 0             | 0           | 1             |