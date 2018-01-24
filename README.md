# trackingPixel
This service should be implemented by a customer adding the following tag "<img src=“http//stores.example.com/ tracking/track.gif”></img>” to each of their product pages, where those product pages are marked up with Open Graph tags and where stores.example.com is a unique domain for that customer that points to our multi-tenant SaaS environment.  Assuming that all requests for /tracking/track.gif are already routed to the service, the service needs to respond to the request with a 1px by 1px transparent gif as quickly as possible.  The service also needs to scrape the page that requested the tracking pixel for all open graph markup such that the scraped data can be shown to an admin user at some point in the future.  Assume that there is no urgency for the admin user to see the data.  Do not worry about how the data is eventually presented to the user or how they get access to it, only that it needs to be stored in such a way that an application independent of this service will be able to access it and present it to the user.

# environment & build
 require Go

 $ go get github.com/dyatlov/go-opengraph/opengraph
 
 $ go build ../src/github.com/gengwensu/trackingPixel/trackingPixel.go
 $./trackingPixel -outfile="result.txt" &
