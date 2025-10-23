package common

// 模型类型

const (
	ModelClassText      int = 1 // 文本生成
	ModelClassEmbedding int = 2 // 向量嵌入
	ModelClassReranker  int = 3 // 重排序
	ModelClassVL        int = 4 // 视觉理解
	ModelClassImage     int = 5 // 图像生成
	ModelClassVideo     int = 6 // 视频生成
	ModelClassASR       int = 7 // 语音识别
	ModelClassTTS       int = 8 // 语音合成
	ModelClassMuti      int = 9 // 多模态模型
)

// 扩展能力

const (
	ModelExtAbilityTools     int = 1 // 工具调用
	ModelExtAbilityJson      int = 2 // 结构化输出
	ModelExtAbilityCache     int = 3 // Cache缓存
	ModelExtAbilityBatch     int = 4 // 批量推理
	ModelExtAbilityTry       int = 5 // 模型体验
	ModelExtAbilityFinetune  int = 6 // 模型调优
	ModelExtAbilityWebSearch int = 7 // 联网搜索
)
