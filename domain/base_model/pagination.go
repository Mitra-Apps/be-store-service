package base_model

type Pagination struct {
	Records      int32
	TotalRecords int32
	Limit        int32
	Page         int32
	TotalPage    int32
}

func (p *Pagination) GetOffset() int32 {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int32 {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int32 {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}
