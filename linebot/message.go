package main

import (
	"linebot/gurunavi"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/dustin/go-humanize"
)

// レストラン情報をテキストで返す
func TextRestaurants(g *gurunavi.GurunaviResponseBody) string {
	var t string
	for _, r := range g.Rest{
		t += r.Name + "\n" + r.URL + "\n"
	}
	return t
}


// レストラン情報をCarouselContainerポインター型で返す
func FlexRestaurants(g *gurunavi.GurunaviResponseBody) *linebot.CarouselContainer{
	var bcs []*linebot.BubbleContainer

	// 個々のレストラン情報をBubbleContainerにセットしてスライス格納
	for _, r := range g.Rest{
		b := linebot.BubbleContainer{
			Type: linebot.FlexContainerTypeBubble,
			Hero: setHero(r),
			Body: setBody(r),
			Footer: setFooter(r),
		}
		bcs = append(bcs, &b)
	}

	// BubbleContainerスライスをCarouselContainerに格納して返却
	return &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: bcs,
	}
}

func setHero(r *gurunavi.Rest) *linebot.ImageComponent {
	// 飲食店の画像URLがない場合
	if r.ImageURL.ShopImage1 == "" {
		return nil
	}

	return &linebot.ImageComponent{
		Type:        linebot.FlexComponentTypeImage,
		URL:         r.ImageURL.ShopImage1,
		Size:        linebot.FlexImageSizeTypeFull,
		AspectRatio: linebot.FlexImageAspectRatioType20to13,
		AspectMode:  linebot.FlexImageAspectModeTypeCover,
	}
}


func setBody(r *gurunavi.Rest) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			setRestaurantName(r),
			setLocation(r),
			setCategory(r),
			setBudget(r),
		},
	}
}


func setRestaurantName(r *gurunavi.Rest) *linebot.TextComponent {
	return &linebot.TextComponent{
		Type:   linebot.FlexComponentTypeText,
		Text:   r.Name,
		Wrap:   true,
		Weight: linebot.FlexTextWeightTypeBold,
		Size:   linebot.FlexTextSizeTypeMd,
	}
}

func setLocation(r *gurunavi.Rest) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeBaseline,
		Margin:  linebot.FlexComponentMarginTypeLg,
		Spacing: linebot.FlexComponentSpacingTypeSm,
		Contents: []linebot.FlexComponent{
			setSubtitle("エリア"),
			setDetail(r.Access.Station + "徒歩" + r.Access.Walk + "分"),
		},
	}
}

func setSubtitle(t string) *linebot.TextComponent {
	return &linebot.TextComponent{
		Type:  linebot.FlexComponentTypeText,
		Text:  t,
		Color: "#aaaaaa",
		Size:  linebot.FlexTextSizeTypeXs,
		Flex:  linebot.IntPtr(4),
	}
}


func setDetail(t string) *linebot.TextComponent {
	return &linebot.TextComponent{
		Type:  linebot.FlexComponentTypeText,
		Text:  t,
		Wrap:  true,
		Color: "#666666",
		Size:  linebot.FlexTextSizeTypeXs,
		Flex:  linebot.IntPtr(12),
	}
}

func setCategory(r *gurunavi.Rest) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeBaseline,
		Margin:  linebot.FlexComponentMarginTypeLg,
		Spacing: linebot.FlexComponentSpacingTypeSm,
		Contents: []linebot.FlexComponent{
			setSubtitle("ジャンル"),
			setDetail(r.Category),
		},
	}
}


func setBudget(r *gurunavi.Rest) *linebot.BoxComponent {
	var detail string
	// if文での型アサーションの使用　r.Budgetがfloat64型であれば
	if b, ok := r.Budget.(float64); ok {
		detail = "¥" + humanize.Comma(int64(b))
	} else {
		detail = "-"
	}

	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeBaseline,
		Margin:  linebot.FlexComponentMarginTypeLg,
		Spacing: linebot.FlexComponentSpacingTypeSm,
		Contents: []linebot.FlexComponent{
			setSubtitle("予算"),
			setDetail(detail),
		},
	}
}

func setFooter(r *gurunavi.Rest) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeVertical,
		Spacing: linebot.FlexComponentSpacingTypeXs,
		Contents: []linebot.FlexComponent{
			setButton("地図を見る", "https://www.google.com/maps"+"?q="+r.Latitude+","+r.Longitude),
			setButton("電話する", "tel:"+r.Tel),
			setButton("詳しく見る", r.URL),
		},
	}
}


func setButton(label string, uri string) *linebot.ButtonComponent {
	return &linebot.ButtonComponent{
		Type:   linebot.FlexComponentTypeButton,
		Style:  linebot.FlexButtonStyleTypeLink,
		Height: linebot.FlexButtonHeightTypeSm,
		Action: linebot.NewURIAction(label, uri),
	}
}
