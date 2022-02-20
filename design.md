# Web Crawler Design

I have assumed, inferred or made up some requirements for the purposes of the test and to fit the implementation within scope (4 hours). In reality, the requirements of this project would need a lot more thought and clarification.

## Requirements
- Accept a starting seed URL
- Crawl a single subdomain (or root domain), visiting each URL it finds on the domain, without following external URLs
- Each page visited visited should be printed, with a list of *links* found on that page
- It should be polite - it wouldn't look very good if the crawler took down monzo.com

**Scope**:
- This iteration of the crawler will only visit HTML pages. CSS, JS, PDF, TXT resources will not be accessed
- The crawler will visit each URL it finds anywhere on the HTML page
- On each page, the crawler will print links. I am taking 'links' to mean anchor tags within the pages HTML, so links to CSS, JS or other resources will not be printed
- Politeness should only cover rate limiting, not respecting robots.txt or any other directives

### Non-functional
- **Simple**: With small scope/budget, this project needs to be simple to run, operate and view the results of
- **Efficient**: The crawler should make good use of it's resources, including time


## Solution

### Decisions
- **Written in Go**: Simple, reliable, compiles to a single portable binary, easy to run, good concurrency model which will help with efficiency as we'll likely be constrained by IO
- **Single process**: Simple to run and deploy, no infrastructure required
- **Text interface**: Accepts a seed URL via a CLI, prints the sitemap to CLI output or a text file
- **Stateless**: As the crawler will be used on single domains, we can keep all state in memory, saving us from any complexities with persistence. We may hit a limit on large domains, such as Amazon, but it won't be used for that. Monzo.com had ~63k pages according to the Google query `site:monzo.com -site:community.monzo.com`

### Components
#### CLI
Basic CLI that accepts a seed URL and reports on progress. Progress will be reported using logging from across the whole system. Each URL visited should be logged.

#### URL Queue
Queue that stores URL targets for the fetchers to retrieve but is also responsible for rate limiting how often the domain is requested by limiting how often URLs are dequeued.

#### Fetchers
HTTP downloaders. Responsible for DNS resolution, downloading HTML pages, and queuing the content.

#### Page queue
Receives HTML pages and fans them out to parallel processors.

#### Sitemap Manager
Receives HTML pages, parses them and builds a `Sitemap` incrementally.

#### Sitemap
A model of a website, recording how pages are linked together.

#### Sitemap writer
An interface that accepts a `Sitemap` and writes it to media. In this instance our implementation will use a text file.

#### URL extractor
Receives HTML pages from the queue and parses any URLs from them, passing them to the URL filter. It will also standardise

#### URL filter
Receives URLs, filters out any previously before seen URLs and any non-HTML pages, and then adds them to the URL queue for crawling.

This component
