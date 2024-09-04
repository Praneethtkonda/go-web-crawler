# go-web-crawler
Implementation of web crawler in golang that crawls links and further continues the process for the found links


## Running the code
```bash
go-web-crawler git:(main) ✗ go run main.go
2024/09/04 11:57:51 Starting worker
2024/09/04 11:57:51 From worker 14
2024/09/04 11:57:51 From worker 0
2024/09/04 11:57:51 From worker 1
2024/09/04 11:57:51 From worker 28
2024/09/04 11:57:51 From worker 2
2024/09/04 11:57:51 From worker 22
2024/09/04 11:57:51 From worker 23
2024/09/04 11:57:51 From worker 24
2024/09/04 11:57:51 From worker 3
2024/09/04 11:57:51 From worker 4
2024/09/04 11:57:51 From worker 21
2024/09/04 11:57:51 From worker 59
2024/09/04 11:57:51 From worker 26
2024/09/04 11:57:51 From worker 60
2024/09/04 11:57:51 From worker 27
2024/09/04 11:57:51 From worker 25
2024/09/04 11:57:51 From worker 17
2024/09/04 11:57:51 Waiting for workers to finish
2024/09/04 11:57:52 [Worker 59] processing URL: http://quotes.toscrape.com/tag/simile/
2024/09/04 11:57:52 Crawling URL: http://quotes.toscrape.com/tag/simile/
2024/09/04 11:57:52 [Worker 17] processing URL: http://quotes.toscrape.com/tag/humor/page/1/
2024/09/04 11:57:52 Crawling URL: http://quotes.toscrape.com/tag/humor/page/1/
```