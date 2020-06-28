package collector

var (
	Scrapers = map[Scraper]bool{
		CpuLoad{}: true,
	}
)

