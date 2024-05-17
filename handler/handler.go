package handler

import (
	"context"
	"github.com/riba2534/openwechat"
	"io"
	"log"
	"net/http"
	"openai-wechat/ai"
	"openai-wechat/config"
	"openai-wechat/consts"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func MessageHandler(msg *openwechat.Message) {
	if !msg.IsText() {
		return
	}
	ctx := context.Background()
	systemPrompt := config.Prompt
	switch {
	case strings.HasPrefix(msg.Content, config.C.WechatConfig.TextConfig.TriggerPrefix):
		// 文字回复
		if config.C.ContextConfig.SwitchOn {
			//go textSessionReplyHandler(ctx, msg, config.C.WechatConfig.TextConfig.TriggerPrefix, openai.GPT4, systemPrompt)
			go textSessionReplyHandler(ctx, msg, config.C.WechatConfig.TextConfig.TriggerPrefix, openai.GPT3Dot5Turbo, systemPrompt)
		} else {
			//go textReplyHandler(ctx, msg, config.C.WechatConfig.TextConfig.TriggerPrefix, openai.GPT4, systemPrompt)
			go textReplyHandler(ctx, msg, config.C.WechatConfig.TextConfig.TriggerPrefix, openai.GPT3Dot5Turbo, systemPrompt)
		}
	case strings.HasPrefix(msg.Content, config.C.WechatConfig.ImageConfig.TriggerPrefix):
		// 图片回复
		go imageReplyHandler(ctx, msg, config.C.WechatConfig.ImageConfig.TriggerPrefix)
	}
}

// 文字回复处理
func textReplyHandler(ctx context.Context, msg *openwechat.Message, prefix, model, systemPrompt string) {
	log.Printf("[text] Request: %s", msg.Context) // 输出请求消息到日志
	q := strings.TrimSpace(strings.TrimPrefix(msg.Content, prefix))
	reply := ai.CreateChatCompletion(ctx, model, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: q,
		},
	})
	log.Printf("[text] Response: %s", reply) // 输出回复消息到日志
	_, err := msg.ReplyText(reply)
	if err != nil {
		log.Printf("msg.ReplyText Error: %+v", err)
	}
}

// 带有上下文的文字回复
func textSessionReplyHandler(ctx context.Context, msg *openwechat.Message, prefix, model, systemPrompt string) {
	log.Printf("[text session] Request: %s", msg.Content) // 输出请求消息到日志
	user := func() string {
		s := msg.FromUserName
		if msg.IsSendBySelf() {
			s = msg.ToUserName
		}
		return s
	}()
	q := strings.TrimSpace(strings.TrimPrefix(msg.Content, prefix))
	reply := ai.GetSessionOpenAITextReply(ctx, q, user, model, systemPrompt)
	log.Printf("[text session] Response: %s", reply) // 输出回复消息到日志

	// 超过限长1420不能回复问题
	temp := []rune(reply)
	length := len(temp)
	mod := length / consts.MaxLength
	rdd := length % consts.MaxLength

	if mod == 0 {
		_, err := msg.ReplyText(reply)
		if err != nil {
			log.Printf("msg.ReplyText Error:%+v", err)
		}
	} else {
		for i := 0; i < mod; i++ {
			_, err := msg.ReplyText(string(temp[i*consts.MaxLength : consts.MaxLength*(i+1)]))
			if err != nil {
				log.Printf("msg.ReplyText Error:%+v", err)
			}
		}
		if rdd != 0 {
			_, err := msg.ReplyText(string(temp[mod*consts.MaxLength:]))
			if err != nil {
				log.Printf("msg.ReplyText Error:%+v", err)
			}
		}
	}
}

// 回复图片
func imageReplyHandler(ctx context.Context, msg *openwechat.Message, prefix string) {
	log.Printf("[image] Request: %s", msg.Content)
	q := strings.TrimSpace(strings.TrimPrefix(msg.Content, prefix))
	url := ai.CreateImageReply(ctx, q)
	if url == "" {
		log.Printf("[image] Response: url为空")
		msg.ReplyText(consts.ErrTips)
		return
	}
	log.Printf("[image] Response: url = %s", url)
	image, err := downloadImage(url)
	if err != nil {
		log.Printf("[image] downloadImage err, err=%+v", err)
		msg.ReplyText(consts.ErrTips)
		return
	}
	_, err = msg.ReplyImage(image)
	if err != nil {
		log.Printf("msg.ReplyImage Error: %+v", err)
	}
}

func downloadImage(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("downloadImage failed, err=%+v", err)
		return nil, err
	}
	return response.Body, nil
}
