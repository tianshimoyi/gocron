package constants

const (
	TaskResourceTag      = "Task"
	NodeResourceTag      = "Node"
	HTTP200              = "It is used to indicate nonspecific success. The response body contains a representation of the resource."
	HTTP204              = "It is used to indicate nonspecific success. The response body contains nothing."
	HTTP201              = "Upon success, the HTTP response shall include a Location HTTP header that contains the resource URI of the created resource."
	HTTP400              = "Bad Request. It is used to indicate that incorrect parameters were passed to the request."
	HTTP403              = "Forbidden. The operation is not allowed given the current status of the resource."
	HTTP404              = "Not Found. It is used when a client provided a URI that cannot be mapped to a valid resource URI."
	HTTP409              = "Already exists"
	HTTP412              = "Precondition Failed. It is used when a condition has failed during conditional requests, e.g. when using ETags to avoid write conflicts."
	HTTP414              = "It is used to indicate that the server is refusing to process the request because the request URI is longer than the server is willing or able to process."
	HTTP500              = "Internal Error"
	GoCronUsernameHeader = "X-GOCRON-UID"
)
