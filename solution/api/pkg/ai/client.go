package ai

import "context"

const (
	GenerateAdTextPrompt = "Сгенерируй текст рекламного объявления из 2-3 предложений на основе предоставленных данных. Учитывай название рекламной кампании, название рекламодателя, определённый язык и контекст. Тон объявления должен быть информативным и привлекающим внимание."
	GenerateAdTextMsg    = "Название кампании: %s\nРекламодатель: %s\nКонтекст: %s\nЯзык: %s"
)

type CallRequest struct {
	Prompt         string
	AdTitle        string
	AdvertiserName string
	Message        string
	Context        string
}

type Client interface {
	Call(ctx context.Context, req CallRequest) (string, error)
}
