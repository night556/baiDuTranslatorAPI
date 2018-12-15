#BAIDUTRANSLATORAPI

该项目调用百度翻译的API 实现终端翻译。可选择翻译语言，源语言采用自动识别。

语法：
1. `全英文`  默认翻译为中文
2. `中文`   默认翻译为英文。中文识别采用go的unicode包的scrits自动识别中文。
3. `语句/目标语言` 翻译成目标语言，目标语言需匹配百度翻译API文档中的标识。