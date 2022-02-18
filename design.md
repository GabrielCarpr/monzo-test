# Web Crawler Design

I have assumed, inferred or made up some requirements for the purposes of the test and to fit the implementation within scope (4 hours).

## Requirements
- Accept a starting seed URL
- Crawl a single subdomain (or root domain), visiting each URL it finds on the domain, but not follow external URLs
- Each URL visited should be printed, with a list of *links* found on that page
- It should be polite - it wouldn't look very good if the crawler took down monzo.com

**Scope**:
- The crawler will visit each URL it finds anywhere on the page
- On each page, the crawler will print links. I am taking 'links' to mean anchor tags within the pages HTML, so links to external CSS and JS resources will not be printed
- This iteration of the crawler will only visit HTML pages. CSS, JS, PDF, TXT resources will not be accessed


### Non-functional
- **Simple**: With small scope/budget, this project needs to be simple to run, operate and view the results of
- **Efficient**: The crawler should make good use of it's resources, including time


## Solution