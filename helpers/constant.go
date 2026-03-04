package helpers

const (
	// API Response
	Success = "success"
	Error   = "error"

	// Error JWT
	AuthenticationFailed = "Authentication failed"
	MissingToken         = "Missing Token"
	ValidateJWTFailed    = "Failed to Validate JWT"

	// Login Response
	LoginSuccess       = "Login successful"
	InvalidCredentials = "Invalid email or password"
	InvalidRequest     = "Invalid request payload"

	// Order Response
	OrderPlaced              = "Order placed successfully"
	OrderPlacementFailed     = "Failed to place order"
	OrderListRetrieved       = "Order list retrieved successfully"
	OrderListRetrievalFailed = "Failed to retrieve order list"

	// Trade Response
	TradeListRetrieved       = "Trade list retrieved successfully"
	TradeListRetrievalFailed = "Failed to retrieve trade list"

	// Markets Response
	MarketSnapshotRetrieved       = "Market snapshot retrieved successfully"
	MarketSnapshotRetrievalFailed = "Failed to retrieve market snapshot"
)
