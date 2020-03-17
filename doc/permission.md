# 權限設計

權限項目分為：temperature（體溫登記與瀏覽） and account（人員權限變更）
人員角色分為：管理員、老師、學生，只有職稱功能，不代表權限

每種權限可有三種視野：
- 個人
- 班級
- 全校

各種角色舉例（身份、體溫權限、人員權限）：
admin：管理員、全校、全校
班級導師：老師、班級、班級
學生：學生、個人、個人
衛生股長：學生、班級、個人
護理師：老師、全校、個人
一般老師：老師、個人、個人