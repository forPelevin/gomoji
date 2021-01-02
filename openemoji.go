package gomoji

import (
	"context"
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
)

var (
	errInvalidJSONInResponse = errors.New("invalid JSON in response body")
)

// OpenEmojiProvider is an emoji provider that gets data from Open Emoji API (https://emoji-api.com).
type OpenEmojiProvider struct {
	accessKey string
}

// NewOpenEmojiProvider creates a new instance of OpenEmojiProvider.
func NewOpenEmojiProvider(accessKey string) *OpenEmojiProvider {
	return &OpenEmojiProvider{
		accessKey: accessKey,
	}
}

// AllEmojis gets all emojis from emoji-api.com.
func (o *OpenEmojiProvider) AllEmojis(ctx context.Context) ([]Emoji, error) {
	respBody, err := o.doRequest("https://emoji-api.com/emojis?access_key=" + o.accessKey)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if !gjson.ValidBytes(respBody) {
		return nil, errInvalidJSONInResponse
	}
	res := gjson.ParseBytes(respBody)
	gjsonEmojiList := res.Array()

	emojis := make([]Emoji, len(gjsonEmojiList))
	for i, gjonEmoji := range gjsonEmojiList {
		emojis[i] = Emoji{
			Slug:        gjonEmoji.Get("slug").String(),
			Character:   gjonEmoji.Get("character").String(),
			UnicodeName: gjonEmoji.Get("unicodeName").String(),
			CodePoint:   gjonEmoji.Get("codePoint").String(),
			Group:       gjonEmoji.Get("group").String(),
			SubGroup:    gjonEmoji.Get("subGroup").String(),
		}
	}

	return emojis, nil
}

func (o *OpenEmojiProvider) doRequest(uri string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(uri)

	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("do fasthttp request: %w", err)
	}

	return resp.Body(), nil
}
