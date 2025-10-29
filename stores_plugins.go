package coze

import (
	"context"
	"net/http"
)

func (r *storesPlugins) List(ctx context.Context, req *ListStoresPluginsReq, options ...CozeAPIOption) (NumberPaged[ProductPlugin], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[ProductPlugin], error) {
			response := new(listStoresPluginsResp)
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method:  http.MethodGet,
				URL:     "/v1/stores/plugins",
				Body:    req.toReq(request),
				options: options,
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[ProductPlugin]{
				response: response.HTTPResponse,
				HasMore:  response.Data.HasMore,
				Data:     response.Data.Items,
				LogID:    response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

type ProductPlugin struct {
	baseModel
	Metainfo   *ProductMetaInfo   `json:"metainfo,omitempty"`
	PluginInfo *ProductPluginInfo `json:"plugin_info,omitempty"`
}

type ProductMetaInfo struct {
	Name        string           `json:"name,omitempty"`        // 插件的名称。
	Description string           `json:"description,omitempty"` // 插件的描述信息，用于说明插件的功能、用途和特点。
	Category    *ProductCategory `json:"category,omitempty"`    // 插件所属的分类信息，包含分类 ID 和分类名称。
	IconURL     string           `json:"icon_url,omitempty"`    // 插件的图标 URL。
	EntityType  string           `json:"entity_type,omitempty"` // 实体类型，当前仅支持 plugin，表示插件。
	EntityID    string           `json:"entity_id,omitempty"`   // 插件 ID。
	ListedAt    int64            `json:"listed_at,omitempty"`   // 插件上架时间，以 Unix 时间戳格式表示，单位为秒。
	PaidType    string           `json:"paid_type,omitempty"`   // 插件的付费类型，固定为 paid，即付费插件。
	ProductID   string           `json:"product_id,omitempty"`  // 商品的 ID，用于在插件商店中唯一标识该插件商品。
	IsOfficial  bool             `json:"is_official,omitempty"` // 标识插件是否为官方发布。
}

type ProductCategory struct {
	ID   string `json:"id,omitempty"`   // 插件分类 ID。
	Name string `json:"name,omitempty"` // 插件分类名称。
}

type ProductPluginInfo struct {
	Heat                   int64   `json:"heat,omitempty"`                      // 插件的热度值，用于表示该插件的受欢迎程度或使用频率。数值越大表示热度越高。
	CallCount              int64   `json:"call_count,omitempty"`                // 插件的调用量，表示该插件的累计调用次数。
	Description            string  `json:"description,omitempty"`               // 插件的描述信息，用于说明插件的功能、用途和特点。
	SuccessRate            float64 `json:"success_rate,omitempty"`              // 插件的调用成功率，以小数形式表示，数值范围为 0 到 1。
	BotsUseCount           int64   `json:"bots_use_count,omitempty"`            // 该插件在智能体或工作流中的累计关联次数。
	FavoriteCount          int64   `json:"favorite_count,omitempty"`            // 插件的收藏量，表示该插件被用户收藏的总次数。
	IsCallAvailable        bool    `json:"is_call_available,omitempty"`         // 标识该插件当前是否可被调用。
	TotalToolsCount        int64   `json:"total_tools_count,omitempty"`         // 插件包含的工具总数。
	AvgExecDurationMs      float64 `json:"avg_exec_duration_ms,omitempty"`      // 插件执行的平均耗时，单位为毫秒。
	AssociatedBotsUseCount int64   `json:"associated_bots_use_count,omitempty"` // 当前扣子商店中关联了该插件的智能体数量。
}

type ListStoresPluginsReq struct {
	Keyword     *string  `query:"keyword" json:"-"`      // 插件搜索的关键词，支持模糊匹配。
	IsOfficial  *bool    `query:"is_official" json:"-"`  // 是否为扣子官方插件。默认返回官方插件和三方插件。
	CategoryIDs []string `query:"category_ids" json:"-"` // 插件分类 ID 列表，用于筛选指定多个分类下的插件。默认为空，即返回所有分类下的插件。可以通过查询插件分类API 获取对应的插件分类 ID。
	SortType    *string  `query:"sort_type" json:"-"`    // 排序类型，用于指定返回插件的排序方式。支持的排序方式如下所示：
	PageNum     int      `query:"page_num" json:"-"`
	PageSize    int      `query:"page_size" json:"-"`
}

type ListStoresPluginsResp struct {
	Items   []*ProductPlugin `json:"items"`
	HasMore bool             `json:"has_more"`
}

func (r ListStoresPluginsReq) toReq(page *pageRequest) *ListStoresPluginsReq {
	return &ListStoresPluginsReq{
		Keyword:     r.Keyword,
		IsOfficial:  r.IsOfficial,
		CategoryIDs: r.CategoryIDs,
		SortType:    r.SortType,
		PageNum:     page.PageNum,
		PageSize:    page.PageSize,
	}
}

type listStoresPluginsResp struct {
	baseResponse
	Data *ListStoresPluginsResp `json:"data"`
}

type storesPlugins struct {
	core *core
}

func newStoresPlugins(core *core) *storesPlugins {
	return &storesPlugins{core: core}
}
